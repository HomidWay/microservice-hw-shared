package ratelimiterinterceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (rl *RateLimitInterceptor) unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	userID := rl.extractUserID(ctx)

	if userID == "" {
		return handler(ctx, req)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = rl.rateLimiter.WaitForQuota(ctxTimeout, userID)
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}

	return handler(ctx, req)
}
