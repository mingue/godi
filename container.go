package godi

import (
	"fmt"
	"reflect"
	"sync"
)

const (
	ErrDecoratedMustBeInterface = "item to decorate must be an interface"
	ErrFactoryAlreadyRegistered = "factory already registered"
	ErrFactoryNotRegistered     = "factory not registered"
	ErrDecoratorBeforeFactory   = "a factory needs to be registered before a decorator"
)

type Container struct {
	definitions    map[reflect.Type]*definition
	singletonCache *lifetimeCache
	scopedCache    *lifetimeCache
}

type definition struct {
	lifetime lifetime
	/*
		This is a slice of factories to contain the factory for the instance and the decorators
		The decorator have the following definition: func(decorated T, c *Container) T
		The item to create has the definition: func(c *Container) T
	*/
	f []any
}

func New() *Container {
	var singletonCache = lifetimeCache{
		entries: make(map[reflect.Type]*cacheEntry),
		mx:      sync.Mutex{},
	}

	var scopedCache = lifetimeCache{
		entries: make(map[reflect.Type]*cacheEntry),
		mx:      sync.Mutex{},
	}

	return &Container{
		definitions:    make(map[reflect.Type]*definition),
		singletonCache: &singletonCache,
		scopedCache:    &scopedCache,
	}
}

func Singleton[T any](c *Container, f func(c *Container) T) error {
	return add(c, LifetimeSingleton, f)
}

func Scoped[T any](c *Container, f func(c *Container) T) error {
	return add(c, LifetimeScoped, f)
}

func Transient[T any](c *Container, f func(c *Container) T) error {
	return add(c, LifetimeTransient, f)
}

func add[T any](c *Container, lifetime lifetime, f func(c *Container) T) error {
	factoryName := getKeyFromT[T]()

	_, exists := c.definitions[factoryName]

	if exists {
		return fmt.Errorf(ErrFactoryAlreadyRegistered)
	}

	c.definitions[factoryName] = &definition{
		lifetime: lifetime,
		f:        []any{f},
	}

	return nil
}

func Decorate[T any](c *Container, f func(decorated T, c *Container) T) error {
	target := getKeyFromT[T]()

	_, found := c.definitions[target]
	if !found {
		return fmt.Errorf(ErrDecoratorBeforeFactory)
	}

	if target.Elem().Kind() != reflect.Interface {
		return fmt.Errorf(ErrDecoratedMustBeInterface)
	}

	for key, def := range c.definitions {
		if target == key {
			def.f = append(def.f, f)
		}
	}

	return nil
}

func Get[T any](c *Container) (T, error) {
	key := getKeyFromT[T]()

	definition, exists := c.definitions[key]

	if !exists {
		var result T
		return result, fmt.Errorf(ErrFactoryNotRegistered)
	}

	if definition.lifetime == LifetimeSingleton {
		value, found := getFromCacheOrBuild[T](c, c.singletonCache, key, definition)
		return value, found
	}

	if definition.lifetime == LifetimeScoped {
		value, found := getFromCacheOrBuild[T](c, c.scopedCache, key, definition)
		return value, found
	}

	return buildItem[T](c, definition), nil
}

func (c *Container) NewScope() *Container {
	var newScopedCache = lifetimeCache{
		entries: make(map[reflect.Type]*cacheEntry),
		mx:      sync.Mutex{},
	}

	return &Container{
		definitions:    c.definitions,
		singletonCache: c.singletonCache,
		scopedCache:    &newScopedCache,
	}
}

func getFromCacheOrBuild[T any](c *Container, cache *lifetimeCache, key reflect.Type, d *definition) (T, error) {
	value, found := cache.entries[key]

	if !found {
		cache.mx.Lock()

		value, found = cache.entries[key]

		if !found {
			entry := cacheEntry{}
			cache.entries[key] = &entry
			value = &entry
		}

		cache.mx.Unlock()

		if !value.initialized {
			value.mx.Lock()

			if !value.initialized {
				value.initialized = true
				value.instance = buildItem[T](c, d)
			}

			value.mx.Unlock()
		}
	}

	return value.instance.(T), nil
}

func buildItem[T any](c *Container, d *definition) T {
	var factory = d.f[0].(func(c *Container) T)
	value := factory(c)

	for i := 1; i <= len(d.f)-1; i++ {
		decoratorF := d.f[i].(func(decorated T, c *Container) T)
		value = decoratorF(value, c)
	}

	return value
}

func getKeyFromT[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil))
}
