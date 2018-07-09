package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	service "github.com/tkeech1/gowebsvc/svc"

	"github.com/go-kit/kit/endpoint"
)

// https://gokit.io/examples/stringsvc.html

// endpoints
func makeGreetingEndpoint(svc service.Greeter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.GreetRequest)
		v, err := svc.Greet(ctx, req.S)
		if err != nil {
			return service.GreetResponse{"", err.Error()}, nil
		}
		return service.GreetResponse{v, ""}, nil
	}
}

func makeExpensiveEndpoint(svc service.Greeter) endpoint.Endpoint {
	var (
		init sync.Once
	)
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		v := "already initialized"
		init.Do(func() {
			// do an expensive operation here - it will only occur on the first invocation of the handler
			req := request.(service.ExpensiveRequest)
			v, err = svc.Expensive(ctx, req.C, req.U, req.P)
		})

		if err != nil {
			return service.ExpensiveResponse{"", err.Error()}, nil
		}
		return service.ExpensiveResponse{v, ""}, nil
	}
}

// transports
func decodeGreetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request service.GreetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeExpensiveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request service.ExpensiveRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
