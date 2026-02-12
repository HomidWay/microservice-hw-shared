package authorizationserviceclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/HomidWay/microservice-hw-proto/pb/authorizationservice"
	"google.golang.org/grpc"
)

var (
	ErrEmptyCredentials   = errors.New("credentials are empty")
	ErrInvalidCredentials = errors.New("username or Password is invalid")
)

type AuthorizationServiceClient struct {
	rpcClient authorizationservice.AuthorizationServiceClient
}

func NewAuthorizationServiceClient(host string, port int) (*AuthorizationServiceClient, error) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to authserver: %w", err)
	}

	client := authorizationservice.NewAuthorizationServiceClient(conn)

	return &AuthorizationServiceClient{rpcClient: client}, nil
}

func (c *AuthorizationServiceClient) AuthorizeUser(ctx context.Context, userName, password string) (*AuthorizationResponse, error) {

	if userName == "" || password == "" {
		return nil, ErrEmptyCredentials
	}

	response, err := c.rpcClient.AuthorizeUser(context.Background(), &authorizationservice.AuthorizeUserRequest{Username: userName, Password: password})
	if err != nil {
		return nil, fmt.Errorf("Failed to authorize user with error: %w", err)
	}

	return NewAuthorizationResponse(
		response.SessionId,
		response.UserId,
	), nil
}

func (c *AuthorizationServiceClient) ValidateSession(ctx context.Context, sessionID string) (*ValidateSessionResponse, error) {

	if sessionID == "" {
		return nil, ErrEmptyCredentials
	}

	response, err := c.rpcClient.ValidateSession(ctx, &authorizationservice.ValidateSessionRequest{SessionId: sessionID})
	if err != nil {
		return nil, fmt.Errorf("Failed to validate session with error: %w", err)
	}

	return NewValidateSessionResponse(
		response.SessionId,
		response.UserId,
		SessionState(response.State),
		response.ValidUntil.AsTime(),
	), nil
}
