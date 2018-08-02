package endpoint

import (
	"context"

	"github.com/siyu-nuro/go-example/pkg/mapper/json"
	"github.com/siyu-nuro/go-example/pkg/model"
	"github.com/siyu-nuro/go-example/pkg/service"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

// RequestHandler handles all incoming requests with corresponding services
type RequestHandler struct {
	CreateOrder               endpoint.Endpoint
}

// NewRequestHandler constructs a new Handler
func NewRequestHandler(o service.OrderService) RequestHandler {
	requestHandler := RequestHandler{}
	{
		var createOrder endpoint.Endpoint
		createOrder = MakeCreateOrder(o)
		createOrder = jwt.NewParser(validateRetailerAuth, stdjwt.SigningMethodRS256, jwt.MapClaimsFactory)(createOrder)
		requestHandler.CreateOrder = createOrder
	}

	return requestHandler
}

func MakeCreateOrder(o service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*model.CreateOrderRequest)
		if !req.HasRequiredFields() {
			return nil, model.ErrCreateOrderMissingRequiredFields
		}
		orderDetail := json.MapModelOrderDetailToInputOrderDetail(req.OrderDetail)
		orderID, err := o.CreateOrder(ctx, orderDetail)
		return json.MapCreateOrderResponse(orderID, err), nil
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}
