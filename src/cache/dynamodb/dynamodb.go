package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/utils/anonymize"
)

type DRV struct {
	ctx       context.Context
	keyPrefix string
	client    *dynamodb.Client
	tableName string
}

func New(ctx context.Context, keyPrefix string, conf config.DynamoDB) (DRV, error) {
	var cfg aws.Config
	var err error

	if conf.Region == "" || conf.Region == "global" {
		log.Info().
			Msg("Using a global DynamoDB table")
		cfg, err = awsconfig.LoadDefaultConfig(ctx)
	} else {
		log.Info().
			Str("region", conf.Region).
			Msg("Using a regional DynamoDB table")
		cfg, err = awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(conf.Region))
	}

	if err != nil {
		log.Error().Err(err).Msg("Error loading AWS config")
		return DRV{}, err
	}

	var client *dynamodb.Client
	if conf.Endpoint != "" {
		log.Warn().
			Str("endpoint", conf.Endpoint).
			Msg("Using custom endpoint")

		cfg.BaseEndpoint = &conf.Endpoint
		resolver := dynamodb.NewDefaultEndpointResolverV2()
		client = dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolverV2(resolver))
	} else {
		client = dynamodb.NewFromConfig(cfg)
	}

	return DRV{ctx, keyPrefix, client, conf.Table}, nil
}

func (drv DRV) Close() {}

func (drv DRV) Set(k string, v any, ttl ...time.Duration) error {
	// Create a hash of the key to store in dynamodb.
	keyHash := anonymize.CalculateHashBase64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	// Serialize the value to JSON.
	valJSON, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("dynamodb.Set(): error marshaling value: %w", err)
	}

	// Encrypt the value using the original key (not the hash).
	val, err := anonymize.Encrypt(valJSON, k)
	if err != nil {
		return fmt.Errorf("dynamodb.Set(): error encrypting value: %w", err)
	}

	// Set the key-value pair in dynamodb.
	attributes := map[string]types.AttributeValue{
		"Key":   &types.AttributeValueMemberS{Value: keyHash},
		"Value": &types.AttributeValueMemberS{Value: val},
	}

	if len(ttl) > 0 {
		expirationTime := time.Now().Add(ttl[0]).Unix()
		attributes["TTL"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", expirationTime)}
	}

	_, err = drv.client.PutItem(drv.ctx, &dynamodb.PutItemInput{
		TableName: aws.String(drv.tableName),
		Item:      attributes,
	})
	if err != nil {
		return fmt.Errorf("dynamodb.Set(): error setting KV in dynamodb: %w", err)
	}

	return nil
}

func (drv DRV) Get(k string, o any) error {
	// Create a hash of the key to retrieve from dynamodb.
	keyHash := anonymize.CalculateHashBase64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	// Get the value from dynamodb.
	result, err := drv.client.GetItem(drv.ctx, &dynamodb.GetItemInput{
		TableName: aws.String(drv.tableName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: keyHash},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb.Get(): error getting value from dynamodb for key %v: %w", keyHash, err)
	}

	if result.Item == nil {
		log.Trace().
			Str("key_hash", keyHash).
			Msg("Found no value in dynamodb")
		return nil
	}

	// Required because TTL isn't guaranteed to remove the item immediately
	if result.Item["TTL"] != nil {
		expirationTime, err := strconv.ParseInt(result.Item["TTL"].(*types.AttributeValueMemberN).Value, 10, 64)
		if err != nil {
			return fmt.Errorf("dynamodb.Get(): error parsing TTL value for key hash %v: %w", keyHash, err)
		}

		expiresIn := time.Until(time.Unix(expirationTime, 0))
		if expiresIn < 0 {
			log.Trace().
				Str("key_hash", keyHash).
				Msg("Value has expired")
			return nil
		}
	}

	val := result.Item["Value"].(*types.AttributeValueMemberS).Value

	// Decrypt the value using the original key (not the hash).
	valJSON, err := anonymize.Decrypt(val, k)
	if err != nil {
		return fmt.Errorf("dynamodb.Get(): error decrypting value: %w", err)
	}

	// Deserialize the value from JSON.
	if err := json.Unmarshal(valJSON, o); err != nil {
		return fmt.Errorf("dynamodb.Get(): error unmarshaling value from dynamodb for key %v: %w", keyHash, err)
	}

	return nil
}

func (drv DRV) GetTTL(k string) (time.Duration, error) {
	// Create a hash of the key to retrieve from dynamodb.
	keyHash := anonymize.CalculateHashBase64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	// Get the value from dynamodb.
	result, err := drv.client.GetItem(drv.ctx, &dynamodb.GetItemInput{
		TableName: aws.String(drv.tableName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: keyHash},
		},
	})
	if err != nil {
		return 0, fmt.Errorf("dynamodb.GetTTL(): error getting value from dynamodb for key hash %v: %w", keyHash, err)
	}

	if result.Item == nil {
		log.Trace().
			Str("key_hash", keyHash).
			Msg("Found no value in dynamodb")
		return 0, nil
	}

	ttlAttribute, exists := result.Item["TTL"]
	if !exists {
		return 0, nil
	}

	expirationTime, err := strconv.ParseInt(ttlAttribute.(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("dynamodb.GetTTL(): error parsing TTL value for key hash %v: %w", keyHash, err)
	}

	expiresIn := time.Until(time.Unix(expirationTime, 0))
	if expiresIn < 0 {
		expiresIn = 0
	}

	return expiresIn, nil
}
