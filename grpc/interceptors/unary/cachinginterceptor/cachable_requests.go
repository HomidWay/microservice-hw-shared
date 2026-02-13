package cachinginterceptor

import (
	"github.com/HomidWay/microservice-hw-proto/pb/orderservice"
	"github.com/HomidWay/microservice-hw-proto/pb/spotintrumentservice"
	"github.com/HomidWay/microservice-hw-proto/pb/userservice"
	"google.golang.org/protobuf/proto"
)

var cachableRequests = map[string]ResponseFactory{

	orderservice.OrderService_GetAllOrderIDs_FullMethodName: func() proto.Message {
		return &orderservice.GetAllOrderIDsResponse{}
	},

	orderservice.OrderService_GetOrderStatus_FullMethodName: func() proto.Message {
		return &orderservice.GetOrderStatusResponse{}
	},

	spotintrumentservice.SpotInstrumentService_ViewMarkets_FullMethodName: func() proto.Message {
		return &spotintrumentservice.ViewMarketsResponse{}
	},

	userservice.UserService_GetUserData_FullMethodName: func() proto.Message {
		return &userservice.GetUserDataResponse{}
	},
}
