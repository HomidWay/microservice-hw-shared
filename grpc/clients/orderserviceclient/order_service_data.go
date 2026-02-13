package orderserviceclient

import (
	"sync"
	"time"
)

type OrderType int

const (
	OrderTypeUndefined OrderType = iota
	OrderTypeMarketOrder
	OrderTypeLimitOrder
	OrderTypeStopLimitOrder
)

type OrderStatus int

const (
	OrderStatusUndefined OrderStatus = iota
	OrderStatusPendingPayment
	OrderStatusAwaitingConfirmation
	OrderStatusConfirmed
	OrderStatusComplete
	OrderStatusFailed OrderStatus = 99
)

type OrderStatusChangedEvent struct {
	orderID   string
	oldStatus OrderStatus
	newStatus OrderStatus
	timestamp time.Time
}

func NewOrderStatusChangedEvent(orderID string, oldStatus, newStatus OrderStatus, timestamp time.Time) OrderStatusChangedEvent {
	return OrderStatusChangedEvent{
		orderID:   orderID,
		oldStatus: oldStatus,
		newStatus: newStatus,
		timestamp: timestamp,
	}
}

func (e OrderStatusChangedEvent) OrderID() string {
	return e.orderID
}

func (e OrderStatusChangedEvent) OldStatus() OrderStatus {
	return e.oldStatus
}

func (e OrderStatusChangedEvent) NewStatus() OrderStatus {
	return e.newStatus
}

func (e OrderStatusChangedEvent) Timestamp() time.Time {
	return e.timestamp
}

type Order struct {
	orderID     string
	userID      string
	marketID    string
	orderType   OrderType
	createdAt   time.Time
	orderStatus OrderStatus
	price       float64
	quantity    uint

	mu        sync.Mutex
	listeners map[chan OrderStatusChangedEvent]struct{}
}

func NewOrder(orderID, userID, marketID string, orderType OrderType, orderStatus OrderStatus, price float64, quantity uint, createdAt time.Time) *Order {

	order := &Order{
		orderID:     orderID,
		userID:      userID,
		marketID:    marketID,
		orderType:   orderType,
		orderStatus: orderStatus,
		createdAt:   createdAt,
		price:       price,
		quantity:    quantity,
		listeners:   make(map[chan OrderStatusChangedEvent]struct{}),
	}

	return order
}

func (o *Order) OrderID() string {
	return o.orderID
}

func (o *Order) UserID() string {
	return o.userID
}

func (o *Order) MarketID() string {
	return o.marketID
}

func (o *Order) OrderType() OrderType {
	return o.orderType
}

func (o *Order) OrderStatus() OrderStatus {
	return o.orderStatus
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) Price() float64 {
	return o.price
}

func (o *Order) Quantity() uint {
	return o.quantity
}
