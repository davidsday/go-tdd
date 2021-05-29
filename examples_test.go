package main

import (
	"fmt"
	"strconv"
	"testing"
)

// Output: Hello, World
func ExampleTestXXXX() {
	fmt.Println("Hello, World")
	// Output: What's shakin
}

//===========================================================================
//TestExampleError
//===========================================================================

//TestExampleError ....
func TestExampleError(t *testing.T) {
	input := `{"Time":"2021-05-27T10:05:48.703313416-04:00","Action":"output","Package":"example","Test":"ExampleHW","Output":"--- FAIL: ExampleHW (0.00s)\n"}`
	jlo := JLObject{}
	jlo.unmarshal(input)

	got := exampleError(jlo.getOutput())
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
//TestFindExample
//===========================================================================

//TestFindExampleFunc ....
func TestFindExampleFunc(t *testing.T) {
	setupLogging()
	exampleFuncDecl := `func ExampleHW\\(\\) {`
	plugDir := `/home/dave/.config/nvim/plugged/go-tdd`
	pkgDir := `/home/dave/sw/go/go-tdd`
	results := newResults()
	results.GocycloIgnore = `vendor|testdata`
	got1, _, _ := findExampleFunc(plugDir, exampleFuncDecl, pkgDir, results.GocycloIgnore)

	want := `/home/dave/sw/go/go-tdd/examples_test.go`
	if got1 != want {
		t.Errorf("got '%s' want '%s'", got1, want)
	}
}
