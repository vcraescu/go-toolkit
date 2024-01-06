package echo

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

type LambdaFunc func(
	ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

func NewLambda(srv *Server) LambdaFunc {
	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return echoadapter.NewV2(srv.Echo).ProxyWithContext(ctx, req)
	}
}
