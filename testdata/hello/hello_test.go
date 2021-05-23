package main

import "testing"

//TestHello ....
func TestHello(t *testing.T) {
	got := hello()
	want := "hello World"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
