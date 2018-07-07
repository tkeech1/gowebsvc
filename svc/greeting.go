package svc

import (
	"context"
	"errors"
	"time"
)

type Greeter interface {
	Greet(context.Context, string) (string, error)
	Expensive(context.Context, string, string, string) (string, error)
}

type GreetingService struct{}

func (g GreetingService) Greet(ctx context.Context, greeting string) (string, error) {
	ch := make(chan string)

	// in case this is a long-runnning operation, use a go function
	go func(s string) {
		ch <- s
	}(greeting)

	select {
	case response := <-ch:
		if response == "" {
			return "", errors.New("empty greeting")
		}
		return response, nil
	case <-ctx.Done():
		return "", errors.New("request cancelled")
	}
}

func (g GreetingService) Expensive(ctx context.Context, connectionString, username, password string) (string, error) {
	if connectionString == "" {
		return "", errors.New("missing connectionString")
	}
	if username == "" {
		return "", errors.New("missing username")
	}
	if password == "" {
		return "", errors.New("missing password")
	}

	ch := make(chan string)
	go func(c, u, p string) {
		ch <- connectionString + username + password
	}(connectionString, username, password)

	select {
	case response := <-ch:
		if response == "" {
			return "", errors.New("empty greeting")
		}
		return response, nil
	case <-ctx.Done():
		return "", errors.New("request cancelled")
	case <-time.After(1 * time.Second):
		return "", errors.New("request timed out")
	}
}

/*
func (s *server) handleExpensive(text string) http.HandlerFunc {
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
