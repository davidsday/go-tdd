package main

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"testing"
)

//===========================================================================

//===========================================================================
//===========================================================================
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
// func splitIntoLines(b []byte) [][]byte {
//===========================================================================

//TestSplitIntoLines ....
func TestSplitIntoLines(t *testing.T) {
	s := "Strings\nIntegers\nFloats\nBooleans\n"
	got := splitIntoLines(s)
	want := []string{"Strings", "Integers", "Floats", "Booleans"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestSplitIntoLines_with_non_empty_last_line ....
func TestSplitIntoLines_with_non_empty_last_line(t *testing.T) {
	s := "Strings\nIntegers\nFloats\nBooleans\nBytes"
	got := splitIntoLines(s)
	want := []string{"Strings", "Integers", "Floats", "Booleans", "Bytes"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func thisIsTheFirstFailure(pgmdata PgmData) bool {
//===========================================================================

//TestThisIsTheFirstFailure_true ....
func TestThisIsTheFirstFailure_true(t *testing.T) {
	results := newResults()
	got := thisIsTheFirstFailure(&results)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestThisIsTheFirstFailure_false ....
func TestThisIsTheFirstFailure_false(t *testing.T) {
	results := newResults()

	results.Counts["fail"] = 1
	got := thisIsTheFirstFailure(&results)
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
	results := newResults()
	testName := "thisTest"
	filename := "filename"
	linenum := "12"

	// func takeNoteOfFirstFailure(filename, linenum, testName string, results GtpResults) {
	takeNoteOfFirstFailure(filename, linenum, testName, &results)

	if results.FirstFail.Tname != testName && results.FirstFail.Lineno != linenum && results.FirstFail.Fname != filename {
		t.Errorf("Filename: %s, LineNo: %s, TestName: %s", results.FirstFail.Fname, results.FirstFail.Lineno, results.FirstFail.Tname)
	}
}

//===========================================================================
//TestProcessStdErrMsg
//===========================================================================

//TestProcessStdErrMsg
func TestProcessStdErrMsg(t *testing.T) {
	results := newResults()
	results.Args.setScreenColumns(135)
	Barmessage := newBarMessage()
	want := BarMessage{Color: "yellow", Message: "STDERR: This is my message from STDERR.", QuickFixList: GtpQfList{}}
	msg := "This is my message from STDERR."
	packageDir := "."
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, packageDir)
	processStdErr(msg, &results, PackageDirsToSearch, &Barmessage)
	if !reflect.DeepEqual(want, Barmessage) {
		t.Errorf("Barmessage: '%#v', Want: '%#v'", Barmessage, want)
	}
}

//TestDoStdErrMsgTooLong ....
func TestProcessStdErrMsgTooLong(t *testing.T) {
	results := newResults()
	results.Args.setScreenColumns(72)
	Barmessage := newBarMessage()
	want := BarMessage{Color: "yellow", Message: "STDERR: This is my message from STDERR. xx, [See pkgdir/StdErr.txt]", QuickFixList: GtpQfList{}}
	msg := "This is my message from STDERR. xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	packageDir := "."
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, packageDir)
	processStdErr(msg, &results, PackageDirsToSearch, &Barmessage)
	if !reflect.DeepEqual(want, Barmessage) {
		t.Errorf("Barmessage: '%#v', Want: '%#v'", Barmessage, want)
	}
	_ = os.Remove(PackageDirsToSearch[0] + "/StdErr.txt")
}

//===========================================================================
//TestProcessStdOutMsg
//===========================================================================

//TestProcessStdOutMsg
func TestProcessStdOutMsg1(t *testing.T) {
	results := newResults()
	results.Args.setScreenColumns(135)
	Barmessage := newBarMessage()
	want := BarMessage{Color: "green", Message: "1 Run, 1 Passed, Test Coverage: 0.0%, Average Complexity: NaN, in 0.001s", QuickFixList: GtpQfList{}}
	out := `{"Time":"2021-05-10T09:00:49.114179156-04:00","Action":"run","Package":"github.com/davidsday/hello","Test":"TestHello"}
{"Time":"2021-05-10T09:00:49.114321584-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"=== RUN   TestHello\n"}
{"Time":"2021-05-10T09:00:49.114344537-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"--- PASS: TestHello (0.00s)\n"}
{"Time":"2021-05-10T09:00:49.114350428-04:00","Action":"pass","Package":"github.com/davidsday/hello","Test":"TestHello","Elapsed":0}
{"Time":"2021-05-10T09:00:49.11435592-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"PASS\n"}
{"Time":"2021-05-10T09:00:49.114360603-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"coverage: 0.0% of statements\n"}
{"Time":"2021-05-10T09:00:49.114412878-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"ok  \tgithub.com/davidsday/hello\t0.001s\n"}
{"Time":"2021-05-10T09:00:49.114430234-04:00","Action":"pass","Package":"github.com/davidsday/hello","Elapsed":0.001}`
	packageDir := "./testdata/hello"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, packageDir)
	processStdOut(out, &results, PackageDirsToSearch, &Barmessage)
	if !reflect.DeepEqual(Barmessage, want) {
		t.Errorf("'%v'|'%v'", Barmessage, want)
	}
	// _ = os.Remove(PackageDirsToSearch[0] + "/StdErr.txt")
}

//TestProcessStdOutMsg2
func TestProcessStdOutMsg2(t *testing.T) {
	want := []byte(`{"color":"green","message":"1 Run, 1 Passed, Test Coverage: 50.0%, Average Complexity: NaN, in 0.001s","quickfixlist":[]}`)

	input := `{"Time":"2021-05-10T20:55:51.269642384-04:00","Action":"run","Package":"github.com/davidsday/hello","Test":"Example"}
{"Time":"2021-05-10T20:55:51.269779248-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"Example","Output":"=== RUN   Example\n"}
{"Time":"2021-05-10T20:55:51.269793185-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"Example","Output":"--- PASS: Example (0.00s)\n"}
{"Time":"2021-05-10T20:55:51.269798857-04:00","Action":"pass","Package":"github.com/davidsday/hello","Test":"Example","Elapsed":0}
{"Time":"2021-05-10T20:55:51.269808731-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"PASS\n"}
{"Time":"2021-05-10T20:55:51.269815385-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"coverage: 50.0% of statements\n"}
{"Time":"2021-05-10T20:55:51.269858272-04:00","Action":"output","Package":"github.com/davidsday/hello","Output":"ok  \tgithub.com/davidsday/hello\t0.001s\n"}
{"Time":"2021-05-10T20:55:51.269867388-04:00","Action":"pass","Package":"github.com/davidsday/hello","Elapsed":0.001}`

	results := newResults()
	results.Args.setScreenColumns(135)
	Barmessage := newBarMessage()
	packageDir := "./testdata/hello"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, packageDir)

	processStdOut(input, &results, PackageDirsToSearch, &Barmessage)
	if !reflect.DeepEqual(Barmessage.marshalToByteString(), want) {
		t.Errorf("'%v'|'%v'", Barmessage, want)
	}
}

//TestProcessStdOutMsg3
func TestProcessStdOutMsg3(t *testing.T) {
	want := []byte(`{"color":"red","message":"1 Run, 0 Passed, 1 Failed, 1st in hello_test.go, on line 10, in 0.001s","quickfixlist":[{"filename":"./testdata/hello/hello_test.go","lnum":10,"col":1,"vcol":1,"pattern":"TestHello","text":" got 'hello world' want 'hello World'"}]}`)

	input := `{"Time":"2021-05-23T11:06:16.636557509-04:00","Action":"run","Package":"hello","Test":"TestHello"}
{"Time":"2021-05-23T11:06:16.636882472-04:00","Action":"output","Package":"hello","Test":"TestHello","Output":"=== RUN   TestHello\n"}
{"Time":"2021-05-23T11:06:16.636892791-04:00","Action":"output","Package":"hello","Test":"TestHello","Output":"    hello_test.go:10: got 'hello world' want 'hello World'\n"}
{"Time":"2021-05-23T11:06:16.63690239-04:00","Action":"output","Package":"hello","Test":"TestHello","Output":"--- FAIL: TestHello (0.00s)\n"}
{"Time":"2021-05-23T11:06:16.636906919-04:00","Action":"fail","Package":"hello","Test":"TestHello","Elapsed":0}
{"Time":"2021-05-23T11:06:16.636912727-04:00","Action":"output","Package":"hello","Output":"FAIL\n"}
{"Time":"2021-05-23T11:06:16.636916975-04:00","Action":"output","Package":"hello","Output":"coverage: 50.0% of statements\n"}
{"Time":"2021-05-23T11:06:16.637072805-04:00","Action":"output","Package":"hello","Output":"exit status 1\n"}
{"Time":"2021-05-23T11:06:16.637085928-04:00","Action":"output","Package":"hello","Output":"FAIL\thello\t0.001s\n"}
{"Time":"2021-05-23T11:06:16.6371028-04:00","Action":"fail","Package":"hello","Elapsed":0.001}`

	results := newResults()
	results.Args.setScreenColumns(135)
	Barmessage := newBarMessage()
	packageDir := "./testdata/hello"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, packageDir)

	processStdOut(input, &results, PackageDirsToSearch, &Barmessage)
	if !reflect.DeepEqual(Barmessage.marshalToByteString(), want) {
		t.Errorf("Got: '%v'|Want: '%v'", string(Barmessage.marshalToByteString()), string(want))
	}
}

//TestProcessStdOutMsg4
func TestProcessStdOutMsg4(t *testing.T) {

	want := []byte(`{"color":"yellow","message":"In package: ./testdata/hello, [Test Files, but No Tests to Run]","quickfixlist":[]}`)

	input := `{"Time":"2021-05-11T22:20:14.727345713-04:00","Action":"output","Package":"hello","Output":"testing: warning: no tests to run\n"}
{"Time":"2021-05-11T22:20:14.727527656-04:00","Action":"output","Package":"hello","Output":"PASS\n"}
{"Time":"2021-05-11T22:20:14.727601779-04:00","Action":"output","Package":"hello","Output":"ok  \thello\t0.001s\n"}
{"Time":"2021-05-11T22:20:14.727621504-04:00","Action":"pass","Package":"hello","Elapsed":0.001}`

	results := newResults()
	results.Args.setScreenColumns(135)
	Barmessage := newBarMessage()
	packageDir := "./testdata/hello"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, packageDir)

	processStdOut(input, &results, PackageDirsToSearch, &Barmessage)
	if !reflect.DeepEqual(Barmessage.marshalToByteString(), want) {
		t.Errorf("Got: '%v'|Want: '%v'", string(Barmessage.marshalToByteString()), string(want))
	}
}

////TestProcessStdOutMsg5
func TestProcessStdOutMsg5(t *testing.T) {

	want := []byte(`{"color":"yellow","message":"In package: ./testdata/hello, [No Test Files]","quickfixlist":[]}`)

	input := `{"Time":"2021-05-11T22:30:46.844883222-04:00","Action":"output","Package":"./testdata/hello","Output":"?   \thello\t[no test files]\n"}
{"Time":"2021-05-11T22:30:46.84499945-04:00","Action":"skip","Package":"hello","Elapsed":0}`
	results.Args.setScreenColumns(135)
	Barmessage := newBarMessage()
	packageDir := "./testdata/hello"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, packageDir)

	processStdOut(input, &results, PackageDirsToSearch, &Barmessage)
	if !reflect.DeepEqual(Barmessage.marshalToByteString(), want) {
		t.Errorf("'%v'|'%v'", string(Barmessage.marshalToByteString()), string(want))
	}
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
	results := newResults()
	results.Args.setScreenColumns(135)
	buildAndAppendAnErrorForInvalidJSON(&results)
	if len(results.Errors) <= 0 {
		t.Errorf("pgmdata.Perrors has %d elements", len(results.Errors))
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

//TestSplitOnColons(str string) []string ....
func TestSplitOnColons(t *testing.T) {
	str := "FAIL:1st useful part:2nd useful part:3rd useful part"
	got := splitOnColons(str)
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
// func checkErrorCandidates(results *GtpResults, output string) bool {
//===========================================================================

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_no_test_files(t *testing.T) {
	results := newResults()
	packageDir := "."
	// output := "FAIL:Part1:Part2:Part3"
	output := "[no test files]"
	got := checkErrorCandidates(&results, output, packageDir)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_yes(t *testing.T) {
	packageDir := "."
	results := newResults()
	// output := "FAIL:Part1:Part2:Part3"
	output := "panic: "
	got := checkErrorCandidates(&results, output, packageDir)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_panic_recovered(t *testing.T) {
	packageDir := "."
	results := newResults()
	// output := "FAIL:Part1:Part2:Part3"
	output := "panic: blah blah blah [recovered]"
	got := checkErrorCandidates(&results, output, packageDir)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_fatal_error(t *testing.T) {
	packageDir := "."
	results := newResults()
	output := "fatal error: "
	got := checkErrorCandidates(&results, output, packageDir)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_fatal_error_recovered(t *testing.T) {
	packageDir := "."
	results := newResults()
	output := "fatal error: blah blah blah [recovered]"
	got := checkErrorCandidates(&results, output, packageDir)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestCheckErrorCandidates ....
func TestCheckErrorCandidates_no(t *testing.T) {
	packageDir := "."
	results := newResults()
	output := "Part0:Part1:Part2:Part3"
	got := checkErrorCandidates(&results, output, packageDir)
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
// func adjustOutSuperfluousFinalResult(action string, results *GtpResults) (int, int) {
//===========================================================================

//TestAdjustOutSuperfuousFinalResult ....
func TestAdjustOutSuperfuousFinalResult(t *testing.T) {
	action := "pass"
	results.Counts["pass"] = 21

	passes, fails := adjustOutSuperfluousFinalResult(action, &results)

	if passes != 20 || fails != 0 {
		t.Errorf("passes:'%d', fails: '%d', wanted: passes: '20', fails: 0", passes, fails)
	}
}

//===========================================================================
// func HandleOutputLines(results *GtpResults, jlo JLObject, prevJlo JLObject,
//	packageDir string) (bool, error) {
//===========================================================================

//TestHandleOutputLines ....
func TestHandleOutputLines(t *testing.T) {
	Barmessage := newBarMessage()
	jlo, prevJlo := JLObject{}, JLObject{}
	jsonlinePrevJlo := []byte(`{"Time": "2021-05-07T23:32:18.412171038-04:00", "Action": "output", "Package": "github.com/davidsday/go-tdd", "Output": "PASS\n"}`)
	jsonlineJlo := []byte(`{"Time": "2021-05-07T23:32:18.412174016-04:00", "Action": "output", "Package": "github.com/davidsday/go-tdd", "Output": "coverage: 77.9% of statements\n"}`)
	err := json.Unmarshal(jsonlinePrevJlo, &prevJlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	err = json.Unmarshal(jsonlineJlo, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	jloSlice := []JLObject{prevJlo, jlo}

	doBreak, _ := HandleOutputLines(&results, jloSlice, 1, results.Args.PackageDir, results.Args.PluginDir, &Barmessage)
	if doBreak != false {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(doBreak), strconv.FormatBool(false))
	}
}

////TestHandleOutputLines ....
func TestHandleOutputLines_FAIL(t *testing.T) {
	Barmessage := newBarMessage()
	jlo, prevJlo := JLObject{}, JLObject{}
	jsonlinePrevJlo := []byte(`{"Time":"2021-05-08T08:06:40.543663129-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"    main_test.go:12: got = \"Hello, World!\", want \"!Hello, World!\"\n"}`)
	jsonlineJlo := []byte(`{"Time":"2021-05-08T08:06:40.543669982-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"--- FAIL: TestHello (0.00s)\n"}`)
	err := json.Unmarshal(jsonlinePrevJlo, &prevJlo)
	chkErr(err, "Error Unmarshaling jsonline_prevJlo")
	err = json.Unmarshal(jsonlineJlo, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine_jlo")
	jloSlice := []JLObject{prevJlo, jlo}

	doBreak, _ := HandleOutputLines(&results, jloSlice, 1, results.Args.PackageDir, results.Args.PluginDir, &Barmessage)
	if doBreak != false {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(doBreak), strconv.FormatBool(false))
	}
}

////TestHandleOutputLines ....
func TestHandleOutputLines_TestFileRef(t *testing.T) {
	Barmessage := newBarMessage()
	jlo, prevJlo := JLObject{}, JLObject{}
	jlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543663129-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"    main_test.go:12: got = \"Hello, World!\", want \"!Hello, World!\"\n"}`)
	prevJlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543669982-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"--- FAIL: TestHello (0.00s)\n"}`)
	jloSlice := []JLObject{prevJlo, jlo}

	doBreak, _ := HandleOutputLines(&results, jloSlice, 1, results.Args.PackageDir, results.Args.PluginDir, &Barmessage)
	if doBreak != false {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(doBreak), strconv.FormatBool(false))
	}
}

////TestHandleOutputLines ....
func TestHandleOutputLines_received_a_panic(t *testing.T) {
	Barmessage := newBarMessage()
	jlo, prevJlo := JLObject{}, JLObject{}
	prevJlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543663129-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"    main_test.go:12: got = \"Hello, World!\", want \"!Hello, World!\"\n"}`)
	jlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543669982-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"panic: "}`)
	jloSlice := []JLObject{prevJlo, jlo}

	doBreak, _ := HandleOutputLines(&results, jloSlice, 1, results.Args.PackageDir, results.Args.PluginDir, &Barmessage)
	if doBreak != true {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(doBreak), strconv.FormatBool(true))
	}
}

// func safeToLookAhead(jloSlice []JLObject, i, incr int) bool {

//TestSafeToLookAhead_true
func TestSafeToLookAhead_true(t *testing.T) {
	jlo, prevJlo := JLObject{}, JLObject{}
	jlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543663129-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"    main_test.go:12: got = \"Hello, World!\", want \"!Hello, World!\"\n"}`)
	prevJlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543669982-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"--- FAIL: TestHello (0.00s)\n"}`)
	jloSlice := []JLObject{prevJlo, jlo}
	got := safeToLookAhead(jloSlice, 0, 1)
	want := true
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}

//TestSafeToLookAhead_false
func TestSafeToLookAhead_false(t *testing.T) {
	jlo, prevJlo := JLObject{}, JLObject{}
	jlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543663129-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"    main_test.go:12: got = \"Hello, World!\", want \"!Hello, World!\"\n"}`)
	prevJlo.unmarshal(`{"Time":"2021-05-08T08:06:40.543669982-04:00","Action":"output","Package":"github.com/davidsday/hello","Test":"TestHello","Output":"--- FAIL: TestHello (0.00s)\n"}`)
	jloSlice := []JLObject{prevJlo, jlo}
	got := safeToLookAhead(jloSlice, 0, 3)
	want := false
	if got != want {
		t.Errorf("got '%s' want '%s'", strconv.FormatBool(got), strconv.FormatBool(want))
	}
}
