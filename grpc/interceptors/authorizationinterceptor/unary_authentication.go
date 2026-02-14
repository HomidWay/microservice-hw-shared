package authorizationinterceptor

import (
	"context"
	"errors"

	"github.com/HomidWay/microservice-hw-shared/sessionvalidation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ai *AuthorizationInterceptor) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if _, ok := publicMethods[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	sessionIDVal := ctx.Value("Authorization")
	if sessionIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "session is empty")
	}

	sessionID, ok := sessionIDVal.(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "session is empty")
	}

	session, err := ai.sessionManager.ValidateSession(sessionID)
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
