package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

var results GtpResults

func TestMain(m *testing.M) {
	//2021/06/02 10:36:52 results.Args: 'main.GtpArgs{PackageDir:"/home/dave/sw/go/go-tdd", ScreenColumns:"144", GocycloIgnore:"vendor|testdata", GoTddDebug:true, PluginDir:"/home/dave/sw/go/go-tdd", Timeout:"10s"}'
	setupLogging()
	results = newResults()
	// for testing, we can have PackageDir and PluginDir the same and since we are running the tests in the PackageDir, they can both be "."
	// for testing purposes.
	results.Args = GtpArgs{
		PackageDir:    ".",
		ScreenColumns: "144",
		GocycloIgnore: "vendor|testdata",
		GoTddDebug:    true,
		PluginDir:     ".",
		Timeout:       "10s",
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

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
	exampleFuncDecl := `func ExampleHW`

	got1, _, _ := findExampleFunc(results.Args.PluginDir, exampleFuncDecl, results.Args.PackageDir, results.Args.GocycloIgnore)

	want := `examples_test.go`
	if got1 != want {
		t.Errorf("got '%s' want '%s'", got1, want)
	}
}

//TestFindExampleFunc ....
func TestFindExampleFunc_XXXX(t *testing.T) {
	exampleFuncDecl := `func ExampleTestXXXX\(\) {`

	got1, _, _ := findExampleFunc(results.Args.PluginDir, exampleFuncDecl, results.Args.PackageDir, results.Args.GocycloIgnore)

	want := `examples_test.go`
	if got1 != want {
		t.Errorf("got '%s' want '%s'", got1, want)
	}
}
