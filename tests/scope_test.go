package test

import (
	"strconv"
	"testing"

	"github.com/mingue/godi"
)

type (
	SomeInterface interface{}
	SomeStruct    struct {
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

func TestGetTransientNamedInstance(t *testing.T) {
	var cont = godi.New()
	name := "name"

	godi.TransientNamed(cont, name, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	x, err := godi.GetNamed[SomeInterface](cont, name)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if x == nil {
		t.Fatalf("x should not be nil")
	}
}

func TestGetSingletonNamedInstance(t *testing.T) {
	var cont = godi.New()
	name := "name"
	godi.SingletonNamed(cont, name, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	x, err := godi.GetNamed[SomeInterface](cont, name)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if x == nil {
		t.Fatalf("x should not be nil")
	}
}

func TestGetScopedNamedInstance(t *testing.T) {
	var cont = godi.New()
	name := "name"
	godi.ScopedNamed(cont, name, func(c *godi.Container) SomeInterface {
		return &SomeStruct{}
	})

	x, err := godi.GetNamed[SomeInterface](cont, name)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if x == nil {
		t.Fatalf("x should not be nil")
	}
}

func TestGetUnnamedAndNamedForSameType(t *testing.T) {
	var cont = godi.New()
	name := "name"

	godi.Singleton(cont, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})
	godi.SingletonNamed(cont, name, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	unnamed, err := godi.Get[*SomeStruct](cont)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if unnamed == nil {
		t.Fatalf("y should not be nil")
	}

	named, err := godi.GetNamed[*SomeStruct](cont, name)
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if named == nil {
		t.Fatalf("x should not be nil")
	}

	unnamed.data = "some"

	if unnamed.data == named.data {
		t.Fatalf("It should be a different instance")
	}
}

func TestGetDifferentNamedInstances(t *testing.T) {
	var cont = godi.New()

	godi.SingletonNamed(cont, "1", func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})
	godi.SingletonNamed(cont, "2", func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	named1, err := godi.GetNamed[*SomeStruct](cont, "1")
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if named1 == nil {
		t.Fatalf("named1 should not be nil")
	}

	named2, err := godi.GetNamed[*SomeStruct](cont, "2")
	if err != nil {
		t.Fatalf("Failed to get instance: %v", err.Error())
	}

	if named2 == nil {
		t.Fatalf("named2 should not be nil")
	}

	named1.data = "some"

	if named1.data == named2.data {
		t.Fatalf("It should be a different instance")
	}
}

func TestGetScopedReturnDifferentInstanceOnDifferentScopeForNamedDefinitions(t *testing.T) {
	var cont = godi.New()
	godi.ScopedNamed(cont, "1", func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})
	godi.ScopedNamed(cont, "2", func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	for i := 1; i < 3; i++ {
		x, _ := godi.GetNamed[*SomeStruct](cont, strconv.Itoa(i))
		x.data = "x"

		scopedCont := cont.NewScope()
		y, _ := godi.GetNamed[*SomeStruct](scopedCont, strconv.Itoa(i))

		if y.data != "" {
			t.Fatalf("y.data should be empty")
		}
	}
}

func TestSameDefinitionReturnErrorIfRegisteredInGlobalAndScopedContainer(t *testing.T) {
	var cont = godi.New()
	godi.Scoped(cont, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	var secondScope = cont.NewScope()
	err := godi.Scoped(secondScope, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	if err == nil {
		t.Fatal("If definition already registered in global container it should not allow to register in scoped container")
	}
}

func TestSameDefinitionCanBeRegisteredOncePerScopedContainer(t *testing.T) {
	var cont = godi.New()
	firstScope := cont.NewScope()
	godi.Scoped(firstScope, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	var secondScope = cont.NewScope()
	err := godi.Scoped(secondScope, func(c *godi.Container) *SomeStruct {
		return &SomeStruct{}
	})

	if err != nil {
		t.Fatal("It should allow to register instance on new scope")
	}
}
