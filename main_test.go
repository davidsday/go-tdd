package main

import (
	"encoding/json"
	"os"
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
	stdErrMsgPrefix := "STDERR:"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	got := stdErrMsgTooLongForOneLine(msg, stdErrMsgPrefix, stdErrMsgTrailer, 144)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestStdErrMsgTooLongForOneLine_80_cols ...
func TestStdErrMsgTooLongForOneLine_80_cols(t *testing.T) {
	msg := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	// msg := "xxxxxxxxxx"
	stdErrMsgPrefix := "STDERR:"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	got := stdErrMsgTooLongForOneLine(msg, stdErrMsgPrefix, stdErrMsgTrailer, 80)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestStdErrMsgTooLongForOneLine_80_cols ...
func TestStdErrMsgTooLongForOneLine_80_cols_short_msg(t *testing.T) {
	msg := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	stdErrMsgPrefix := "STDERR:"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	got := stdErrMsgTooLongForOneLine(msg, stdErrMsgPrefix, stdErrMsgTrailer, 80)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
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
	var Results GtpResults
	Results.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 0, "fail": 0, "output": 0}
	Results.Counts["fail"] = 0
	got := thisIsTheFirstFailure(&Results)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestThisIsTheFirstFailure_false ....
func TestThisIsTheFirstFailure_false(t *testing.T) {
	var Results GtpResults
	Results.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 0, "fail": 0, "output": 0}
	Results.Counts["fail"] = 1
	got := thisIsTheFirstFailure(&Results)
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
	Results := GtpResults{}
	testName := "thisTest"
	parts := []string{"firstPart", "secondPart"}

	takeNoteOfFirstFailure(&Results, parts, testName)

	if Results.FirstFail.Tname != testName && Results.FirstFail.Lineno != parts[1] && Results.FirstFail.Fname != parts[0] {
		t.Errorf("Filename: %s, LineNo: %s, TestName: %s", Results.FirstFail.Fname, Results.FirstFail.Lineno, Results.FirstFail.Tname)
	}
}

//TestUnneededFAILPrefix_Has_FAIL ....
func TestUnneededFAILPrefix_Has_FAIL(t *testing.T) {
	output := "FAIL:Part1:Part2:Part3"
	list := splitOnSemiColons(output)
	got := removeUnneededFAILPrefix(list)
	want := []string{"Part1", "Part2", "Part3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}

//TestUnneededFAILPrefix_Has_No_FAIL ....
func TestUnneededFAILPrefix_Has_No_FAIL(t *testing.T) {
	output := "Part1:Part2:Part3"
	list := splitOnSemiColons(output)
	got := removeUnneededFAILPrefix(list)
	want := []string{"Part1", "Part2", "Part3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}

//TestDoStdErrMsg ....
func TestDoStdErrMsg(t *testing.T) {
	Results := GtpResults{}
	Results.VimColumns = 135
	msg := "STDERR: This is my message from STDERR."
	PackageDir := "/home/dave/sw/go/goTestParser"
	doStdErrMsg(msg, &Results, PackageDir)
	if len(Results.Errors) == 0 {
		t.Errorf("Length of pgmdata.Perrors = %s\n", strconv.Itoa(len(Results.Errors)))
	}
}

//TestDoStdErrMsgTooLong ....
func TestDoStdErrMsgTooLong(t *testing.T) {
	Results := GtpResults{}
	Results.VimColumns = 135
	msg := "STDERR: This is my message from STDERR. xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	PackageDir := "/home/dave/sw/go/goTestParser"
	doStdErrMsg(msg, &Results, PackageDir)
	if len(Results.Errors) == 0 {
		t.Errorf("Length of pgmdata.Perrors = %s\n", strconv.Itoa(len(Results.Errors)))
	}
	_ = os.Remove(PackageDir + "/StdErr.txt")
}

//===========================================================================
// metricsMsg()
//===========================================================================

//TestMetricsMsg() ....
func TestMetricsMsg(t *testing.T) {
	skips := 0
	fails := 0
	coverage := "32.0%"
	complexity := "7.29"
	got := metricsMsg(skips, fails, coverage, complexity)
	want := ", Test Coverage: 32.0%, Average Complexity: 7.29"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestMetricsMsg() ....
func TestMetricsMsg_no_skips_2_fails(t *testing.T) {
	skips := 0
	fails := 2
	coverage := "32.0%"
	complexity := "7.29"
	got := metricsMsg(skips, fails, coverage, complexity)
	want := ""
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// buildAndAppendAnErrorForInvalidJSON(&PD)
//===========================================================================

//TestBuildAndAppendAnErrorForInvalidJSON ....
func TestBuildAndAppendAnErrorForInvalidJSON(t *testing.T) {
	Results := GtpResults{}
	Results.VimColumns = 135
	buildAndAppendAnErrorForInvalidJSON(&Results)
	if len(Results.Errors) <= 0 {
		t.Errorf("pgmdata.Perrors has %d elements", len(Results.Errors))
	}
}

//===========================================================================
// func ifFinalActionWasPass(jlo.Action string) bool
//===========================================================================

//TestIfFinalActionWasPass ....
func TestFinalActionWasPass_pass(t *testing.T) {
	got := finalActionWasPass("pass")
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestIfFinalActionWasPass ....
func TestFinalActionWasPass_fail(t *testing.T) {
	got := finalActionWasPass("fail")
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// func ifFinalActionWasFail(jlo.Action string) bool
//===========================================================================

//TestIfFinalActionWasFail ....
func TestFinalActionWasFail_fail(t *testing.T) {
	got := finalActionWasFail("fail")
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestIfFinalActionWasFail_pass ....
func TestFinalActionWasFail_pass(t *testing.T) {
	got := finalActionWasFail("pass")
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// func WeHaveHadMoreThanOnePass(passes int) bool
//===========================================================================

//Test ....
func TestWeHaveHadMoreThanOnePass_five(t *testing.T) {
	passes := 5
	got := weHaveHadMoreThanOnePass(passes)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

func TestWeHaveHadMoreThanOnePass_zero(t *testing.T) {
	passes := 0
	got := weHaveHadMoreThanOnePass(passes)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// func WeHaveHadMoreThanOneFail(fails int) bool
//===========================================================================

//Test ....
func TestWeHaveHadMoreThanOneFail_five(t *testing.T) {
	fails := 5
	got := weHaveHadMoreThanOneFail(fails)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

func TestWeHaveHadMoreThanOneFail_zero(t *testing.T) {
	fails := 0
	got := weHaveHadMoreThanOneFail(fails)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestSplitOnSemicolons(str string) []string ....
func TestSplitOnSemicolons(t *testing.T) {
	str := "FAIL:1st useful part:2nd useful part:3rd useful part"
	got := splitOnSemiColons(str)
	want := []string{"FAIL", "1st useful part", "2nd useful part", "3rd useful part"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}

//TestShellOut ....
func TestShellOut(t *testing.T) {
	got, _, _ := Shellout("cat ./testdata/hello.txt")
	want := "Hello World!\n"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func hasTestFileReferences(output string) bool {
//===========================================================================

////TesthasTestFileReferences_true ....
func TestHasTestFileReferences_true(t *testing.T) {
	output := "Part1:main_test.go:12:Part4"
	got := hasTestFileReferences(output)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TesthasTestFileReferences_none ....
func TestHasTestFileReferences_none(t *testing.T) {
	output := "Part0:Part1:Part2:Part3"
	got := hasTestFileReferences(output)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// func hasTestCoverage(output string) bool {
//===========================================================================

//TestHasTestCoverage ....
func TestHasTestCoverage_yes(t *testing.T) {
	output := "coverage: 69.3% of statements"
	got := hasTestCoverage(output)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestHasTestCoverage ....
func TestHasTestCoverage_no(t *testing.T) {
	output := "FAIL:Part1:Part2:Part3"
	got := hasTestCoverage(output)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// func checkErrorCandidates(Results *GtpResults, output string) bool {
//===========================================================================

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_no_test_files(t *testing.T) {
	results := GtpResults{}
	// output := "FAIL:Part1:Part2:Part3"
	output := "[no test files]"
	got := checkErrorCandidates(&results, output)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_yes(t *testing.T) {
	results := GtpResults{}
	// output := "FAIL:Part1:Part2:Part3"
	output := "panic:"
	got := checkErrorCandidates(&results, output)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_no(t *testing.T) {
	results := GtpResults{}
	output := "Part0:Part1:Part2:Part3"
	got := checkErrorCandidates(&results, output)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//===========================================================================
// func adjustOutSuperfluousFinalPass(action string, passCount int) int {
//===========================================================================

//TestAdjustOutSuperfluousFinalPass ....
func TestAdjustOutSuperfluousFinalPass(t *testing.T) {
	passCount := 21
	got := adjustOutSuperfluousFinalPass("pass", passCount)
	want := 20
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//TestAdjustOutSuperfluousFinalPass ....
func TestAdjustOutSuperfluousFinalPass_1(t *testing.T) {
	passCount := 1
	got := adjustOutSuperfluousFinalPass("pass", passCount)
	want := 1
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//===========================================================================
// func adjustOutSuperfluousFinalFail(action string, passCount int) int {
//===========================================================================

//TestAdjustOutSuperfluousFinalPass ....
func TestAdjustOutSuperfluousFinalFail(t *testing.T) {
	failCount := 21
	got := adjustOutSuperfluousFinalFail("fail", failCount)
	want := 20
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//TestAdjustOutSuperfluousFinalPass ....
func TestAdjustOutSuperfluousFinalFail_1(t *testing.T) {
	failCount := 1
	got := adjustOutSuperfluousFinalFail("fail", failCount)
	want := 1
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//===========================================================================
// func adjustOutSuperfluousFinalResult(action string, Results *GtpResults) (int, int) {
//===========================================================================

//TestAdjustOutSuperfuousFinalResult ....
func TestAdjustOutSuperfuousFinalResult(t *testing.T) {
	Results := GtpResults{}
	Results.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 21, "fail": 0, "output": 0}
	action := "pass"

	passes, fails := adjustOutSuperfluousFinalResult(action, &Results)

	if passes != 20 || fails != 0 {
		t.Errorf("passes:'%d', fails: '%d', wanted: passes: '20', fails: 0", passes, fails)
	}
}

//===========================================================================
// func HandleOutputLines(Results *GtpResults, jlo JLObject, prevJlo JLObject,
//	PackageDirFromVim string) (bool, error) {
//===========================================================================

//TestHandleOutputLines ....
func TestHandleOutputLines(t *testing.T) {
	Results := GtpResults{}
	Results.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 21, "fail": 0, "output": 0}
	jlo, prevJlo := JLObject{}, JLObject{}
	jsonlinePrevJlo := []byte(`{"Time": "2021-05-07T23:32:18.412171038-04:00", "Action": "output", "Package": "github.com/davidsday/goTestParser", "Output": "PASS\n"}`)
	jsonlineJlo := []byte(`{"Time": "2021-05-07T23:32:18.412174016-04:00", "Action": "output", "Package": "github.com/davidsday/goTestParser", "Output": "coverage: 77.9% of statements\n"}`)
	err := json.Unmarshal(jsonlinePrevJlo, &prevJlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	err = json.Unmarshal(jsonlineJlo, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	packageDir := PackageDirFromVim

	doBreak, _ := HandleOutputLines(&Results, jlo, prevJlo, packageDir)
	if doBreak != false {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(doBreak), strconv.FormatBool(false))
	}
}

////TestHandleOutputLines ....
func TestHandleOutputLines_FAIL(t *testing.T) {
	Results := GtpResults{}
	Results.Counts = map[string]int{"run": 25, "pause": 0, "continue": 0, "skip": 0, "pass": 21, "fail": 4, "output": 31}
	jlo, prevJlo := JLObject{}, JLObject{}
	jsonlinePrevJlo := []byte(`{"Time":"2021-05-08T08:06:40.543663129-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"    main_test.go:12: got = \"Hello, World!\", want \"!Hello, World!\"\n"}`)
	jsonlineJlo := []byte(`{"Time":"2021-05-08T08:06:40.543669982-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"--- FAIL: TestHello (0.00s)\n"}`)
	err := json.Unmarshal(jsonlinePrevJlo, &prevJlo)
	chkErr(err, "Error Unmarshaling jsonline_prevJlo")
	err = json.Unmarshal(jsonlineJlo, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine_jlo")
	packageDir := PackageDirFromVim

	doBreak, _ := HandleOutputLines(&Results, jlo, prevJlo, packageDir)
	if doBreak != false {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(doBreak), strconv.FormatBool(false))
	}
}

// {"Time":"2021-05-08T08:06:40.543522241-04:00","Action":"run","Package":"github.com/davidsday/hello","Test":"TestHello"}
// {"Time":"2021-05-08T08:06:40.543653446-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"=== RUN   TestHello\n"}
// {"Time":"2021-05-08T08:06:40.543663129-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"    main_test.go:12: got = \"Hello, World!\", want \"!Hello, World!\"\n"}
// {"Time":"2021-05-08T08:06:40.543669982-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"--- FAIL: TestHello (0.00s)\n"}
// {"Time":"2021-05-08T08:06:40.543673412-04:00","Action":"fail","Package":"github.com/davidsday/hello","Test":"TestHello","Elapsed":0}
// {"Time":"2021-05-08T08:06:40.543678142-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"FAIL\n"}
// {"Time":"2021-05-08T08:06:40.543681244-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"coverage: 0.0% of statements\n"}
// {"Time":"2021-05-08T08:06:41.043682815-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"exit status 1\n"}
// {"Time":"2021-05-08T08:06:41.043783023-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"FAIL\tgithub.com/davidsday/hello\t0.501s\n"}
// {"Time":"2021-05-08T08:06:41.043819391-04:00","Action":"fail","Package":"github.com/davidsday/hello","Elapsed":0.501}
