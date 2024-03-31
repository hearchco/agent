package router

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
)

func (rw RouterWrapper) StartLambda() {
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return chiadapter.NewV2(rw.mux).ProxyWithContextV2(ctx, req)
	})
}
