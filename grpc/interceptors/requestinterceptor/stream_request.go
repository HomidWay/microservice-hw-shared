package requestinterceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ri *RequestInterceptor) streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	requestID := ri.getRequestID(ss.Context())
	if len(requestID) > 0 {
		return handler(srv, ss)
	}

	id := uuid.New()

	ctx := context.WithValue(ss.Context(), "x-request-id", id.String())

	if err := grpc.SetHeader(ctx, metadata.Pairs("x-request-id", id.String())); err != nil {
		return status.Errorf(codes.Internal, "failed to set x-request-id: %v", err)
	}

	return handler(srv, ss)
}
