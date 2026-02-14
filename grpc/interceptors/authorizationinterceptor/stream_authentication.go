package authorizationinterceptor

import (
	"errors"

	"github.com/HomidWay/microservice-hw-shared/sessionvalidation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ai *AuthorizationInterceptor) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	if _, ok := publicMethods[info.FullMethod]; ok {
		return handler(srv, ss)
	}

	sessionID, err := ai.extractSessionID(ss.Context())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid session")
	}

	session, err := ai.sessionValidator.ValidateSession(sessionID)
	if err != nil {
		if errors.Is(err, sessionvalidation.ErrSessionNotFound) {
			return status.Errorf(codes.Unauthenticated, "invalid session")
		}

		return status.Errorf(codes.Internal, "failed to get session data")
	}

	if !session.IsValid() {
		return status.Errorf(codes.Unauthenticated, "invalid session")
	}

	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	md.Set("x-user-id", session.UserID())

	err = ss.SetHeader(md)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to set header")
	}

	wrapped := &authenticatedStream{
		ServerStream:   ss,
		sessionManager: ai.sessionValidator,
		sessionID:      sessionID,
		userID:         session.UserID(),
	}

	return handler(srv, wrapped)
}

type authenticatedStream struct {
	grpc.ServerStream
	sessionManager sessionvalidation.SessionRepository
	sessionID      string
	userID         string
}

func (as *authenticatedStream) RecvMsg(m interface{}) error {

	session, err := as.sessionManager.ValidateSession(as.sessionID)
	if err != nil {
		if errors.Is(err, sessionvalidation.ErrSessionNotFound) {
			return status.Errorf(codes.Unauthenticated, "invalid session")
		}

		return status.Errorf(codes.Internal, "failed to get session data")
	}

	if !session.IsValid() {
		return status.Errorf(codes.Unauthenticated, "invalid session")
	}

	return as.ServerStream.RecvMsg(m)
}

func (as *authenticatedStream) SendMsg(m interface{}) error {

	session, err := as.sessionManager.ValidateSession(as.sessionID)
	if err != nil {
		if errors.Is(err, sessionvalidation.ErrSessionNotFound) {
			return status.Errorf(codes.Unauthenticated, "invalid session")
		}

		return status.Errorf(codes.Internal, "failed to get session data")
	}

	if !session.IsValid() {
		return status.Errorf(codes.Unauthenticated, "invalid session")
	}

	return as.ServerStream.SendMsg(m)
}
