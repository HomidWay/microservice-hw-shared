package ratelimiterinterceptor

import (
	"context"

	"github.com/HomidWay/microservice-hw-shared/ratelimiter"
	"google.golang.org/grpc"
)

type RateLimitInterceptor struct {
	rateLimiter *ratelimiter.RateLimiter
}

func NewRateLimitInterceptor(ctx context.Context, limitPerSecond, limitPerMinute int) *RateLimitInterceptor {

	rateLimiter := ratelimiter.NewRateLimiter(ctx, limitPerSecond, limitPerMinute)

	return &RateLimitInterceptor{
		rateLimiter: rateLimiter,
	}
}

func (rl *RateLimitInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return rl.unaryInterceptor
}

func (rl *RateLimitInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return rl.streamInterceptor
}

func (rl *RateLimitInterceptor) extractUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}

	return ""
}
