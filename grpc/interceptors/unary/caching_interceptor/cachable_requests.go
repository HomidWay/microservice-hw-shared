package cachinginterceptor

import (
	orderservice "github.com/HomidWay/micro-service-hw-shared/grpc/pb/order_service"
	spotintrumentservice "github.com/HomidWay/micro-service-hw-shared/grpc/pb/spot_instrument_service"
	userservice "github.com/HomidWay/micro-service-hw-shared/grpc/pb/user_service"
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
