package requestinterceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ri *RequestInterceptor) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	requestID := ri.getRequestID(ctx)
	if len(requestID) > 0 {
		return handler(ctx, req)
	}

	id := uuid.New()

	ctx = context.WithValue(ctx, "x-request-id", id.String())

	if err := grpc.SetHeader(ctx, metadata.Pairs("x-request-id", id.String())); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set x-request-id: %v", err)
	}

	return handler(ctx, req)
}
