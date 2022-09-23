package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mingue/godi"
	"github.com/mingue/godi/example/pkg/invoice"
)

func main() {
	fmt.Println("Starting execution")

	// Create Container
	cont := godi.New()

	// Register dependencies
	godi.Scoped(cont, func(c *godi.Container) http.Handler {
		svc, _ := godi.Get[invoice.InvoiceService](c)
		return invoice.NewInvoiceHandler(svc)
	})

	godi.Scoped(cont, func(c *godi.Container) invoice.InvoiceRepository {
		ctx, _ := godi.Get[context.Context](c)
		return invoice.NewInvoiceRepositoryImpl(ctx)
	})

	godi.Scoped(cont, func(c *godi.Container) invoice.InvoiceService {
		repo, _ := godi.Get[invoice.InvoiceRepository](c)
		return invoice.NewInvoiceServiceImpl(repo)
	})

	godi.Scoped(cont, func(c *godi.Container) *log.Logger {
		return log.New(os.Stdout, "App: ", log.Default().Flags())
	})

	godi.Decorate(cont, func(d http.Handler, c *godi.Container) http.Handler {
		logger, _ := godi.Get[*log.Logger](c)
		return NewRequestLoggingDecorator(d, logger)
	})

	// Start the http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestCont := cont.NewScope()
		ctx := context.WithValue(r.Context(), "Some", "Value")

		godi.Scoped(requestCont, func(c *godi.Container) context.Context {
			return ctx
		})
		h, _ := godi.Get[http.Handler](requestCont)
		h.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Couldn't start the server")
	}
}
