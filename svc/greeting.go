package svc

import (
	"context"
	"errors"
	"time"
)

type Greeter interface {
	Greet(context.Context, string) (string, error)
	Expensive(context.Context, string, string, string) (string, error)
}

type GreeterGRPC interface {
	GreetGRPC(context.Context, *GRPCGreetRequest) (*GRPCGreetResponse, error)
}

type GreetingService struct{}

type GreetingServiceGRPC struct{}

func (g *GreetingServiceGRPC) GreetGRPC(ctx context.Context, in *GRPCGreetRequest) (*GRPCGreetResponse, error) {
	return &GRPCGreetResponse{Greeting: "GRPC - " + in.S}, nil
}

func (g GreetingService) Greet(ctx context.Context, greeting string) (string, error) {
	ch := make(chan string)

	// in case this is a long-runnning operation, use a go function
	go func(s string) {
		ch <- s
	}(greeting)

	select {
	case response := <-ch:
		if response == "" {
			return "", errors.New("empty greeting")
		}
		return response, nil
	case <-ctx.Done():
		return "", errors.New("request cancelled")
	}
}

func (g GreetingService) Expensive(ctx context.Context, connectionString, username, password string) (string, error) {
	if connectionString == "" {
		return "", errors.New("missing connectionString")
	}
	if username == "" {
		return "", errors.New("missing username")
	}
	if password == "" {
		return "", errors.New("missing password")
	}

	ch := make(chan string)
	go func(c, u, p string) {
		ch <- connectionString + username + password
	}(connectionString, username, password)

	select {
	case response := <-ch:
		if response == "" {
			return "", errors.New("empty greeting")
		}
		return response, nil
	case <-ctx.Done():
		return "", errors.New("request cancelled")
	case <-time.After(1 * time.Second):
		return "", errors.New("request timed out")
	}
}
