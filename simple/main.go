package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tkeech1/gowebsvc/middleware"
	service "github.com/tkeech1/gowebsvc/svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO - Router
//https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831

type server struct {
	svc       service.Greeter
	transport HttpJsonCoderDecoder
}

func (s *server) handleGreeting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var response service.GreetResponse

		gr, err := s.transport.DecodeGreetingServiceRequest(r)
		if err != nil {
			response = service.GreetResponse{
				V:   "",
				Err: err.Error(),
			}
			s.transport.EncodeGreetingServiceRequest(&w, response)
			return
		}

		greeting, err := s.svc.Greet(ctx, gr.S)
		if err != nil {
			response = service.GreetResponse{
				V:   "",
				Err: err.Error(),
			}
			s.transport.EncodeGreetingServiceRequest(&w, response)
			return
		}

		response = service.GreetResponse{
			V:   greeting,
			Err: "",
		}
		s.transport.EncodeGreetingServiceRequest(&w, response)
	}
}

func (s *server) handleExpensive() http.HandlerFunc {
	var (
		init sync.Once
	)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var response service.ExpensiveResponse

		gr, err := s.transport.DecodeExpensiveServiceRequest(r)
		if err != nil {
			response = service.ExpensiveResponse{
				V:   "",
				Err: err.Error(),
			}
			s.transport.EncodeExpensiveServiceRequest(&w, response)
			return
		}

		expensive := "already initialized"
		init.Do(func() {
			// do an expensive operation here - it will only occur on the first invocation of the handler
			expensive, err = s.svc.Expensive(ctx, gr.C, gr.U, gr.P)
		})

		if err != nil {
			response = service.ExpensiveResponse{
				V:   "",
				Err: err.Error(),
			}
			s.transport.EncodeExpensiveServiceRequest(&w, response)
			return
		}

		response = service.ExpensiveResponse{
			V:   expensive,
			Err: "",
		}
		s.transport.EncodeExpensiveServiceRequest(&w, response)
	}
}

func main() {

	//GRPC 
	go func() {

		logMiddleware := middleware.LoggingMiddlewareGRPC{
			Logger: log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			Next:   &service.GreetingServiceGRPC{},
		}

		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		service.RegisterGreetingServiceServer(s, &logMiddleware)

		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// end GRPC

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "Test_GreetingServiceCancelContext",
		Subsystem: "greeting_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "Test_GreetingServiceCancelContext",
		Subsystem: "greeting_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	instrumentingMiddleware := middleware.InstrumentingMiddleware{
		RequestCount:   requestCount,
		RequestLatency: requestLatency,
		Next:           service.GreetingService{},
	}
	logMiddleware := middleware.LoggingMiddleware{
		Logger: log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
		Next:   instrumentingMiddleware,
	}

	s := server{transport: HttpJson{}, svc: logMiddleware}
	http.HandleFunc("/greeting", s.handleGreeting())
	http.HandleFunc("/expensive", s.handleExpensive())
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))

}
