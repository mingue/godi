package test

import (
	"log"
	"testing"

	"github.com/mingue/godi"
)

type Doer interface {
	Do()
}

type (
	SimpleDoer struct {
		count int
	}
	CallCountDecorator struct {
		d         Doer
		preCount  uint8
		postCount uint8
	}
)

func (d *SimpleDoer) Do() {
	d.count++
}

func (d *CallCountDecorator) Do() {
	d.preCount++
	d.d.Do()
	d.postCount++
}

func TestDecorateErrorsIfFactoryNotRegisteredPreviously(t *testing.T) {
	var cont = godi.New()

	err := godi.Decorate(cont, func(d Doer, c *godi.Container) Doer {
		return &CallCountDecorator{
			d: d,
		}
	})

	if err == nil || err.Error() != godi.ErrDecoratorBeforeFactory {
		log.Fatal("Expecting factory first error")
	}
}

func TestDecorateOnlyAcceptsInterfaces(t *testing.T) {
	var cont = godi.New()

	godi.Transient(cont, func(c *godi.Container) *SimpleDoer {
		return &SimpleDoer{}
	})

	err := godi.Decorate(cont, func(d *SimpleDoer, c *godi.Container) *SimpleDoer {
		return &SimpleDoer{}
	})

	if err == nil || err.Error() != godi.ErrDecoratedMustBeInterface {
		log.Fatal("Validation failed")
	}
}

func TestDecorateATransient(t *testing.T) {
	var cont = godi.New()
	godi.Transient(cont, func(c *godi.Container) Doer {
		return &SimpleDoer{}
	})

	err := godi.Decorate(cont, func(d Doer, c *godi.Container) Doer {
		return &CallCountDecorator{
			d: d,
		}
	})

	if err != nil {
		log.Fatalf("Couldn't decorate: %v", err.Error())
	}

	x, _ := godi.Get[Doer](cont)
	x.Do()

	decorator := x.(*CallCountDecorator)

	if decorator.preCount != 1 {
		log.Fatal("Precount should be 1")
	}

	if decorator.postCount != 1 {
		log.Fatal("Postcount should be 1")
	}

	decorated := decorator.d.(*SimpleDoer)

	if decorated.count != 1 {
		log.Fatal("Count should be 1")
	}
}
