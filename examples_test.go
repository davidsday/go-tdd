package main

import (
	"fmt"
	"reflect"
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
	exampleFuncDecl := `func ExampleHW`
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

//===========================================================================
//TestSplitExampleFuncSearchResults
//===========================================================================

func TestSplitExampleFuncSearchResults(t *testing.T) {
	results := `example_test.go:7:1:func ExampleHW() {`
	want := "example_test.go"
	got1, _, _ := splitExampleFuncSearchResults(results)

	if !reflect.DeepEqual(got1, want) {
		t.Errorf("Got: '%s', Want: '%s'", got1, want)
	}
}
