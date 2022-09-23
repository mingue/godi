package test

import (
	"testing"

	"github.com/mingue/godi"
)

type (
	SomeInterface interface {}
	SomeStruct struct {
		data string
	}
)

func TestGetTransientInstance(t *testing.T) {
	var cont = godi.New()
	godi.Transient(cont, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	x, err := godi.Get[SomeInterface](cont)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if x == nil {
		t.Fatalf("x should not be nil")
	}
}

func TestGetSingletonInstance(t *testing.T) {
	var cont = godi.New()
	godi.Singleton(cont, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	x, err := godi.Get[SomeInterface](cont)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if x == nil {
		t.Fatalf("x should not be nil")
	}
}

func TestGetScopedInstance(t *testing.T) {
	var cont = godi.New()
	godi.Scoped(cont, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	x, err := godi.Get[SomeInterface](cont)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if x == nil {
		t.Fatalf("x should not be nil")
	}
}

func TestGetTransientReturnNewInstancesEveryTime(t *testing.T) {
	var cont = godi.New()
	godi.Transient(cont, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	x, _ := godi.Get[*SomeStruct](cont)
	y, _ := godi.Get[*SomeStruct](cont)

	x.data = "x"

	if y.data != "" {
		t.Fatalf("y.data should be empty")
	}
}

func TestGetSingletonReturnSameInstanceEveryTime(t *testing.T) {
	var cont = godi.New()
	godi.Singleton(cont, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	x, _ := godi.Get[*SomeStruct](cont)
	y, _ := godi.Get[*SomeStruct](cont)

	x.data = "x"

	if y.data != "x" {
		t.Fatalf("y.data should be x")
	}
}

func TestGetSingletonReturnSameInstanceEvenOnNewScopeEveryTime(t *testing.T) {
	var cont = godi.New()
	godi.Singleton(cont, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	x, _ := godi.Get[*SomeStruct](cont)
	x.data = "x"

	scopedCont := cont.NewScope()
	y, _ := godi.Get[*SomeStruct](scopedCont)

	if y.data != "x" {
		t.Fatalf("y.data should be x")
	}
}

func TestGetScopedReturnSameInstanceOnSameContainerScope(t *testing.T) {
	var cont = godi.New()
	godi.Scoped(cont, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	x, _ := godi.Get[*SomeStruct](cont)
	y, _ := godi.Get[*SomeStruct](cont)

	x.data = "x"

	if y.data != "x" {
		t.Fatalf("y.data should be x")
	}
}

func TestGetScopedReturnDifferentInstanceOnDifferentScope(t *testing.T) {
	var cont = godi.New()
	godi.Scoped(cont, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	x, _ := godi.Get[*SomeStruct](cont)
	x.data = "x"

	scopedCont := cont.NewScope()
	y, _ := godi.Get[*SomeStruct](scopedCont)

	if y.data != "" {
		t.Fatalf("y.data should be empty")
	}
}
