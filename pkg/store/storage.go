package store

import (
	"context"
	"github.com/siyu-nuro/go-example/pkg/entity"
)

// StorageManager defines the data access layer
type StorageManager interface {
	CreateOrder(ctx context.Context, orderDetail *entity.InputOrderDetail) (*entity.OrderID, error)
}
