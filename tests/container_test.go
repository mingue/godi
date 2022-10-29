package test

import (
	"log"
	"testing"

	"github.com/mingue/godi"
)

func TestGetErrorIfNotRegistered(t *testing.T) {
	var cont = godi.New()
	_, err := godi.Get[SomeInterface](cont)

	if err == nil || err.Error() != godi.ErrFactoryNotRegistered {
		log.Fatal("Expecting factory not registered")
	}
}

func TestGetErrorIfAlreadyRegistered(t *testing.T) {
	var cont = godi.New()
	godi.Transient(cont, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	err := godi.Transient(cont, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	if err == nil || err.Error() != godi.ErrFactoryAlreadyRegistered {
		log.Fatal("Expecting factory not registered")
	}
}
