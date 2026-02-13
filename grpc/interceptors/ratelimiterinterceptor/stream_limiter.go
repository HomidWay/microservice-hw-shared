package ratelimiterinterceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (rl *RateLimitInterceptor) streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	userID := rl.extractUserID(ss.Context())

	if userID == "" {
		return handler(srv, ss)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := rl.rateLimiter.WaitForQuota(ctxTimeout, userID)
	if err != nil {
		return status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}

	return handler(srv, ss)
}
