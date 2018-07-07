package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	service "github.com/tkeech1/gowebsvc/svc"
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
	logMiddleware := loggingMiddleware{
		logger: log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
		next:   service.GreetingService{},
	}
	s := server{transport: HttpJson{}, svc: logMiddleware}

	http.HandleFunc("/greeting", s.handleGreeting())
	http.HandleFunc("/expensive", s.handleExpensive())
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
