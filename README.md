# GODI

Simple & Performant Dependency Injection Container for Go  

## Objectives  

- Remove the need to pass context.Context to intermediate dependencies for cancellations and avoid breaking changes to existing APIs
- Allow injection of scoped requests without the need to use context.WithValue which is not strongly typed
- Be able to easily decorate an interface to extend it's functionality
- Be able to create custom scopes to match different lifetimes: request, session...
- Be compatible with std APIs so that it can be used in all the solutions
- Strongly typed API, making use of generics
- Thread safe object creation
- Easy to test, no globals or init functions

## Usage  

### Create a container

```go
// Create an empty container
cont := godi.New()
```

### Register factories for dependencies

```go

// Transient instances return new instances every time the are being requested
godi.Transient(cont, func(c *godi.Container) SomeInterface {
    return &SomeStruct{}
})

// Singleton instances return always the same instance for the duration of the process
godi.Singleton(cont, func(c *godi.Container) SomeInterface {
    return &SomeStruct{}
})

// Scoped instances return the same instance per container scope
godi.Scoped(cont, func(c *godi.Container) SomeInterface {
    return &SomeStruct{}
})

```

### Create new container scopes per request, session... as required

```go

newScopedContainer := cont.NewScope()

// Registered scoped definitions are only available for the scoped container
// and could be injected on any object of the stack

// Like passing the context only where needed for cancellations 
// without having to pass it across the whole stack, 
// or inject distributed tracing information where needed

godi.Scoped(newScopedContainer, func(c *godi.Container) context.Context {
    return r.Context()
})

```

### Decorate definitions to easily wrap and extend their functionality

```go

// Decorate each previous definition for the interface extending its functionality
godi.Decorate(cont, func(d http.Handler, c *godi.Container) http.Handler {
    logger, _ := godi.Get[*log.Logger](c)
    return NewRequestLoggingDecorator(d, logger)
})

```

### Get instances from the container

```go

h, _ := godi.Get[http.Handler](requestCont)
h.ServeHTTP(w, r)

```

### Register several implementations of the same interface by using the Named options

```go

// Register several factories for the same interface http.Handler
godi.ScopedNamed(cont, "invoiceHandler", func(c *godi.Container) http.Handler {
    svc, _ := godi.Get[invoice.InvoiceService](c)
    return invoice.NewInvoiceHandler(svc)
})

godi.ScopedNamed(cont, "readyHandler", func(c *godi.Container) http.Handler {
    ctx, _ := godi.Get[context.Context](c)
    return &invoice.ReadyHandler{Ctx: ctx}
})

// Obtain the named registrations from the container
invoiceHandler, _ := godi.GetNamed[http.Handler](cont, "invoiceHandler")
readyHandler, _ := godi.GetNamed[http.Handler](cont, "readyHandler")

// Decorators will apply to all named and not named registrations for the interface
```

## Example Application  

See <https://github.com/mingue/godi/blob/main/example/cmd/server/main.go>

## Things pending to investigate or implement  

- [x] Allow to register several items for the same interface, like http.Handlers
- [] Add syntactic sugar for http handler registration to reduce boilerplate
- [] Ensure that instances with limited lifetimes: scoped or transient, are not injected into Singletons
- [] Investigate usage of interface to enable function overload on existing APIs, factory func, func or T
- [] Allow to register with constructors as per dig Invoke, requires benchmarking
- [] Container interceptors or hooks for debugging or visibility
