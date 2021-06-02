package main

import (
	"fmt"
	"strconv"
	"testing"
)

// Output: What's shakin
func ExampleTestXXXX() {
	fmt.Println("Hello, World")
	// Output: Hello, World
}

//===========================================================================
//TestExampleError
//===========================================================================

//TestExampleError ....
func TestExampleError_true(t *testing.T) {
	input := `{"Time":"2021-05-27T10:05:48.703313416-04:00","Action":"output","Package":"example","Test":"ExampleHW","Output":"--- FAIL: ExampleHW (0.00s)\n"}`
	jlo := JLObject{}
	jlo.unmarshal(input)

	got := exampleError(jlo.getOutput())
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestExampleError_false ....
func TestExampleError_false(t *testing.T) {
	input := `{"Time":"2021-05-27T10:05:48.703313416-04:00","Action":"output","Package":"example","Test":"ExampleHW","Output":"PASS: ExampleHW (0.00s)\n"}`
	jlo := JLObject{}
	jlo.unmarshal(input)

	got := exampleError(jlo.getOutput())
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
//TestFindExample
//===========================================================================

//TestFindExampleFunc ....
func TestFindExampleFunc(t *testing.T) {
	if debug {
		setupLogging()
	}
	exampleFuncDecl := `func ExampleHW`
	// plugDir & pkgDir the same for testing
	plugDir := `/home/dave/sw/go/go-tdd`
	pkgDir := `/home/dave/sw/go/go-tdd`

	results := newResults()
	results.Args.GocycloIgnore = `vendor|testdata`
	got1, _, _ := findExampleFunc(plugDir, exampleFuncDecl, pkgDir, results.Args.GocycloIgnore)

	want := `/home/dave/sw/go/go-tdd/examples_test.go`
	if got1 != want {
		t.Errorf("got '%s' want '%s'", got1, want)
	}
}

//TestFindExampleFunc ....
func TestFindExampleFunc_XXXX(t *testing.T) {
	if debug {
		setupLogging()
	}
	exampleFuncDecl := `func ExampleTestXXXX\(\) {`
	// plugDir & pkgDir the same for testing
	plugDir := `/home/dave/sw/go/go-tdd`
	pkgDir := `/home/dave/sw/go/go-tdd`
	results := newResults()
	results.Args.GocycloIgnore = `vendor|testdata`
	got1, _, _ := findExampleFunc(plugDir, exampleFuncDecl, pkgDir, results.Args.GocycloIgnore)

	want := `/home/dave/sw/go/go-tdd/examples_test.go`
	if got1 != want {
		t.Errorf("got '%s' want '%s'", got1, want)
	}
}
