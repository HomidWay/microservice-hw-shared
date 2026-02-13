package orderserviceclient

import (
	"context"
	"fmt"

	"github.com/HomidWay/microservice-hw-proto/pb/orderservice"
	"google.golang.org/grpc"
)

type OrderServiceClient struct {
	rpcClient orderservice.OrderServiceClient
}

func NewOrderServiceClient(host string, port int, options ...grpc.DialOption) (*OrderServiceClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to userservice: %w", err)
	}

	client := orderservice.NewOrderServiceClient(conn)

	return &OrderServiceClient{rpcClient: client}, nil
}

func (o OrderServiceClient) CreateOrder(ctx context.Context, userID, marketID string, orderType OrderType, price float64, quantity uint) (*Order, error) {

	request := &orderservice.CreateOrderRequest{
		UserId:    userID,
		MarketId:  marketID,
		OrderType: orderservice.OrderType(orderType),
		Price:     float32(price),
		Quantity:  uint64(quantity),
	}

	response, err := o.rpcClient.CreateOrder(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return NewOrder(
		response.OrderId,
		response.UserId,
		response.MarketId,
		OrderType(response.OrderType),
		OrderStatus(response.Status),
		float64(response.Price),
		uint(response.Quantity),
		response.CreatedAt.AsTime(),
	), nil
}

func (o OrderServiceClient) GetOrderStatus(ctx context.Context, orderID string) (*Order, error) {

	request := &orderservice.GetOrderStatusRequest{
		OrderId: orderID,
	}

	response, err := o.rpcClient.GetOrderStatus(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return NewOrder(
		response.OrderId,
		response.UserId,
		response.MarketId,
		OrderType(response.OrderType),
		OrderStatus(response.Status),
		float64(response.Price),
		uint(response.Quantity),
		response.CreatedAt.AsTime(),
	), nil
}

func (o OrderServiceClient) GetAllOrderIDs(ctx context.Context, userID string) ([]string, error) {

	request := &orderservice.GetAllOrderIDsRequest{
		UserId: userID,
	}

	response, err := o.rpcClient.GetAllOrderIDs(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return response.OrderIDs, nil
}

func (o OrderServiceClient) SubscribeToOrderUpdates(ctx context.Context, orderID string) (chan OrderStatusChangedEvent, error) {

	responseStream, err := o.rpcClient.StreamOrderUpdates(ctx, &orderservice.StreamOrderUpdatesRequest{
		OrderId: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	responseChan := make(chan OrderStatusChangedEvent)

	go func() {
		for {

			eventData, err := responseStream.Recv()
			if err != nil {
				fmt.Println("failed to receive order update event: ", err.Error())
				close(responseChan)
				return
			}

			event := NewOrderStatusChangedEvent(
				eventData.OrderId,
				OrderStatus(eventData.OldStatus),
				OrderStatus(eventData.NewStatus),
				eventData.Timestamp.AsTime(),
			)

			select {
			case responseChan <- event:
			default:
			}
		}
	}()

	return responseChan, nil
}
