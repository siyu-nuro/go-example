package service

import (
	"context"

	"github.com/siyu-nuro/go-example/pkg/entity"
)

// OrderService manages all order-related logic
type OrderService interface {
	CreateOrder(ctx context.Context, orderDetail *entity.InputOrderDetail) (*entity.OrderID, error)
}
