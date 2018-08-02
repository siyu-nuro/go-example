package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"github.com/siyu-nuro/go-example/pkg/endpoint"
	"github.com/siyu-nuro/go-example/pkg/model"
	"github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

)

var (
	errBadRequest       = errors.New("bad request")
)

func NewHTTPHandler(handler endpoint.RequestHandler) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(errorEncoder),
	}
	m := mux.NewRouter()

	m.Methods("POST").Path("/orders/create").Handler(httptransport.NewServer(
		handler.CreateOrder,
		decodeHTTPCreateOrderRequest,
		encodeHTTPGenericResponse,
		options...,
	))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	maxAgeOk  := handlers.MaxAge(86400)
	return handlers.CORS(originsOk, headersOk, methodsOk, maxAgeOk)(m)
}

func decodeHTTPCreateOrderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req model.CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return &req, err
}


func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

// TODO: add more error codes
func err2code(err error) int {
	switch err {
	case errBadRequest:
		return http.StatusBadRequest
	case jwt.ErrTokenContextMissing, jwt.ErrUnexpectedSigningMethod, jwt.ErrTokenMalformed, jwt.ErrTokenExpired,
		jwt.ErrTokenNotActive, jwt.ErrTokenInvalid, endpoint.ErrAuthFailed:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
