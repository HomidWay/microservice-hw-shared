package requestinterceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RequestInterceptor struct{}

func NewRequestInterceptor() *RequestInterceptor {
	return &RequestInterceptor{}
}

func (ri RequestInterceptor) Intercept(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	id := uuid.New()

	ctx = context.WithValue(ctx, "request_id", id.String())

	if err := grpc.SetHeader(ctx, metadata.Pairs("x-request-id", id.String())); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set x-request-id: %v", err)
	}

	return handler(ctx, req)
}
