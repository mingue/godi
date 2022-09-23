# GODI

Simple & Performant Dependency Injection Container for Go  

## Objectives  

- Remove the need to pass context.Context everywhere and avoid breaking changes to existing APIs
- Be able to easily decorate an interface to extend it's functionality
- Be able to create custom scopes to match different lifetimes: request, session...
- Be compatible with std APIs for easy usage
- Strongly typed API, making use of generics
- Thread safe
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

// Scoped instances return the same instance per container
godi.Scoped(cont, func(c *godi.Container) SomeInterface {
    return &SomeStruct{}
})

```

### Create new container scopes per request, session... as required

```go

newScopedContainer := cont.NewScope()

```

### Decorate definitions to easily wrap and extend their functionality

```go

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

## Example Application  

See <https://github.com/mingue/godi/blob/main/example/cmd/server/main.go>

## Things pending to investigate or implement  

- [] Allow to register several items for the same interface, like http.Handlers
- [] Add syntactic sugar for http handler registration to reduce boilerplate
- [] Ensure that instances with limited lifetimes: scoped or transient, are not injected into Singletons
- [] Investigate usage of interface to enable function overload on existing APIs, factory func, func or T
- [] Allow to register with constructors as per dig Invoke, requires benchmarking
- [] Container interceptors or hooks for debugging or visibility
