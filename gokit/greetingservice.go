package main

import (
	"log"
	"os"

	"net/http"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
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
	logger := kitlog.NewLogfmtLogger(os.Stdout)

	var svc service.Greeter
	svc = service.GreetingService{}
	svc = loggingMiddleware{logger, svc}

	greetingHandler := getGreetingHandler(svc)
	expensiveHandler := getExpensiveHandler(svc)

	http.Handle("/greeting", greetingHandler)
	http.Handle("/expensive", expensiveHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
