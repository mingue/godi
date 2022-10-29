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
	rootPath := "/"
	readyPath := "/ready"

	godi.ScopedNamed(cont, rootPath, func(c *godi.Container) http.Handler {
		svc, _ := godi.Get[invoice.InvoiceService](c)
		return invoice.NewInvoiceHandler(svc)
	})

	godi.ScopedNamed(cont, readyPath, func(c *godi.Container) http.Handler {
		ctx, _ := godi.Get[context.Context](c)
		return &invoice.ReadyHandler{Ctx: ctx}
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
	http.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
		requestCont := cont.NewScope()

		// Only as example, you can now register a struct to be injected anywhere on the stack
		ctx := context.WithValue(r.Context(), "Some", "OnRootPath")
		godi.Scoped(requestCont, func(c *godi.Container) context.Context {
			return ctx
		})
		// if err != nil {
		// 	fmt.Printf("Error registering the context: %v \n", err)
		// 	r.Response.StatusCode = 500
		// 	fmt.Fprint(w, fmt.Sprintf("Error registering the context: %v \n", err))
		// 	return
		// }

		h, _ := godi.GetNamed[http.Handler](requestCont, rootPath)
		h.ServeHTTP(w, r)
	})

	http.HandleFunc(readyPath, func(w http.ResponseWriter, r *http.Request) {
		requestCont := cont.NewScope()

		// Only as example, you can now register a struct to be injected anywhere on the stack
		ctx := context.WithValue(r.Context(), "Some", "OnReadyPath")
		err := godi.Scoped(requestCont, func(c *godi.Container) context.Context {
			return ctx
		})
		if err != nil {
			r.Response.StatusCode = 500
			fmt.Fprint(w, fmt.Sprintf("Error registering the context: %v \n", err))
			return
		}

		h, _ := godi.GetNamed[http.Handler](requestCont, readyPath)
		h.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Couldn't start the server")
	}
}
