package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	middleware "github.com/tkeech1/gowebsvc/middleware"
	service "github.com/tkeech1/gowebsvc/svc"

	"github.com/stretchr/testify/assert"
)

func Test_Greet(t *testing.T) {
	tests := map[string]struct {
		ctx              context.Context
		svc              service.GreetingService
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

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "Test_GreetingService",
		Subsystem: "greeting_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "Test_GreetingService",
		Subsystem: "greeting_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	tests := map[string]struct {
		svc                service.Greeter
		logger             *log.Logger
		greeting           []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			greeting:           []byte(`{"s":"hello"}`),
			expectedResponse:   `{"greeting":"hello"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_nogreeting": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			greeting:           []byte(`{"s":""}`),
			expectedResponse:   `{"greeting":"","err":"empty greeting"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptyjson": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			greeting:           []byte(`{}`),
			expectedResponse:   `{"greeting":"","err":"empty greeting"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptymessage": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			greeting:           []byte(``),
			expectedResponse:   `{"greeting":"","err":"EOF"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		req, err := http.NewRequest("POST", "/greeting", bytes.NewBuffer(test.greeting))
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()

		instrumentingMiddleware := middleware.InstrumentingMiddleware{
			RequestCount:   requestCount,
			RequestLatency: requestLatency,
			Next:           test.svc,
		}
		logMiddleware := middleware.LoggingMiddleware{
			Logger: test.logger,
			Next:   instrumentingMiddleware,
		}
		s := server{transport: HttpJson{}, svc: logMiddleware}

		handler := s.handleGreeting()
		handler.ServeHTTP(w, req)
		assert.Equal(t, test.expectedResponse, w.Body.String())
		assert.Equal(t, test.httpStatusResponse, w.Code)
	}

	// check prometheus stats
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	w := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(w, req)

	parser := expfmt.TextParser{}
	parsedData, err := parser.TextToMetricFamilies(w.Body)
	if err != nil {
		t.Fatal(" unable to get prometheus metrics ")
	}

	var errorCount, successCount float64
	for _, metric := range parsedData["Test_GreetingService_greeting_service_request_count"].GetMetric() {
		for _, label := range metric.GetLabel() {
			if label.GetName() == "error" && label.GetValue() == "true" {
				errorCount = metric.GetCounter().GetValue()
			}
			if label.GetName() == "error" && label.GetValue() == "false" {
				successCount = metric.GetCounter().GetValue()
			}
		}
	}
	assert.Equal(t, 2.0, errorCount)
	assert.Equal(t, 1.0, successCount)
}

func Test_GreetingServiceCancelContext(t *testing.T) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 0*time.Millisecond)
	defer cancel()

	tests := map[string]struct {
		svc                service.GreetingService
		logger             *log.Logger
		greeting           []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			greeting:           []byte(`{"s":"hello"}`),
			expectedResponse:   `{"greeting":"","err":"request cancelled"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		req, err := http.NewRequest("POST", "/greeting", bytes.NewBuffer(test.greeting))
		req = req.WithContext(ctx)
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()

		logMiddleware := middleware.LoggingMiddleware{Logger: test.logger, Next: test.svc}
		s := server{transport: HttpJson{}, svc: logMiddleware}

		handler := s.handleGreeting()
		handler.ServeHTTP(w, req)
		assert.Equal(t, test.expectedResponse, w.Body.String())
		assert.Equal(t, test.httpStatusResponse, w.Code)
	}
}

func Test_ExpensiveService(t *testing.T) {

	tests := map[string]struct {
		svc                service.GreetingService
		logger             *log.Logger
		expensive          []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			expensive:          []byte(`{"connection_string":"c1","username":"u1","password":"p1"}`),
			expectedResponse:   `{"status":"c1u1p1"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_noconnection": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			expensive:          []byte(`{"connection_string":"","username":"u1","password":"p1"}`),
			expectedResponse:   `{"status":"","err":"missing connectionString"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_nousername": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			expensive:          []byte(`{"connection_string":"c1","username":"","password":"p1"}`),
			expectedResponse:   `{"status":"","err":"missing username"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_nopassword": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			expensive:          []byte(`{"connection_string":"c1","username":"u1","password":""}`),
			expectedResponse:   `{"status":"","err":"missing password"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptyjson": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			expensive:          []byte(`{}`),
			expectedResponse:   `{"status":"","err":"missing connectionString"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
		"error_emptymessage": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
			expensive:          []byte(``),
			expectedResponse:   `{"status":"","err":"EOF"}` + "\n",
			httpStatusResponse: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		req, err := http.NewRequest("POST", "/expensive", bytes.NewBuffer(test.expensive))
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()

		logMiddleware := middleware.LoggingMiddleware{Logger: test.logger, Next: test.svc}
		s := server{transport: HttpJson{}, svc: logMiddleware}

		handler := s.handleExpensive()
		handler.ServeHTTP(w, req)
		assert.Equal(t, test.expectedResponse, w.Body.String())
		assert.Equal(t, test.httpStatusResponse, w.Code)
	}
}

func Test_ExpensiveServiceMultipleTries(t *testing.T) {

	tests := map[string]struct {
		svc                service.GreetingService
		logger             *log.Logger
		expensive          []byte
		expectedResponse   string
		httpStatusResponse int
	}{
		"success": {
			svc:                service.GreetingService{},
			logger:             log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
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

	logMiddleware := middleware.LoggingMiddleware{Logger: tests["success"].logger, Next: tests["success"].svc}
	s := server{transport: HttpJson{}, svc: logMiddleware}
	handler := s.handleExpensive()

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
