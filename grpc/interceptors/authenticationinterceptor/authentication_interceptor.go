package requestinterceptor

import (
	"context"
	"fmt"

	"github.com/HomidWay/microservice-hw-shared/sessionvalidation"
	"google.golang.org/grpc"
)

type AuthenticationInterceptor struct {
	sessionManager sessionvalidation.SessionRepository
}

func NewAuthenticationInterceptor() *AuthenticationInterceptor {
	return &AuthenticationInterceptor{}
}

func (ai *AuthenticationInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return ai.unaryInterceptor
}

func (ai *AuthenticationInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return ai.streamInterceptor
}

func (ai *AuthenticationInterceptor) extractSessionID(ctx context.Context) (string, error) {
	if sessionID, ok := ctx.Value("Authentication").(string); ok {
		return sessionID, nil
	}

	return "", fmt.Errorf("no session id found")
}
