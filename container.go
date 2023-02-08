// GODI a idiomatic DI container for go
package godi

import (
	"errors"	
	"reflect"
	"sync"
)

var (
	ErrDecoratedMustBeInterface = errors.New("item to decorate must be an interface")
	ErrFactoryAlreadyRegistered = errors.New("factory already registered")
	ErrFactoryNotRegistered     = errors.New("factory not registered")
	ErrDecoratorBeforeFactory   = errors.New("a factory needs to be registered before a decorator")
)

type Container struct {
	globalDef      map[reflect.Type]map[string]*definition
	scopedDef      map[reflect.Type]map[string]*definition
	singletonCache *lifetimeCache
	scopedCache    *lifetimeCache
}

type definition struct {
	lifetime lifetime
	/*
		This is a slice of factories to contain the factory for the instance and the decorators
		The decorators have the following definition: func(decorated T, c *Container) T
		The item to create has the definition: func(c *Container) T
	*/
	f []any
}

func New() *Container {
	var singletonCache = lifetimeCache{
		entries: make(map[reflect.Type]map[string]*cacheEntry),
		mx:      sync.Mutex{},
	}

	var scopedCache = lifetimeCache{
		entries: make(map[reflect.Type]map[string]*cacheEntry),
		mx:      sync.Mutex{},
	}

	return &Container{
		globalDef:      make(map[reflect.Type]map[string]*definition),
		singletonCache: &singletonCache,
		scopedCache:    &scopedCache,
	}
}

func Singleton[T any](c *Container, f func(c *Container) T) error {
	return add(c, "", LifetimeSingleton, f)
}

func Scoped[T any](c *Container, f func(c *Container) T) error {
	return add(c, "", LifetimeScoped, f)
}

func Transient[T any](c *Container, f func(c *Container) T) error {
	return add(c, "", LifetimeTransient, f)
}

func SingletonNamed[T any](c *Container, name string, f func(c *Container) T) error {
	return add(c, name, LifetimeSingleton, f)
}

func ScopedNamed[T any](c *Container, name string, f func(c *Container) T) error {
	return add(c, name, LifetimeScoped, f)
}

func TransientNamed[T any](c *Container, name string, f func(c *Container) T) error {
	return add(c, name, LifetimeTransient, f)
}

func add[T any](c *Container, name string, lifetime lifetime, f func(c *Container) T) error {
	factoryName := getKeyFromT[T]()

	var typeDef map[string]*definition
	foundTypeDef := false

	// Search for existing type def in scoped def if exists or on the global map
	if c.scopedDef != nil {
		typeDef, foundTypeDef = c.scopedDef[factoryName]
	}

	if !foundTypeDef {
		typeDef, foundTypeDef = c.globalDef[factoryName]
	}

	// If a definition exist for the same type and name throw err
	if foundTypeDef {
		_, foundNamedDef := typeDef[name]

		if foundNamedDef {
			return ErrFactoryAlreadyRegistered
		}
	}

	// If typeDef not found create a new one in either scopedDef or globalDef
	if !foundTypeDef {
		typeDef = make(map[string]*definition)

		if lifetime == LifetimeScoped && c.scopedDef != nil {
			c.scopedDef[factoryName] = typeDef
		} else {
			c.globalDef[factoryName] = typeDef
		}
	}

	typeDef[name] = &definition{
		lifetime: lifetime,
		f:        []any{f},
	}

	return nil
}

func Decorate[T any](c *Container, f func(decorated T, c *Container) T) error {
	target := getKeyFromT[T]()

	var typeDef map[string]*definition
	foundTypeDef := false

	if c.scopedDef != nil {
		typeDef, foundTypeDef = c.scopedDef[target]
	}

	if !foundTypeDef {
		typeDef, foundTypeDef = c.globalDef[target]
	}

	if !foundTypeDef {
		return ErrDecoratorBeforeFactory
	}

	if target.Elem().Kind() != reflect.Interface {
		return ErrDecoratedMustBeInterface
	}

	// For each registration for the type we add the decorated as we can have named definitions
	for _, namedDef := range typeDef {
		namedDef.f = append(namedDef.f, f)
	}

	return nil
}

func Get[T any](c *Container) (T, error) {
	return GetNamed[T](c, "")
}

func GetNamed[T any](c *Container, name string) (T, error) {
	key := getKeyFromT[T]()

	var typeDef map[string]*definition
	foundTypeDef := false

	if c.scopedDef != nil {
		typeDef, foundTypeDef = c.scopedDef[key]
	}

	if !foundTypeDef {
		typeDef, foundTypeDef = c.globalDef[key]
	}

	if !foundTypeDef {
		var result T

		return result, ErrFactoryNotRegistered
	}

	if len(typeDef) == 0 {
		var result T

		return result, ErrFactoryNotRegistered
	}

	namedDef, foundNamedDef := typeDef[name]

	if !foundNamedDef {
		var result T

		return result, ErrFactoryNotRegistered
	}

	if namedDef.lifetime == LifetimeSingleton {
		value, found := getFromCacheOrBuild[T](c, c.singletonCache, key, name, namedDef)

		return value, found
	}

	if namedDef.lifetime == LifetimeScoped {
		value, found := getFromCacheOrBuild[T](c, c.scopedCache, key, name, namedDef)

		return value, found
	}

	return buildItem[T](c, namedDef), nil
}

func GetNoAlloc[T any](c *Container, x *T) error {
	return GetNamedNoAlloc(c, x, "")
}

func GetNamedNoAlloc[T any](c *Container, x *T, name string) error {
	key := getKeyFromT[T]()

	var typeDef map[string]*definition
	foundTypeDef := false

	if c.scopedDef != nil {
		typeDef, foundTypeDef = c.scopedDef[key]
	}

	if !foundTypeDef {
		typeDef, foundTypeDef = c.globalDef[key]
	}

	if !foundTypeDef {
		return ErrFactoryNotRegistered
	}

	if len(typeDef) == 0 {
		return ErrFactoryNotRegistered
	}

	namedDef, foundNamedDef := typeDef[name]

	if !foundNamedDef {
		return ErrFactoryNotRegistered
	}

	if namedDef.lifetime == LifetimeSingleton {
		found := getFromCacheOrBuildNoAlloc(c, c.singletonCache, key, name, namedDef, x)

		return found
	}

	if namedDef.lifetime == LifetimeScoped {
		found := getFromCacheOrBuildNoAlloc(c, c.scopedCache, key, name, namedDef, x)

		return found
	}

	buildItemNoAlloc(c, x, namedDef)

	return nil
}

func (c *Container) NewScope() *Container {
	var newScopedCache = lifetimeCache{
		entries: make(map[reflect.Type]map[string]*cacheEntry),
		mx:      sync.Mutex{},
	}

	return &Container{
		globalDef:      c.globalDef,
		scopedDef:      make(map[reflect.Type]map[string]*definition),
		singletonCache: c.singletonCache,
		scopedCache:    &newScopedCache,
	}
}

// We use a thread safe from getting items from the cache or build new ones
// There are 2 level caches one to lock the creation of containers for the type
// and the second one for the instance itself
// as it might need to resolved other dependencies and types.
func getFromCacheOrBuild[T any](
	c *Container,
	cache *lifetimeCache,
	key reflect.Type,
	name string,
	d *definition) (T, error) {
	typeCache, found := cache.entries[key]

	// If cache not found for the type create one
	if !found {
		cache.mx.Lock()

		typeCache, found = cache.entries[key]

		if !found {
			entry := make(map[string]*cacheEntry)
			cache.entries[key] = entry
			typeCache = entry
		}

		cache.mx.Unlock()
	}

	// Now we search for the named cache in the existing type cache
	namedCache, found := typeCache[name]

	if !found {
		cache.mx.Lock()

		namedCache, found = typeCache[name]

		if !found {
			entry := cacheEntry{}
			typeCache[name] = &entry
			namedCache = &entry
		}

		cache.mx.Unlock()
	}

	if !namedCache.initialized {
		namedCache.mx.Lock()

		if !namedCache.initialized {
			namedCache.initialized = true
			namedCache.instance = buildItem[T](c, d)
		}

		namedCache.mx.Unlock()
	}

	val, ok := namedCache.instance.(T)
	if !ok {
		panic("Couldn't cast the type")
	}

	return val, nil
}

// We use a thread safe from getting items from the cache or build new ones
// There are 2 level caches one to lock the creation of containers for the type
// and the second one for the instance itself
// as it might need to resolved other dependencies and types.
func getFromCacheOrBuildNoAlloc[T any](
	c *Container,
	cache *lifetimeCache,
	key reflect.Type,
	name string,
	d *definition,
	x *T) error {
	typeCache, found := cache.entries[key]

	// If cache not found for the type create one
	if !found {
		cache.mx.Lock()

		typeCache, found = cache.entries[key]

		if !found {
			entry := make(map[string]*cacheEntry)
			cache.entries[key] = entry
			typeCache = entry
		}

		cache.mx.Unlock()
	}

	// Now we search for the named cache in the existing type cache
	namedCache, found := typeCache[name]

	if !found {
		cache.mx.Lock()

		namedCache, found = typeCache[name]

		if !found {
			entry := cacheEntry{}
			typeCache[name] = &entry
			namedCache = &entry
		}

		cache.mx.Unlock()
	}

	if !namedCache.initialized {
		namedCache.mx.Lock()

		if !namedCache.initialized {
			namedCache.initialized = true
			buildItemNoAlloc(c, x, d)
			namedCache.instance = *x
		}

		namedCache.mx.Unlock()
	}

	val, ok := namedCache.instance.(T)
	*x = val
	if !ok {
		panic("Couldn't cast the type")
	}

	return nil
}

func buildItem[T any](c *Container, d *definition) T {
	factory, ok := d.f[0].(func(c *Container) T)
	if !ok {
		panic("factory doesn't match the expected format")
	}

	value := factory(c)

	for i := 1; i <= len(d.f)-1; i++ {
		decoratorF, ok := d.f[i].(func(decorated T, c *Container) T)
		if !ok {
			panic("Couldn't cast the decorator factory")
		}

		value = decoratorF(value, c)
	}

	return value
}

func buildItemNoAlloc[T any](c *Container, x *T, d *definition) {
	factory, ok := d.f[0].(func(c *Container) T)
	if !ok {
		panic("factory doesn't match the expected format")
	}

	val := factory(c)

	for i := 1; i <= len(d.f)-1; i++ {
		decoratorF, ok := d.f[i].(func(decorated T, c *Container) T)
		if !ok {
			panic("Couldn't cast the decorator factory")
		}

		val = decoratorF(val, c)
	}
	*x = val
}

func getKeyFromT[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil))
}
