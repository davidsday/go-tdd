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

//TestStdErrLongerThanScreenWidth ....
func TestStdErrMsgTooLongForOneLine_144_cols(t *testing.T) {
	msg := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	got := stdErrMsgTooLongForOneLine(msg, stdErrMsgTrailer, 144)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestStdErrMsgTooLongForOneLine_80_cols ...
func TestStdErrMsgTooLongForOneLine_80_cols(t *testing.T) {
	msg := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	// msg := "xxxxxxxxxx"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	got := stdErrMsgTooLongForOneLine(msg, stdErrMsgTrailer, 80)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestStdErrMsgTooLongForOneLine_80_cols ...
func TestStdErrMsgTooLongForOneLine_80_cols_short_msg(t *testing.T) {
	msg := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	got := stdErrMsgTooLongForOneLine(msg, stdErrMsgTrailer, 80)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================

//TestgetAverageCyclomaticComplexity ....
func TestGetAverageCyclomaticComplexity(t *testing.T) {
	got, err := getAvgCyclomaticComplexity("./tests/avgCCmplx/main.go")
	if err != nil {
		t.Errorf("Error in get AverageCyclomaticComplexity()")
	}
	want := "7.29"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func rcvdMsgOnStdErr(stderror string) bool
//===========================================================================

//TestRcvdMsgOnStdErr_no ....
func TestRcvdMsgOnStdErr_no(t *testing.T) {
	got := rcvdMsgOnStdErr("")
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestRcvdMsgOnStdErr_yes ....
func TestRcvdMsgOnStdErr_yes(t *testing.T) {
	got := rcvdMsgOnStdErr("Pretend STDERR message")
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestRcvdMsgOnStdErr_yes_one_char ....
func TestRcvdMsgOnStdErr_yes_one_char(t *testing.T) {
	got := rcvdMsgOnStdErr("P")
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// func skipMsg(skips int) string {
//===========================================================================

//TestSkipMsg_0 ....
func TestSkipMsg_0(t *testing.T) {
	got := skipMsg(0)
	want := ""
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestSkipMsg_3 ....
func TestSkipMsg_3(t *testing.T) {
	got := skipMsg(3)
	want := ", 3 skips"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
