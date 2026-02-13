package requestinterceptor

import (
	"github.com/HomidWay/microservice-hw-proto/pb/authorizationservice"
	"github.com/HomidWay/microservice-hw-proto/pb/userservice"
)

var publicMethods = map[string]struct{}{
	authorizationservice.AuthorizationService_AuthorizeUser_FullMethodName:   struct{}{},
	authorizationservice.AuthorizationService_ValidateSession_FullMethodName: struct{}{},
	userservice.UserService_RegisterNewUser_FullMethodName:                   struct{}{},
}
