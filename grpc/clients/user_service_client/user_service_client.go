package userserviceclient

import (
	"context"
	"fmt"
	"time"

	"github.com/HomidWay/microservice-hw-proto/pb/userservice"
	sessionvalidation "github.com/HomidWay/microservice-hw-shared/session_validation"
	usermanagement "github.com/HomidWay/microservice-hw-shared/user_management"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	rpcClient userservice.UserServiceClient
}

func NewUserServiceClient(host string, port int) (*UserServiceClient, error) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to userservice: %w", err)
	}

	client := userservice.NewUserServiceClient(conn)

	return &UserServiceClient{rpcClient: client}, nil
}

func (u UserServiceClient) Login(username, password string) (*sessionvalidation.Session, error) {

	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := u.rpcClient.AuthorizeUser(ctx, &userservice.AuthorizeUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to authorize user: %w", err)
	}

	user := sessionvalidation.NewUser(
		response.UserId,
		response.Username,
		response.FullName,
		usermanagement.UserRole(response.UserRole),
	)

	session := sessionmanager.NewSession(user, response.SessionId)

	return session, nil
}

func (u UserServiceClient) RegisterNewUser(username, password, fullname string, role sessionmanager.UserRole) (*sessionmanager.Session, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := u.rpcClient.RegisterNewUser(ctx, &userservice.RegisterNewUserRequest{
		Username: username,
		Password: password,
		FullName: fullname,
		UserRole: usermanagement.UserRole(role),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register new user: %w", err)
	}

	return u.Login(username, password)
}

func (u UserServiceClient) LoadSessionData(sessionID, userID string) (*sessionmanager.Session, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	response, err := u.rpcClient.ValidateSession(ctx, &userservice.ValidateSessionRequest{
		SessionId:      sessionID,
		Identification: &userservice.ValidateSessionRequest_UserId{UserId: userID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	userData, err := u.rpcClient.GetUserData(ctx, &userservice.GetUserDataRequest{
		Identification: &userservice.GetUserDataRequest_UserId{UserId: userID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	user := sessionmanager.NewUser(
		userData.UserId,
		userData.Username,
		userData.FullName,
		sessionmanager.UserRole(userData.UserRole),
	)

	session := sessionmanager.NewSession(user, response.SessionId)

	return session, nil
}
