package main

import (
	"strconv"
	"testing"
)

//===========================================================================

//TestCheckRegx_panic1() ....
func TestCheckRegx_panic1(t *testing.T) {
	pmsg := "panic: runtime error: index out of range [2] with length 2"
	got := CheckRegx(regexPanic, pmsg)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckRegx_panic1() ....
func TestCheckRegx_panic2(t *testing.T) {
	pmsg := "When I saw it I almost had a panic"
	got := CheckRegx(regexPanic, pmsg)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

// ?    github.com/zchee/nvim-go/pkg/server [no test files]
//TestCheckRegx_no_test_files1() ....
func TestCheckRegx_no_test_files1(t *testing.T) {
	pmsg := "?    github.com/zchee/nvim-go/pkg/server [no test files]"
	got := CheckRegx(regexNoTestFiles, pmsg)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckRegx_no_test_files2_capital_N() ....
func TestCheckRegx_no_test_files2_capital_N(t *testing.T) {
	pmsg := "?    github.com/zchee/nvim-go/pkg/server [No test files]"
	got := CheckRegx(regexNoTestFiles, pmsg)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================

// func BuildBarMessage(runs int, skips int, fails int, passes int, elapsed PD_Elapsed, fname string, lineno string, coverage string) string {

//TestBuildBarMessage_no_fails_no_skips() ....
func TestBuildBarMessage_no_fails_no_skips(t *testing.T) {
	got := BuildBarMessage(10, 0, 0, 10, 0.013, "", "", "1.4%")
	want := "10 Run, 10 Passed, Test Coverage: 1.4%, in 0.013s"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestBuildBarMessage_no_fails_but_skips() ....
func TestBuildBarMessage_no_fails_but_skips(t *testing.T) {
	got := BuildBarMessage(10, 4, 0, 6, 0.015, "", "", "7.4%")
	want := "10 Run, 6 Passed, 4 Skipped, in 0.015s"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestBuildBarMessage_fails_no_skips() ....
func TestBuildBarMessage_fails_no_skips(t *testing.T) {
	got := BuildBarMessage(10, 0, 2, 8, 0.015, "main_test.go", "37", "8.4%")
	want := "10 Run, 8 Passed, 2 Failed, 1st in main_test.go, on line 37, in 0.015s"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
//===========================================================================
