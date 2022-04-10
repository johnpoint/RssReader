package bootstrap

import (
	"context"
	"testing"
)

func TestHelper_AddComponent(t *testing.T) {
	var h Helper
	h.AddComponent(&EmptyComponent{})

	if len(h.components) != 1 {
		t.Fail()
	}
}

func TestHelper_Init(t *testing.T) {
	var h Helper
	h.AddComponent(&EmptyComponent{})
	err := h.Init(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestHelper_Init_Failed(t *testing.T) {
	var h Helper
	h.AddComponent(&EmptyComponent{error: true})
	err := h.Init(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmptyComponent_Init(t *testing.T) {
	var c EmptyComponent
	err := c.Init(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
