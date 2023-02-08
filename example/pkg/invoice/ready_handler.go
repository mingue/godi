package invoice

import (
	"log"
	"net/http"
)

type ReadyHandler struct {
	requestContext RequestContext
}

func NewReadyHandler(requestContext RequestContext) *ReadyHandler {
	return &ReadyHandler{requestContext: requestContext}
}

func (r *ReadyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("Printing value: %v, User Agent: %v, Request count: %v", r.requestContext.SomeValue, r.requestContext.UserAgent, r.requestContext.Counter)
	w.Write([]byte("ok"))
}
