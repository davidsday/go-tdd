package main

import "testing"

//============================================================================
//TestGetCounts(key string) int ....
//============================================================================

func TestGetCounts(t *testing.T) {
	results := GtpResults{}
	results.init()
	results.Counts["run"] = 15
	got := results.getCount("run")
	want := 15
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//============================================================================
//TestIncCounts(key string) ....
//============================================================================

func TestIncCount(t *testing.T) {
	results := GtpResults{}
	results.init()
	results.Counts["run"] = 15
	results.incCount("run")
	got := results.getCount("run")

	want := 16
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//============================================================================
//TestDecCount(key string) ....
//============================================================================

func TestDecCount(t *testing.T) {
	results := GtpResults{}
	results.init()
	results.Counts["run"] = 15
	results.decCount("run")
	got := results.getCount("run")

	want := 14
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//============================================================================
//TestgetAverageCyclomaticComplexity ....
//============================================================================

//TestgetAverageCyclomaticComplexity ....
func TestGetAverageCyclomaticComplexity(t *testing.T) {
	paths := []string{}
	results := GtpResults{}
	paths = append(paths, "../gocyclotests/avgCCmplx/main.go")
	results.Summary.setComplexity(paths)
	got := results.Summary.getComplexity()
	want := "7.29"
	if string(got) != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestgetAverageCyclomaticComplexity ....
func TestGetAverageCyclomaticComplexity_no_go_files(t *testing.T) {
	paths := []string{}
	results := GtpResults{}
	paths = append(paths, "./bin/")
	results.Summary.setComplexity(paths)
	got := results.Summary.getComplexity()
	want := "NaN"
	if string(got) != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//============================================================================
//TestbuildBarMessage()
//============================================================================

//TestBuildBarMessage ....
func TestBuildBarMessage_has_errors(t *testing.T) {
	PackageDirFromVim := "/home/dave/sw/go/goTestParser"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, PackageDirFromVim)
	barmsg := BarMessage{}
	results := GtpResults{}
	results.Errors.Add(GtpError{Name: "NoTestFiles", Regex: regexNoTestFiles, Message: "In package: " + PackageDirFromVim + ", [No Test Files]", Color: "yellow"})

	results.buildBarMessage(&barmsg, PackageDirsToSearch)

	got := results.Errors[0].Color
	want := "yellow"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestBuildBarMessage ....
func TestBuildBarMessage_has_fails(t *testing.T) {
	PackageDirFromVim := "/home/dave/sw/go/goTestParser"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, PackageDirFromVim)
	barmsg := BarMessage{}
	results := GtpResults{}
	results.init()
	results.incCount("fail")
	results.incCount("skip")

	results.buildBarMessage(&barmsg, PackageDirsToSearch)

	got := barmsg.getColor()
	want := "red"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestBuildBarMessage_no_fails_but_skips ....
func TestBuildBarMessage_no_fails_but_skips(t *testing.T) {
	PackageDirFromVim := "/home/dave/sw/go/goTestParser"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, PackageDirFromVim)
	barmsg := BarMessage{}
	results := GtpResults{}
	results.init()
	results.incCount("skip")

	results.buildBarMessage(&barmsg, PackageDirsToSearch)

	got := barmsg.getColor()
	want := "yellow"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestBuildBarMessage_no_fails_but_skips ....
func TestBuildBarMessage_all_pass(t *testing.T) {
	PackageDirFromVim := "/home/dave/sw/go/goTestParser"
	PackageDirsToSearch := []string{}
	PackageDirsToSearch = append(PackageDirsToSearch, PackageDirFromVim)
	barmsg := BarMessage{}
	results := GtpResults{}
	results.init()
	results.buildBarMessage(&barmsg, PackageDirsToSearch)

	got := barmsg.getColor()
	want := "green"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
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
// func elapsedMsg(elapsed GtpElapsed) string
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

//============================================================================
//TestSetCoverage()
//============================================================================

// coverage: 58.3% of statements
// func (s *GtpSummary) setCoverage(coverage string) {

//TestSetCoverage()
func TestSetCoverage(t *testing.T) {
	results := GtpResults{}
	coverage := "coverage: 58.3% of statements"
	results.Summary.setCoverage(coverage)
	got := results.Summary.getCoverage()
	want := "58.3%"
	if string(got) != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//============================================================================
// func (e *GtpError) getColor() string {
//============================================================================

//TestGtpError.GetColor ....
func TestGtpError_GetColor(t *testing.T) {
	PackageDirFromVim := "/home/dave/sw/go/goTestParser"
	results := GtpResults{}
	results.Errors.Add(GtpError{Name: "NoTestFiles", Regex: regexNoTestFiles, Message: "In package: " + PackageDirFromVim + ", [No Test Files]", Color: "yellow"})
	got := results.Errors[0].getColor()
	want := "yellow"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//============================================================================
// func (e *GtpError) getMessage() string {
//============================================================================

//TestGtpError.GetColor ....
func TestGtpError_GetMessage(t *testing.T) {
	PackageDirFromVim := "/home/dave/sw/go/goTestParser"
	results := GtpResults{}
	results.Errors.Add(GtpError{Name: "NoTestFiles", Regex: regexNoTestFiles, Message: "In package: " + PackageDirFromVim + ", [No Test Files]", Color: "yellow"})
	got := results.Errors[0].getMessage()
	want := "In package: /home/dave/sw/go/goTestParser, [No Test Files]"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//============================================================================
// func (s *GtpSummary) setElapsed(elapsed GtpElapsed)
//============================================================================

//TestSummarySetElapsed ....
func TestSummarySetElapsed(t *testing.T) {
	results := GtpResults{}
	results.Summary.setElapsed(0.5)
	got := results.Summary.getElapsed()
	want := GtpElapsed(0.5)
	if got != want {
		t.Errorf("got '%f' want '%f'", got, want)
	}
}
