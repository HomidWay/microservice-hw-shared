package requestinterceptor

import (
	"context"
	"errors"

	sessionvalidation "github.com/HomidWay/microservice-hw-shared/session_validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RequestInterceptor struct {
	sessionManager sessionvalidation.SessionRepository
}

func NewRequestInterceptor() *RequestInterceptor {
	return &RequestInterceptor{}
}

func (ri RequestInterceptor) Intercept(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	sessionIDVal := ctx.Value("Authorization")
	if sessionIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "session is empty")
	}

	sessionID, ok := sessionIDVal.(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "session is empty")
	}

	session, err := ri.sessionManager.ValidateSession(sessionID)
	if err != nil {
		if errors.Is(err, sessionvalidation.ErrSessionNotFound) {
			return nil, status.Error(codes.Unauthenticated, "session is invalid")
		}

		return nil, status.Error(codes.Internal, "failed to validate session")
	}

	if !session.IsValid() {
		return nil, status.Error(codes.Unauthenticated, "session is invalid")
	}

	ctx = context.WithValue(ctx, "x-user-id", session.UserID())

	if err := grpc.SetHeader(ctx, metadata.Pairs("x-user-id", session.UserID())); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set x-user-id: %v", err)
	}

	return handler(ctx, req)
}
