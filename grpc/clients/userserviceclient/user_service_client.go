package userserviceclient

import (
	"context"
	"fmt"

	"github.com/HomidWay/microservice-hw-proto/pb/userrole"
	"github.com/HomidWay/microservice-hw-proto/pb/userservice"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	rpcClient userservice.UserServiceClient
}

func NewUserServiceClient(host string, port int, options ...grpc.DialOption) (*UserServiceClient, error) {

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), options...)
	if err != nil {
		return nil, err
	}

	client := userservice.NewUserServiceClient(conn)

	return &UserServiceClient{rpcClient: client}, nil
}

func (c *UserServiceClient) RegisterNewUser(ctx context.Context, userName, password, passwordConfirm, fullName string, userRole UserRole) (*User, error) {
	request := &userservice.RegisterNewUserRequest{
		Username:        userName,
		Password:        password,
		PasswordConfirm: passwordConfirm,
		FullName:        fullName,
		UserRole:        userrole.UserRole(userRole),
	}

	response, err := c.rpcClient.RegisterNewUser(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to register new user: %w", err)
	}

	return NewUser(
		response.UserId,
		response.Username,
		response.FullName,
		UserRole(response.UserRole),
	), nil
}

func (c *UserServiceClient) GetUserData(ctx context.Context, userID string) (*User, error) {

	identifier := &userservice.GetUserDataRequest_UserId{UserId: userID}

	request := &userservice.GetUserDataRequest{
		Identification: identifier,
	}

	response, err := c.rpcClient.GetUserData(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data with error: %w", err)
	}

	return NewUser(
		response.UserId,
		response.Username,
		response.FullName,
		UserRole(response.UserRole),
	), nil
}

func (c *UserServiceClient) ChangeUserRole(ctx context.Context, userID string, newRole UserRole) (*User, error) {

	identifier := &userservice.ChangeUserRoleRequest_Userid{Userid: userID}

	request := &userservice.ChangeUserRoleRequest{
		Identification: identifier,
		NewRole:        userrole.UserRole(newRole),
	}

	response, err := c.rpcClient.ChangeUserRole(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data with error: %w", err)
	}

	return NewUser(
		response.UserId,
		response.Username,
		response.FullName,
		UserRole(response.UserRole),
	), nil
}
