package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	service "github.com/tkeech1/gowebsvc/svc"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           service.Greeter
}

func (mw instrumentingMiddleware) Greet(ctx context.Context, greeting string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "greeting", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.Greet(ctx, greeting)
	return
}

func (mw instrumentingMiddleware) Expensive(ctx context.Context, connectionString, username, password string) (n string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "expensive", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	n, err = mw.next.Expensive(ctx, connectionString, username, password)
	return
}
