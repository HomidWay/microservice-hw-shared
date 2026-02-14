package requestinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type RequestInterceptor struct{}

func NewRequestInterceptor() *RequestInterceptor {
	return &RequestInterceptor{}
}

func (ri *RequestInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return ri.unaryInterceptor
}

func (ri *RequestInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return ri.streamInterceptor
}

func (ri *RequestInterceptor) getRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		requestID := md.Get("x-request-id")
		if len(requestID) > 0 {
			return requestID[0]
		}
	}

	return ""
}
