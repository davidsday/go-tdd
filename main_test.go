package main

import (
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

//TestDoStdErrMsg ....
func TestDoStdErrMsg(t *testing.T) {
	pgmdata := PgmData{}
	pgmdata.Barmessage.Columns = 135
	msg := "STDERR: This is my message from STDERR."
	PackageDir := "/home/dave/sw/go/goTestParser"
	doStdErrMsg(msg, &pgmdata, PackageDir)
	if pgmdata.Perror.MsgStderr != true {
		t.Errorf("pgmdata.Perror.MsgStderr = %s\n", strconv.FormatBool(pgmdata.Perror.MsgStderr))
	}
}

//TestDoStdErrMsgTooLong ....
func TestDoStdErrMsgTooLong(t *testing.T) {
	pgmdata := PgmData{}
	pgmdata.Barmessage.Columns = 135
	msg := "STDERR: This is my message from STDERR. xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	PackageDir := "/home/dave/sw/go/goTestParser"
	doStdErrMsg(msg, &pgmdata, PackageDir)
	if pgmdata.Perror.MsgStderr != true {
		t.Errorf("pgmdata.Perror.MsgStderr = %s\n", strconv.FormatBool(pgmdata.Perror.MsgStderr))
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
	pgmdata := PgmData{}
	pgmdata.Barmessage.Columns = 135
	buildAndAppendAnErrorForInvalidJSON(&pgmdata)
	if len(pgmdata.Perrors) <= 0 {
		t.Errorf("pgmdata.Perrors has %d elements", len(pgmdata.Perrors))
	}
}

//===========================================================================
// func initializePgmData(pd *PgmData, commandLine string) {
//===========================================================================

//TestInitializePgmData ....
func TestInitializePgmData(t *testing.T) {
	pd := PgmData{}
	commandLine := "go test -v -json -cover " + PackageDirFromVim
	pd.initializePgmData(commandLine)
	host, _ := os.Hostname()
	if pd.Info.Host != host {
		t.Errorf("got '%s' as hostname,  want '%s'", pd.Info.Host, host)
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

//===========================================================================
// func marshallTR(pgmdata PgmData)
//===========================================================================

//TestMarshallTR() ....
//func TestMarshallTR() (t *testing.T) {
// json := {"info":{"host":"dev","user":"dave","begintime":"2021-04-25T07:40:25.484991373-04:00","endtime":"2021-04-25T07:40:26.214428046-04:00","gtp_issued_cmd":"go test -v -json -cover /home/dave/sw/go/hello","gtp_rcvd_args":["/home/dave/.config/nvim/plugged/goTestParser/bin/goTestParser","/home/dave/sw/go/hello"],"test_coverage":"0.0%"},"counts":{"runs":1,"pauses":0,"continues":0,"skips":0,"passes":1,"fails":0,"outputs":8},"firstfailedtest":{"fname":"","tname":"","lineno":""},"elapsed":0.0020000000949949026,"error":{"validjson":true,"notestfiles":false,"panic":false,"buildfailed":false,"msg_stderr":false},"qflist":null,"barmessage":{"color":"green","message":"1 Run, 1 Passed, Test Coverage: 0.0%, in 0.002s"}}

//	pd := PgmData{}
//	pd.Barmessage.Columns = 135

//	marshallTR(pd)
//	want := "What we want"
//	if got != want {
//		t.Errorf("got '%s' want '%s'", got, want)
//	}
//}

//// function to perform marshalling
//func marshallTR(pgmdata PgmData) {
//	// data, err := json.MarshalIndent(pgmdata, "", "    ")
//	data, _ := json.Marshal(pgmdata)
//	_, err := os.Stdout.Write(data)
//	chkErr(err, "Error writing to Stdout in marshallTR()")
//	// err = os.WriteFile("./goTestParserLog.json", data, 0664)
//	//	chkErr(err, "Error writing to ./goTestParserLog.json, in marshallTR()")
//} // end_marshallTR
