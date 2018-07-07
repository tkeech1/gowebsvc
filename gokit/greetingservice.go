package main

import (
	"context"
	"log"
	"sync"

	"encoding/json"
	"net/http"

	service "github.com/tkeech1/gowebsvc/svc"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// https://gokit.io/examples/stringsvc.html

// endpoints
func makeGreetingEndpoint(svc service.GreetingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.GreetRequest)
		v, err := svc.Greet(ctx, req.S)
		if err != nil {
			return service.GreetResponse{"", err.Error()}, nil
		}
		return service.GreetResponse{v, ""}, nil
	}
}

func makeExpensiveEndpoint(svc service.GreetingService) endpoint.Endpoint {
	var (
		init sync.Once
	)
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		v := "already initialized"
		init.Do(func() {
			// do an expensive operation here - it will only occur on the first invocation of the handler
			// initialize database, parse templates, etc.
			req := request.(service.ExpensiveRequest)
			log.Print("middleware - just do this once")
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

func getGreetingHandler(svc service.GreetingService) *httptransport.Server {
	return httptransport.NewServer(
		makeGreetingEndpoint(svc),
		decodeGreetRequest,
		encodeResponse,
	)
}

func getExpensiveHandler(svc service.GreetingService) *httptransport.Server {
	return httptransport.NewServer(
		makeExpensiveEndpoint(svc),
		decodeExpensiveRequest,
		encodeResponse,
	)
}

// main
func main() {
	svc := service.GreetingService{}

	greetingHandler := getGreetingHandler(svc)
	expensiveHandler := getExpensiveHandler(svc)

	http.Handle("/greeting", greetingHandler)
	http.Handle("/expensive", expensiveHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
