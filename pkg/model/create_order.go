package model

import "errors"

var (
	ErrCreateOrderMissingRequiredFields = errors.New("createOrderRequest doesn't have all required fields")
)

type CreateOrderRequest struct {
	OrderDetail *OrderDetail `json:"orderDetail"`
}

type CreateOrderResponse struct {
	OrderID *string `json:"orderId"`
	Error   error   `json:"-"`
}

func (c *CreateOrderRequest) HasRequiredFields() bool {
	return c.OrderDetail != nil && c.OrderDetail.HasRequiredFields()
}

// Failed implements Failer.
func (r CreateOrderResponse) Failed() error {
	return r.Error
}
