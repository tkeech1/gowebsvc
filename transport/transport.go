package transport

import (
	"encoding/json"
	"net/http"

	service "github.com/tkeech1/gowebsvc/svc"
)

type HttpJsonCoderDecoder interface {
	DecodeGreetingServiceRequest(*http.Request) (service.GreetRequest, error)
	EncodeGreetingServiceRequest(*http.ResponseWriter, service.GreetResponse) error
	DecodeExpensiveServiceRequest(*http.Request) (service.ExpensiveRequest, error)
	EncodeExpensiveServiceRequest(*http.ResponseWriter, service.ExpensiveResponse) error
}

type HttpJson struct{}

func (s HttpJson) DecodeGreetingServiceRequest(r *http.Request) (service.GreetRequest, error) {
	var request service.GreetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return service.GreetRequest{}, err
	}
	return request, nil
}

func (s HttpJson) EncodeGreetingServiceRequest(w *http.ResponseWriter, response service.GreetResponse) error {
	return json.NewEncoder(*w).Encode(response)
}

func (s HttpJson) DecodeExpensiveServiceRequest(r *http.Request) (service.ExpensiveRequest, error) {
	var request service.ExpensiveRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return service.ExpensiveRequest{}, err
	}
	return request, nil
}

func (s HttpJson) EncodeExpensiveServiceRequest(w *http.ResponseWriter, response service.ExpensiveResponse) error {
	return json.NewEncoder(*w).Encode(response)
}
