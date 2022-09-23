package godi

import (
	"reflect"
	"sync"
)

type lifetime string

const (
	LifetimeSingleton lifetime = "Singleton"
	LifetimeScoped    lifetime = "Scoped"
	LifetimeTransient lifetime = "Transient"
)

type lifetimeCache struct {
	entries map[reflect.Type]*cacheEntry
	mx      sync.Mutex
}

type cacheEntry struct {
	mx          sync.Mutex
	initialized bool
	instance    any
}