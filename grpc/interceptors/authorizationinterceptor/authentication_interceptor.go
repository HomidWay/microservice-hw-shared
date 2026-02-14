package authorizationinterceptor

import (
	"context"
	"fmt"

	"github.com/HomidWay/microservice-hw-shared/sessionvalidation"
	"google.golang.org/grpc"
)

type AuthorizationInterceptor struct {
	sessionManager sessionvalidation.SessionRepository
}

func NewAuthenticationInterceptor() *AuthorizationInterceptor {
	return &AuthorizationInterceptor{}
}

func (ai *AuthorizationInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return ai.unaryInterceptor
}

func (ai *AuthorizationInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return ai.streamInterceptor
}

func (ai *AuthorizationInterceptor) extractSessionID(ctx context.Context) (string, error) {
	if sessionID, ok := ctx.Value("Authentication").(string); ok {
		return sessionID, nil
	}

	return "", fmt.Errorf("no session id found")
}
