package invoice

import (
	"context"
	"fmt"
	"net/http"
)

type ReadyHandler struct {
	Ctx context.Context
}

func (h *ReadyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	val := h.Ctx.Value("Some")
	fmt.Printf("Printing value from repo: %v\n", val)

	w.Write([]byte("ok"))
}
