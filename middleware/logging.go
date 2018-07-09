package middleware

import (
	"context"
	"log"
	"time"

	service "github.com/tkeech1/gowebsvc/svc"
)

type LoggingMiddleware struct {
	Logger *log.Logger
	Next   service.Greeter
}

type LoggingMiddlewareGRPC struct {
	Logger *log.Logger
	Next   service.GreeterGRPC
}

func (mw LoggingMiddleware) Greet(ctx context.Context, greeting string) (output string, err error) {
	defer func(begin time.Time) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		mw.Logger.Print(
			"method: ", "Greet"+"; ",
			"input: ", greeting+"; ",
			"output: ", output+"; ",
			"err: ", errMsg+"; ",
			"took: ", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Next.Greet(ctx, greeting)
	return
}

func (mw LoggingMiddleware) Expensive(ctx context.Context, connectionString, username, password string) (n string, err error) {
	defer func(begin time.Time) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		mw.Logger.Print(
			"method: ", "Expensive"+"; ",
			"connectionString: ", connectionString+"; ",
			"username: ", username+"; ",
			"password: ", password+"; ",
			"output: ", n+"; ",
			"err: ", errMsg+"; ",
			"took: ", time.Since(begin),
		)
	}(time.Now())

	n, err = mw.Next.Expensive(ctx, connectionString, username, password)
	return
}

func (mw LoggingMiddlewareGRPC) GreetGRPC(ctx context.Context, in *service.GRPCGreetRequest) (output *service.GRPCGreetResponse, err error) {
	defer func(begin time.Time) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		mw.Logger.Print(
			"method: ", "Greet"+"; ",
			"input: ", in.S+"; ",
			"output: ", output.Greeting+"; ",
			"err: ", errMsg+"; ",
			"took: ", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Next.GreetGRPC(ctx, in)
	return
}
