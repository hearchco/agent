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
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(conf.Region))
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
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	key := anonymize.HashToSHA256B64(fmt.Sprintf("%v%v", drv.keyPrefix, k))
	item, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("dynamodb.Set(): error marshaling value: %w", err)
	}

	attributes := map[string]types.AttributeValue{
		"Key":   &types.AttributeValueMemberS{Value: key},
		"Value": &types.AttributeValueMemberS{Value: string(item)},
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

	log.Trace().Dur("duration", time.Since(cacheTimer)).Msg("Cached results")
	return nil
}

func (drv DRV) Get(k string, o any) error {
	key := anonymize.HashToSHA256B64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	result, err := drv.client.GetItem(drv.ctx, &dynamodb.GetItemInput{
		TableName: aws.String(drv.tableName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb.Get(): error getting value from dynamodb for key %v: %w", key, err)
	}

	if result.Item == nil {
		log.Trace().
			Str("key", key).
			Msg("Found no value in dynamodb")
		return nil
	}

	// Required because TTL isn't guaranteed to remove the item immediately
	if result.Item["TTL"] != nil {
		expirationTime, err := strconv.ParseInt(result.Item["TTL"].(*types.AttributeValueMemberN).Value, 10, 64)
		if err != nil {
			return fmt.Errorf("dynamodb.Get(): error parsing TTL value for key %v: %w", key, err)
		}

		expiresIn := time.Until(time.Unix(expirationTime, 0))
		if expiresIn < 0 {
			log.Trace().
				Str("key", key).
				Msg("Value has expired")
			return nil
		}
	}

	value := result.Item["Value"].(*types.AttributeValueMemberS).Value
	if err := json.Unmarshal([]byte(value), o); err != nil {
		return fmt.Errorf("dynamodb.Get(): error unmarshaling value from dynamodb for key %v: %w", key, err)
	}

	return nil
}

func (drv DRV) GetTTL(k string) (time.Duration, error) {
	key := anonymize.HashToSHA256B64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	result, err := drv.client.GetItem(drv.ctx, &dynamodb.GetItemInput{
		TableName: aws.String(drv.tableName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return 0, fmt.Errorf("dynamodb.GetTTL(): error getting value from dynamodb for key %v: %w", key, err)
	}

	if result.Item == nil {
		log.Trace().Str("key", key).Msg("Found no value in dynamodb")
		return 0, nil
	}

	ttlAttribute, exists := result.Item["TTL"]
	if !exists {
		return 0, nil
	}

	expirationTime, err := strconv.ParseInt(ttlAttribute.(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("dynamodb.GetTTL(): error parsing TTL value for key %v: %w", key, err)
	}

	expiresIn := time.Until(time.Unix(expirationTime, 0))
	if expiresIn < 0 {
		expiresIn = 0
	}

	return expiresIn, nil
}
