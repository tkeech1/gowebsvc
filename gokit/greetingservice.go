package main

import (
	"log"
	"os"

	"net/http"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	middleware "github.com/tkeech1/gowebsvc/middleware"
	service "github.com/tkeech1/gowebsvc/svc"
)

func getGreetingHandler(svc service.Greeter) *httptransport.Server {
	return httptransport.NewServer(
		makeGreetingEndpoint(svc),
		decodeGreetRequest,
		encodeResponse,
	)
}

func getExpensiveHandler(svc service.Greeter) *httptransport.Server {
	return httptransport.NewServer(
		makeExpensiveEndpoint(svc),
		decodeExpensiveRequest,
		encodeResponse,
	)
}

// main
func main() {
	//logger := kitlog.NewLogfmtLogger(os.Stdout)
	logger := log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "greeting_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "greeting_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc service.Greeter
	svc = service.GreetingService{}
	svc = middleware.LoggingMiddleware{logger, svc}
	svc = middleware.InstrumentingMiddleware{requestCount, requestLatency, svc}

	greetingHandler := getGreetingHandler(svc)
	expensiveHandler := getExpensiveHandler(svc)

	http.Handle("/greeting", greetingHandler)
	http.Handle("/expensive", expensiveHandler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
