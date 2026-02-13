package cachinginterceptor

import (
	"github.com/HomidWay/microservice-hw-proto/pb/orderservice"
	"github.com/HomidWay/microservice-hw-proto/pb/spotinstrumentservice"
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

	spotinstrumentservice.SpotInstrumentService_ViewMarkets_FullMethodName: func() proto.Message {
		return &spotinstrumentservice.ViewMarketsResponse{}
	},

	userservice.UserService_GetUserData_FullMethodName: func() proto.Message {
		return &userservice.GetUserDataResponse{}
	},
}
