package orderserviceclient

import (
	"context"
	"fmt"
	"time"

	"github.com/HomidWay/microservice-hw-proto/pb/orderservice"
	"google.golang.org/grpc"
)

type OrderServiceClient struct {
	rpcClient orderservice.OrderServiceClient
}

func NewOrderServiceClient(host string, port int) (*OrderServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to userservice: %w", err)
	}

	client := orderservice.NewOrderServiceClient(conn)

	return &OrderServiceClient{rpcClient: client}, nil
}

func (o OrderServiceClient) CreateOrder(ctx context.Context, userID, sessionID, marketID string, orderType OrderType, price float64, quantity uint) (*Order, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := o.rpcClient.CreateOrder(ctx, &orderservice.CreateOrderRequest{
		UserId:    userID,
		SessionId: sessionID,
		MarketId:  marketID,
		OrderType: orderservice.OrderType(orderType),
		Price:     float32(price),
		Quantity:  uint64(quantity),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	order := NewOrder(
		response.OrderId,
		response.UserId,
		sessionID,
		response.MarketId,
		OrderType(response.OrderType),
		OrderStatus(response.Status),
		float64(response.Price),
		uint(response.Quantity),
		response.CreatedAt.AsTime(),
	)

	return order, nil
}

func (o OrderServiceClient) GetOrderInfo(userID, orderID string) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := o.rpcClient.GetOrderStatus(ctx, &orderservice.GetOrderStatusRequest{
		UserId:  userID,
		OrderId: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order := NewOrder(
		response.OrderId,
		response.UserId,
		response.SessionId,
		response.MarketId,
		OrderType(response.OrderType),
		OrderStatus(response.Status),
		float64(response.Price),
		uint(response.Quantity),
		response.CreatedAt.AsTime(),
	)

	return order, nil
}

func (o OrderServiceClient) GetAllOrderIDs(userID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := o.rpcClient.GetAllOrderIDs(ctx, &orderservice.GetAllOrderIDsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return response.OrderIDs, nil
}

func (o OrderServiceClient) SubscribeToOrderUpdates(ctx context.Context, userID, orderID string) (chan OrderStatusChangedEvent, error) {

	responseStream, err := o.rpcClient.StreamOrderUpdates(ctx, &orderservice.StreamOrderUpdatesRequest{
		UserId:  userID,
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
