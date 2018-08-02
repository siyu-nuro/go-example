package order

import (
	"context"

	"github.com/siyu-nuro/go-example/pkg/entity"
	"github.com/siyu-nuro/go-example/pkg/service"
	"github.com/siyu-nuro/go-example/pkg/store"

	"go.uber.org/zap"
)

type orderService struct {
	logger              *zap.SugaredLogger
	storageManager      store.StorageManager
}

// NewOrderService constructs a new OrderService
func NewOrderService(l *zap.SugaredLogger, s store.StorageManager) service.OrderService {
	return &orderService{
		logger:              l,
		storageManager:      s,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, orderDetail *entity.InputOrderDetail) (*entity.OrderID, error) {
	// insert into db order table
	orderID, err := o.storageManager.CreateOrder(ctx, orderDetail)
	if err != nil {
		o.logger.Errorw("failed to insert order into db",
			"error", err.Error())
		return nil, err
	}

	return orderID, nil
}
