package main

import (
	"log"
	"net/http"
)

type RequestLoggingDecorator struct {
	logger  *log.Logger
	handler http.Handler
}

func NewRequestLoggingDecorator(f http.Handler, l *log.Logger) http.Handler {
	return RequestLoggingDecorator{
		logger:  l,
		handler: f,
	}
}

// ServeHTTP implements http.Handler
func (d RequestLoggingDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.logger.Printf("Processing request to: %v - %v", r.Method, r.URL)
	d.handler.ServeHTTP(w, r)
	d.logger.Printf("Processed request to: %v - %v", r.Method, r.URL)
}
