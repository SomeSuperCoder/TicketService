package handlers

import (
	"context"
	"fmt"
)

type HelloHandler struct{}

type HelloRequest struct {
	Name string `query:"name"`
}
type HelloResponse struct {
	Body struct {
		Greeting string `json:"name"`
	}
}

func (h *HelloHandler) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	resp := new(HelloResponse)
	resp.Body.Greeting = fmt.Sprintf("Hello, %s!", req.Name)

	return resp, nil
}
