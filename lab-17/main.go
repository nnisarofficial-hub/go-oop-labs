package main

import (
	"fmt"
	"time"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

type Response struct {
	StatusCode int
	Body       string
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}
type ConsoleLogger struct {
}

func (c ConsoleLogger) Info(msg string) {
	fmt.Println("[INFO]", msg)
}

func (c ConsoleLogger) Warn(msg string) {
	fmt.Println("[WARN]", msg)
}

func (c ConsoleLogger) Error(msg string) {
	fmt.Println("[ERROR]", msg)
}

type Handler func(req Request) Response

func WithLogging(logger Logger, next Handler) Handler {
	return func(req Request) Response {
		logger.Info(fmt.Sprintf("→ %s %s", req.Method, req.Path))
		resp := next(req)
		logger.Info(fmt.Sprintf("← %d %s", resp.StatusCode, req.Path))
		return resp
	}
}

func WithAuth(validToken string, next Handler) Handler {
	return func(req Request) Response {
		token := req.Headers["Authorization"]
		if token != "Bearer "+validToken {
			return Response{StatusCode: 401, Body: "| Body: unauthorized"}
		}
		return next(req)
	}
}

func WithTiming(next Handler) Handler {
	return func(req Request) Response {
		start := time.Now()
		resp := next(req)
		fmt.Printf("Request took: %v\n", time.Since(start))
		return resp
	}
}

func ProductHandler(req Request) Response {
	if req.Method == "GET" && req.Path == "/products" {
		return Response{StatusCode: 200, Body: `| Body: [{"id":"1","name":"Barbari Buck"}]`}
	}
	return Response{StatusCode: 404, Body: "not found"}
}

func main() {
	logger := ConsoleLogger{}
	handler := WithTiming(
		WithLogging(logger,
			WithAuth("secret-token",
				ProductHandler,
			),
		),
	)
	resp := handler(Request{
		Method:  "GET",
		Path:    "/products",
		Headers: map[string]string{"Authorization": "Bearer secret-token"},
	})
	fmt.Println("Response 1:", resp.StatusCode, resp.Body)
	fmt.Println()
	resp2 := handler(Request{
		Method: "GET",
		Path:   "/products",
	})
	fmt.Println("Response 2:", resp2.StatusCode, resp2.Body)
}
