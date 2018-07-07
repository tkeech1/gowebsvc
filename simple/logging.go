package main

import (
	"context"
	"log"
	"time"

	service "github.com/tkeech1/gowebsvc/svc"
)

type loggingMiddleware struct {
	logger *log.Logger
	next   service.Greeter
}

func (mw loggingMiddleware) Greet(ctx context.Context, greeting string) (output string, err error) {
	defer func(begin time.Time) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		mw.logger.Print(
			"method: ", "Greet"+"; ",
			"input: ", greeting+"; ",
			"output: ", output+"; ",
			"err: ", errMsg+"; ",
			"took: ", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Greet(ctx, greeting)
	return
}

func (mw loggingMiddleware) Expensive(ctx context.Context, connectionString, username, password string) (n string, err error) {
	defer func(begin time.Time) {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		mw.logger.Print(
			"method: ", "Expensive"+"; ",
			"connectionString: ", connectionString+"; ",
			"username: ", username+"; ",
			"password: ", password+"; ",
			"output: ", n+"; ",
			"err: ", errMsg+"; ",
			"took: ", time.Since(begin),
		)
	}(time.Now())

	n, err = mw.next.Expensive(ctx, connectionString, username, password)
	return
}
