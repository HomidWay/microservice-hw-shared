package authorizationinterceptor

import (
	"context"
	"fmt"

	"github.com/HomidWay/microservice-hw-shared/sessionvalidation"
	"google.golang.org/grpc"
)

type AuthorizationInterceptor struct {
	sessionValidator sessionvalidation.SessionRepository
}

func NewAuthorizationInterceptor(sessionValidator sessionvalidation.SessionRepository) *AuthorizationInterceptor {
	return &AuthorizationInterceptor{
		sessionValidator: sessionValidator,
	}
}

func (ai *AuthorizationInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return ai.unaryInterceptor
}

func (ai *AuthorizationInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return ai.streamInterceptor
}

func (ai *AuthorizationInterceptor) extractSessionID(ctx context.Context) (string, error) {
	if sessionID, ok := ctx.Value("Authentication").(string); ok {
		return sessionID, nil
	}

	return "", fmt.Errorf("no session id found")
}
