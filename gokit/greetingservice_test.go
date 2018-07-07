package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	service "github.com/tkeech1/gowebsvc/svc"
)

func Test_Greet(t *testing.T) {
	tests := map[string]struct {
		ctx              context.Context
		svc              service.Greeter
		greeting         string
		expectedResponse string
		errorResponse    error
	}{
		"success": {
			ctx:              context.Background(),
			svc:              service.GreetingService{},
			greeting:         "helloaskjdfhsl",
			expectedResponse: "helloaskjdfhsl",
			errorResponse:    nil,
		},
		"error": {
			ctx:              context.Background(),
			svc:              service.GreetingService{},
			greeting:         "",
			expectedResponse: "",
			errorResponse:    errors.New("empty greeting"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		response, err := test.svc.Greet(test.ctx, test.greeting)
		assert.Equal(t, test.expectedResponse, response)
		assert.Equal(t, test.errorResponse, err)
	}

}

func Test_GreetingService(t *testing.T) {

	tests := map[string]struct {
		svc                service.Greeter
		logger             kitlog.Logger
		greeting           []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			greeting:           []byte(`{"s":"hello"}`),
			expectedResponse:   `{"greeting":"hello"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_nogreeting": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			greeting:           []byte(`{"s":""}`),
			expectedResponse:   `{"greeting":"","err":"empty greeting"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptyjson": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			greeting:           []byte(`{}`),
			expectedResponse:   `{"greeting":"","err":"empty greeting"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptymessage": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			greeting:           []byte(``),
			expectedResponse:   `EOF`,
			httpStatusResponse: http.StatusInternalServerError,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		test.svc = loggingMiddleware{test.logger, test.svc}

		req, err := http.NewRequest("POST", "/greeting", bytes.NewBuffer(test.greeting))
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()

		handler := getGreetingHandler(test.svc)
		handler.ServeHTTP(w, req)
		assert.Equal(t, test.expectedResponse, w.Body.String())
		assert.Equal(t, test.httpStatusResponse, w.Code)
	}
}

func Test_GreetingServiceCancelContext(t *testing.T) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 0*time.Millisecond)
	defer cancel()

	tests := map[string]struct {
		svc                service.Greeter
		logger             kitlog.Logger
		greeting           []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			greeting:           []byte(`{"s":"hello"}`),
			expectedResponse:   `{"greeting":"","err":"request cancelled"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		test.svc = loggingMiddleware{test.logger, test.svc}

		req, err := http.NewRequest("POST", "/greeting", bytes.NewBuffer(test.greeting))
		req = req.WithContext(ctx)
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()

		handler := getGreetingHandler(test.svc)
		handler.ServeHTTP(w, req)
		assert.Equal(t, test.expectedResponse, w.Body.String())
		assert.Equal(t, test.httpStatusResponse, w.Code)
	}
}

func Test_ExpensiveService(t *testing.T) {

	tests := map[string]struct {
		svc                service.Greeter
		logger             kitlog.Logger
		expensive          []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			expensive:          []byte(`{"connection_string":"c1","username":"u1","password":"p1"}`),
			expectedResponse:   `{"status":"c1u1p1"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_noconnection": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			expensive:          []byte(`{"connection_string":"","username":"u1","password":"p1"}`),
			expectedResponse:   `{"status":"","err":"missing connectionString"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_nousername": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			expensive:          []byte(`{"connection_string":"c1","username":"","password":"p1"}`),
			expectedResponse:   `{"status":"","err":"missing username"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_nopassword": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			expensive:          []byte(`{"connection_string":"c1","username":"u1","password":""}`),
			expectedResponse:   `{"status":"","err":"missing password"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptyjson": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			expensive:          []byte(`{}`),
			expectedResponse:   `{"status":"","err":"missing connectionString"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptymessage": {
			svc:                service.GreetingService{},
			logger:             kitlog.NewLogfmtLogger(os.Stdout),
			expensive:          []byte(``),
			expectedResponse:   `EOF`,
			httpStatusResponse: http.StatusInternalServerError,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		test.svc = loggingMiddleware{test.logger, test.svc}
		req, err := http.NewRequest("POST", "/expensive", bytes.NewBuffer(test.expensive))
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()

		handler := getExpensiveHandler(test.svc)
		handler.ServeHTTP(w, req)
		assert.Equal(t, test.expectedResponse, w.Body.String())
		assert.Equal(t, test.httpStatusResponse, w.Code)
	}
}

func Test_ExpensiveServiceMultipleTries(t *testing.T) {

	tests := map[string]struct {
		expensive          []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			expensive:          []byte(`{"connection_string":"c1","username":"u2","password":"p3"}`),
			expectedResponse:   `{"status":"c1u2p3"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"2nd_try": {
			expensive:          []byte(`{"connection_string":"","username":"hello","password":"hello"}`),
			expectedResponse:   `{"status":"already initialized"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
	}

	var svc service.Greeter
	svc = service.GreetingService{}
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	svc = loggingMiddleware{logger, svc}
	handler := getExpensiveHandler(svc)

	t.Logf("Running test case: %s", "success")
	req, err := http.NewRequest("POST", "/expensive", bytes.NewBuffer(tests["success"].expensive))
	if err != nil {
		t.Errorf(err.Error())
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, tests["success"].expectedResponse, w.Body.String())
	assert.Equal(t, tests["success"].httpStatusResponse, w.Code)

	t.Logf("Running test case: %s", "2nd_try")
	req, err = http.NewRequest("POST", "/expensive", bytes.NewBuffer(tests["2nd_try"].expensive))
	if err != nil {
		t.Errorf(err.Error())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, tests["2nd_try"].expectedResponse, w.Body.String())
	assert.Equal(t, tests["2nd_try"].httpStatusResponse, w.Code)

}

/*
func TestServer_HandleGreetingMiddlewareCancelContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	srv := server{}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.handleGreeting("test"))
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected response %v\n", w)
	}

	expected := "context deadline exceeded\n"
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
	}
}
*/
