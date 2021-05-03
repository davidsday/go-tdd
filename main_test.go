package main

import (
	"reflect"
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
	got := getAvgCyclomaticComplexity("./tests/avgCCmplx/main.go")
	want := "7.29"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestgetAverageCyclomaticComplexity ....
func TestGetAverageCyclomaticComplexity_no_go_files(t *testing.T) {
	got := getAvgCyclomaticComplexity("./bin")
	want := "NaN"
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
	want := ", 3 Skipped"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func failMsg(skips int) string {
//===========================================================================

// func failMsg(fails int, fname, lineno string) string {

//TestFailMsg with fails....
func TestFailMsgWithFails(t *testing.T) {
	got := failMsg(4, "main_test.go", "87")
	want := ", 4 Failed, 1st in main_test.go, on line 87"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestFailMsg without fails....
func TestFailMsgWithOutFails(t *testing.T) {
	got := failMsg(0, "", "")
	want := ""
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func passMsg(passes int) string {
//===========================================================================

//TestPassMsg_10 ....
func TestPassMsg_10(t *testing.T) {
	got := passMsg(10)
	want := ", 10 Passed"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestPassMsg_0 ....
func TestPassMsg_0(t *testing.T) {
	got := passMsg(0)
	want := ", 0 Passed"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func runMsg(runs int) string {
//===========================================================================

//TestRunMsg_10 ....
func TestRunMsg_10(t *testing.T) {
	got := runMsg(10)
	want := "10 Run"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestRunMsg_0 ....
func TestRunMsg_0(t *testing.T) {
	got := runMsg(0)
	want := "0 Run"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func elapsedMsg(elapsed PDElapsed) string
//===========================================================================

//TestElapsedMsg_0.005 ....
func TestElapsedMsg_0_005(t *testing.T) {
	got := elapsedMsg(0.005)
	want := ", in 0.005s"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestElapsedMsg_0.005 ....
func TestElapsedMsg_0_000(t *testing.T) {
	got := elapsedMsg(0.000)
	want := ", in 0.000s"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func convertStringToBytes(s string) []byte {
//===========================================================================

//TestConvertStringToBytes() ....
func TestConvertStringToBytes(t *testing.T) {
	s := "What we want"
	got := convertStringToBytes(s)
	want := []byte("What we want")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%s' want '%s'", string(got), string(want))
	}
}

//===========================================================================
// func splitBytesIntoLines(b []byte) [][]byte {
//===========================================================================

//TestSplitBytesIntoLines ....
func TestSplitBytesIntoLines(t *testing.T) {
	b := []byte("Strings\nIntegers\nFloats\nBooleans\n")
	got := splitBytesIntoLines(b)
	want := [][]byte{[]byte("Strings"), []byte("Integers"), []byte("Floats"), []byte("Booleans")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestSplitBytesIntoLines_with_non_empty_last_line ....
func TestSplitBytesIntoLines_with_non_empty_last_line(t *testing.T) {
	b := []byte("Strings\nIntegers\nFloats\nBooleans\nBytes")
	got := splitBytesIntoLines(b)
	want := [][]byte{[]byte("Strings"), []byte("Integers"), []byte("Floats"), []byte("Booleans"), []byte("Bytes")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func thisIsTheFirstFailure(pgmdata PgmData) bool {
//===========================================================================

//TestThisIsTheFirstFailure_true ....
func TestThisIsTheFirstFailure_true(t *testing.T) {
	var pgmdata PgmData
	pgmdata.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 0, "fail": 0, "output": 0}
	pgmdata.Counts["fail"] = 0
	got := thisIsTheFirstFailure(&pgmdata)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestThisIsTheFirstFailure_false ....
func TestThisIsTheFirstFailure_false(t *testing.T) {
	var pgmdata PgmData
	pgmdata.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 0, "fail": 0, "output": 0}
	pgmdata.Counts["fail"] = 1
	got := thisIsTheFirstFailure(&pgmdata)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// takeNoteOfFirstFailure(pgmdata PgmData, parts []string, testName string)
//===========================================================================

//TestTakeNoteOfFirstFailure ....
func TestTakeNoteOfFirstFailure(t *testing.T) {
	pd := PgmData{}
	testName := "thisTest"
	parts := []string{"firstPart", "secondPart"}

	takeNoteOfFirstFailure(&pd, parts, testName)

	if pd.Firstfailedtest.Tname != testName && pd.Firstfailedtest.Lineno != parts[1] && pd.Firstfailedtest.Fname != parts[0] {
		t.Errorf("Filename: %s, LineNo: %s, TestName: %s", pd.Firstfailedtest.Fname, pd.Firstfailedtest.Lineno, pd.Firstfailedtest.Tname)
	}
}

//TestAddToQuickFixList ....
func TestAddToQuickFixList(t *testing.T) {
	var pgmdata PgmData
	parts := []string{"firstPart", "secondPart", "thirdPart"}
	args := []string{"programName", "packageDir", "10"}
	jlo := JLObject{}
	jlo.Test = "thisTest"

	addToQuickFixList(&pgmdata, args, parts, jlo)
	if len(pgmdata.QfList) <= 0 {
		t.Errorf("The length of pgmdata.QfList is  '%d'\n", len(pgmdata.QfList))
	}
	if pgmdata.QfList[0].Filename != "packageDir/firstPart" {
		t.Errorf("Filename:  Got: %s, Wanted: %s\n", pgmdata.QfList[0].Filename, "packageDir/firstPart")
	}
}

//TestUnneededFAILPrefix_Has_FAIL ....
func TestUnneededFAILPrefix_Has_FAIL(t *testing.T) {
	output := "FAIL:Part1:Part2:Part3"
	got := removeUnneededFAILPrefix(output)
	want := []string{"Part1", "Part2", "Part3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}

//TestUnneededFAILPrefix_Has_No_FAIL ....
func TestUnneededFAILPrefix_Has_No_FAIL(t *testing.T) {
	output := "Part1:Part2:Part3"
	got := removeUnneededFAILPrefix(output)
	want := []string{"Part1", "Part2", "Part3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}
