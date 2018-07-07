package svc

type GreetRequest struct {
	S string `json:"s"`
}

type GreetResponse struct {
	V   string `json:"greeting"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type ExpensiveRequest struct {
	C string `json:"connection_string"`
	U string `json:"username"`
	P string `json:"password"`
}

type ExpensiveResponse struct {
	V   string `json:"status"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}
