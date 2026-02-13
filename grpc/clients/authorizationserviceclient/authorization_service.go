package authorizationserviceclient

import (
	"context"
	"fmt"

	"github.com/HomidWay/microservice-hw-proto/pb/authorizationservice"
	"google.golang.org/grpc"
)

type AuthenticationServiceClient struct {
	rpcClient authorizationservice.AuthorizationServiceClient
}

func NewAuthenticationServiceClient(host string, port int, options ...grpc.DialOption) (*AuthenticationServiceClient, error) {

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), options...)
	if err != nil {
		return nil, err
	}

	client := authorizationservice.NewAuthorizationServiceClient(conn)

	return &AuthenticationServiceClient{rpcClient: client}, nil
}

func (c *AuthenticationServiceClient) ValidateSession(ctx context.Context, sessionID string) (*Session, error) {

	request := &authorizationservice.ValidateSessionRequest{
		SessionId: sessionID,
	}

	sessionData, err := c.rpcClient.ValidateSession(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("validate session fail: %w", err)
	}

	session := NewSession(
		sessionData.SessionId,
		sessionData.UserId,
		SessionState(sessionData.State),
		sessionData.ValidUntil.AsTime(),
	)

	return session, nil
}
