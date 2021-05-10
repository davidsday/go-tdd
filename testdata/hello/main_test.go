package main

import (
	"testing"

	"github.com/davidsday/hello/lib"
)

func TestHello(t *testing.T) {
	want := "Hello, World!"
	if got := lib.Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
