package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mingue/godi"
	"github.com/mingue/godi/example/pkg/invoice"
)

func main() {
	initServer()
}

func initServer() {
	log.Printf("Starting execution...")

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
		requestContext, _ := godi.Get[invoice.RequestContext](c)
		return invoice.NewReadyHandler(requestContext)
	})

	godi.Scoped(cont, func(c *godi.Container) invoice.InvoiceRepository {
		requestContext, _ := godi.Get[invoice.RequestContext](c)
		return invoice.NewInvoiceRepositoryImpl(requestContext)
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

	// Register http handlers
	rootCounter := 0
	http.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
		// Create a new container with scope for the http request
		requestCont := cont.NewScope()

		rootCounter++

		// Register any http request scoped object, which can be enriched with anything from the http.Request
		godi.Scoped(requestCont, func(c *godi.Container) invoice.RequestContext {
			userAgent := r.Header["User-Agent"][0]
			return invoice.RequestContext{SomeValue: "On root path", UserAgent: userAgent, Counter: rootCounter}
		})

		// When you get an object from the container, it can have dependencies on the http request scoped dependencies registered above
		// There is no need to pass the objects across the stack, can be injected to any object
		h, _ := godi.GetNamed[http.Handler](requestCont, rootPath)
		h.ServeHTTP(w, r)
	})

	readyCounter := 0

	http.HandleFunc(readyPath, func(w http.ResponseWriter, r *http.Request) {
		requestCont := cont.NewScope()

		readyCounter++

		godi.Scoped(requestCont, func(c *godi.Container) invoice.RequestContext {
			userAgent := r.Header["User-Agent"][0]
			return invoice.RequestContext{SomeValue: "On ready path", UserAgent: userAgent, Counter: readyCounter}
		})

		h, _ := godi.GetNamed[http.Handler](requestCont, readyPath)
		h.ServeHTTP(w, r)
	})

	// Start the http server
	log.Printf("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
