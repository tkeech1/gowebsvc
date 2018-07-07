package main

import (
	"log"
	"net/http"
	"sync"

	service "github.com/tkeech1/gowebsvc/svc"
	"github.com/tkeech1/gowebsvc/transport"
)

// TODO - Router
//https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831

type server struct {
	//Shared dependencies are fields of the structure
	//db     *sql.DB
	//router *someRouter
	//email  EmailSender
	svc       service.Greeter
	transport transport.HttpJsonCoderDecoder
}

func (s *server) handleGreeting() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Printf("handler started")
		defer log.Printf("handler ended")

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
		log.Printf("handler started")
		defer log.Printf("handler ended")

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
			// initialize database, parse templates, etc.
			log.Print("middleware - just do this once")
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

// a middleware
func (s *server) MiddleWare(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("middleware - logic executed before the handler")
		h(w, r)
		log.Print("middleware - logic executed after the handler ")
	}
}

/*func (s *server) handleExpensive(text string) http.HandlerFunc {
	var (
		init sync.Once
		err  error
	)
	// custom response
	type response struct {
		Greeting string `json:"response"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// The prepareThing is called only once, so you can use it to do one-time per-handler initialisation, and then use the thing in the handler.
		// Be sure to only read the shared data, if handlers are modifying anything, remember you’ll need a mutex or something to protect it.
		//thing := prepareThing()
		init.Do(func() {
			// do an expensive operation here - it will only occur on the first invocation of the handler
			// initialize database, parse templates, etc.
			log.Print("middleware - just do this once")
		})
		// check the operation for errors
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpResponse := response{
			Greeting: text,
		}
		responseBody, err := json.Marshal(httpResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBody)
		//json.NewEncoder(w).Encode(httpResponse)
		//fmt.Fprintf(w, responseBody)
		// use the thing you just initialized
	}
}
*/

func main() {
	s := server{
		svc:       service.GreetingService{},
		transport: transport.HttpJson{},
	}
	http.HandleFunc("/greeting", s.MiddleWare(s.handleGreeting()))
	http.HandleFunc("/expensive", s.handleExpensive())
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
