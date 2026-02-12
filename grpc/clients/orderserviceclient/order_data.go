package orderserviceclient

import (
	"fmt"
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
	sessionID   string
	marketID    string
	orderType   OrderType
	createdAt   time.Time
	orderStatus OrderStatus
	price       float64
	quantity    uint

	mu        sync.Mutex
	listeners map[chan OrderStatusChangedEvent]struct{}
}

func NewOrder(orderID, userID, sessionID, marketID string, orderType OrderType, orderStatus OrderStatus, price float64, quantity uint, createdAt time.Time) *Order {

	order := &Order{
		orderID:     orderID,
		userID:      userID,
		sessionID:   sessionID,
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

func (o *Order) SessionID() string {
	return o.sessionID
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

func (o *Order) progress() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.orderStatus == OrderStatusComplete || o.orderStatus == OrderStatusFailed {
		return fmt.Errorf("order %s is already complete", o.orderID)
	}

	newStatus := OrderStatus(int(o.OrderStatus()) + 1)

	event := NewOrderStatusChangedEvent(o.OrderID(), o.OrderStatus(), newStatus, time.Now())

	o.orderStatus = newStatus

	o.notifySubscribers(event)

	return nil
}

func (o *Order) SubscribeToUpdates() chan OrderStatusChangedEvent {
	o.mu.Lock()
	defer o.mu.Unlock()

	eventChan := make(chan OrderStatusChangedEvent)
	o.listeners[eventChan] = struct{}{}

	return eventChan
}

func (o *Order) UnsubscribeFromUpdates(c chan OrderStatusChangedEvent) {
	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.listeners, c)
}

func (o *Order) notifySubscribers(event OrderStatusChangedEvent) {
	for listener := range o.listeners {
		select {
		case listener <- event:
		default:
		}
	}
}
