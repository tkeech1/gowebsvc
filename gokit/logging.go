package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	service "github.com/tkeech1/gowebsvc/svc"
)

type loggingMiddleware struct {
	logger log.Logger
	next   service.Greeter
}

func (mw loggingMiddleware) Greet(ctx context.Context, greeting string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Greet",
			"input", greeting,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Greet(ctx, greeting)
	return
}

func (mw loggingMiddleware) Expensive(ctx context.Context, connectionString, username, password string) (n string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Expensive",
			"connectionString", connectionString,
			"username", username,
			"password", password,
			"output", n,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	n, err = mw.next.Expensive(ctx, connectionString, username, password)
	return
}
