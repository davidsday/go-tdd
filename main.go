package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fzipp/gocyclo"
)

// ?    github.com/zchee/nvim-go/pkg/server [no test files]
var (
	regexPanic = regexp.MustCompile(`^panic:`)
	// regexNoTestsToRun   = regexp.MustCompile(`^testing: warning: no tests to run`)
	regexNoTestsToRun   = regexp.MustCompile(`no tests to run`)
	regexNoTestFiles    = regexp.MustCompile(`\?\s*\S*\s*\[no test files\]`)
	regexBuildFailed    = regexp.MustCompile(`\?\s*\S*\s*\[build failed\]`)
	regexFailorTestFile = regexp.MustCompile(`^\s\+FAIL:|_test.go`)
	regexTestCoverage   = regexp.MustCompile(`^coverage:`)
	regexNil            = &regexp.Regexp{}
)

// JLObject -
// go test -json outputs JSON objects instead of lines
// each JSON object looks like this. Not all fields
// are emitted for each line
type JLObject struct {
	Time    string
	Action  string
	Package string
	Test    string
	Output  string
	Elapsed float32
}

// PD -> program data
// This is the whole enchilada
// It holds all the program's
// important data.  It is defined
// in pgmdata.go
//
var PD PgmData

// jlo & JLO -> JSON Line Object
// go test -json spits these out, one at a time, separated by newlines
var jlo JLObject
var prevJlo JLObject

// PackageDir is where the current package lives
// We get it from Vim as os.Args[1]
var PackageDirFromVim string
var PackagesToSearch []string

func main() {

	PackageDirFromVim = os.Args[1]
	PackagesToSearch = append(PackagesToSearch, PackageDirFromVim)

	commandLine := "go test -v -json -cover " + PackageDirFromVim
	PD.initializePgmData(commandLine)

	stdout, stderr, _ := Shellout(commandLine)
	if rcvdMsgOnStdErr(stderr) {
		doStdErrMsg(stderr, &PD, PackageDirFromVim)
	} else {
		// stdout & stderr are strings, we need []byte
		byteString := convertStringToBytes(stdout)
		byteLines := splitBytesIntoLines(byteString)

		// Now we should have valid JSON lines only
		for _, jsonLine := range byteLines {
			// Ensure we're getting valid JSON
			if !json.Valid(jsonLine) {
				buildAndAppendAnErrorForInvalidJSON(&PD)
				break
			} else {
				// Convert line of JSON text to JSON line object (Go struct in this case)
				err := json.Unmarshal(jsonLine, &jlo)
				chkErr(err, "Error Unmarshaling jsonLine")
			}

			PackageDirFromJlo := jlo.Package
			PD.Counts[jlo.Action]++

			var err error
			var doBreak bool

			if jlo.Action == "output" {
				doBreak, err = HandleOutputLines(&PD, jlo, prevJlo, PackageDirFromJlo)
				chkErr(err, "Error in HandleOutputLines()")
				if doBreak {
					break
				}
			}
			// Bottom of for loop - current JSON Line Object now
			// becomes the Previous JSON Line Object,
			// for look back purposes ...
			prevJlo = jlo
		} //endfor

		// Make note of the elapsed time, as reported by go test
		PD.Elapsed = PDElapsed(jlo.Elapsed)

		// We've completed the for loop,
		// The last emitted line (JSON Line Object) announces
		// if the run as a whole was a pass or fail.  It does
		// not represent a test.  So it throws off our counts
		// by one.
		adjustOutSuperfluousFinalResult()
		// Now we check for PD.Errors and create a
		// yellow bar and  message if appropriate
		PD.setBarMessage()
	} //Endif

	// Endtime for PD.Info
	PD.Info.Endtime = time.Now().Format(time.RFC3339Nano)
	// Turn our PD object into JSON and send it to stdout
	marshallTR(PD)

} // endmain()

// Shellout - run a command, capturing stdout, stderr, and errors
func Shellout(command string) (string, string, error) {
	// Force POSIX compliant shell for predictability
	var ShellToUse string = "/bin/sh"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
} //end_Shellout()

// function to perform marshaling
func marshallTR(pgmdata PgmData) {
	// data, err := json.MarshalIndent(pgmdata, "", "    ")
	data, _ := json.Marshal(pgmdata)
	_, err := os.Stdout.Write(data)
	chkErr(err, "Error writing to Stdout in marshallTR()")
	// err = os.WriteFile("./goTestParserLog.json", data, 0664)
	//	chkErr(err, "Error writing to ./goTestParserLog.json, in marshallTR()")
} // end_marshallTR

// HandlOutputLines() does much of our work, very similarly
// to how we would search and grep go test -v output.
// to do a good job reporting to the user, we still have
// to grep through normal go test -v type outputs.
// go test -json emits these in jlo.Output fields we handle
// this task here
func HandleOutputLines(pgmdata *PgmData, jlo JLObject, prevJlo JLObject,
	PackageDirFromVim string) (bool, error) {

	var err error = nil
	doBreak := false

	pgmdata.Counts["output"]++

	doBreak = checkErrorCandidates(&PD)
	if doBreak {
		return doBreak, err
	}

	checkTestCoverage(&PD)

	checkForFAILs(&PD, jlo, prevJlo)

	return doBreak, err
} // End HandleOutputLines()

func passMsg(passes int) string {
	oneSpace := " "
	commaSpace := ", "
	return commaSpace + strconv.Itoa(passes) + oneSpace + "Passed"
}

func runMsg(runs int) string {
	oneSpace := " "
	return strconv.Itoa(runs) + oneSpace + "Run"
}

func elapsedMsg(elapsed PDElapsed) string {
	oneSpace := " "
	commaSpace := ", "
	msg := commaSpace + "in" + oneSpace + strconv.FormatFloat(float64(elapsed), 'f', 3, 32) + "s"
	return msg
}

func metricsMsg(skips, fails int, coverage, complexity string) string {
	oneSpace := " "
	commaSpace := ", "
	if skips == 0 && fails == 0 && len(coverage) > 0 {
		msg := commaSpace + "Test Coverage:" + oneSpace + coverage
		msg += commaSpace + "Average Complexity:" + oneSpace + complexity
		return msg
	}
	return ""
}

func failMsg(fails int, fname, lineno string) string {
	if fails > 0 {
		oneSpace := " "
		commaSpace := ", "
		msg := commaSpace + strconv.Itoa(fails) + oneSpace + "Failed"
		msg += commaSpace + "1st in" + oneSpace + fname
		msg += commaSpace + "on line" + oneSpace + lineno
		return msg
	}
	return ""
}

func skipMsg(skips int) string {
	oneSpace := " "
	commaSpace := ", "
	if skips > 0 {
		return commaSpace + strconv.Itoa(skips) + oneSpace + "Skipped"
	}
	return ""
}

// CheckRegx check for a match described by compiled regx with candidate.
// Returns true if theres a match, false otherwise
func CheckRegx(regx *regexp.Regexp, candidate string) bool {
	match := regx.FindString(candidate)
	return match != ""
}

func getAvgCyclomaticComplexity(paths []string) string {
	// allStats := gocyclo.Analyze(paths, regexp.MustCompile("vendor"))
	allStats := gocyclo.Analyze(paths, regexp.MustCompile(`vendor`))
	return fmt.Sprintf("%.3g", allStats.AverageComplexity())
}

func rcvdMsgOnStdErr(stderror string) bool {
	return len(stderror) > 0
}

func doStdErrMsg(stderr string, pd *PgmData, PackageDir string) {
	oneSpace := " "
	msg := stderr
	stdErrMsgPrefix := "STDERR:"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	pd.Barmessage.Color = "yellow"
	if stdErrMsgTooLongForOneLine(stderr, stdErrMsgPrefix, stdErrMsgTrailer, pd.Barmessage.Columns) {
		writeStdErrMsgToDisk(stderr, PackageDir)
		pd.Barmessage.Message = buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgTrailer, msg, pd.Barmessage.Columns)
	} else {
		pd.Barmessage.Message = stdErrMsgPrefix + oneSpace + strings.ReplaceAll(msg, "\n", "|")
		pd.Barmessage.Message = strings.TrimSuffix(pd.Barmessage.Message, "|") + stdErrMsgTrailer
	}
	gtperror := GTPerror{Name: "StdErrError", Regex: regexNil, Message: pd.Barmessage.Message, Color: "yellow"}
	pd.Perrors = append(PD.Perrors, gtperror)
}

func buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgTrailer, msg string, cols int) string {
	oneSpace := " "
	commaSpace := ", "
	retMsg := stdErrMsgPrefix + oneSpace + strings.ReplaceAll(msg, "\n", "|")
	retMsg = strings.TrimSuffix(msg, "|")
	retMsg = retMsg[0 : cols-(len(stdErrMsgPrefix)+len(stdErrMsgTrailer))]
	retMsg += commaSpace + stdErrMsgTrailer
	return retMsg
}

func stdErrMsgTooLongForOneLine(stderr, stdErrMsgPrefix, stdErrMsgTrailer string, cols int) bool {
	oneSpace := " "
	return (len(stderr) > (cols - (len(stdErrMsgTrailer) + len(stdErrMsgPrefix) + len(oneSpace))))
}

func writeStdErrMsgToDisk(stderr, pkgdir string) {
	path := pkgdir + "/StdErr.txt"
	err := os.WriteFile(path, []byte(stderr), 0664)
	chkErr(err, "error writing "+path)
}

func chkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("Error: %v, %s", err, msg)
	}
}

func buildAndAppendAnErrorForInvalidJSON(pd *PgmData) {
	pd.Perrors = append(pd.Perrors,
		GTPerror{
			Name:    "InvalidJSON",
			Regex:   regexNil,
			Message: "[Invalid JSON]",
			Color:   "yellow",
		})
}
func splitBytesIntoLines(b []byte) [][]byte {
	// stdout & stderr are strings, we need []byte
	lines := bytes.Split(b, []byte("\n"))
	//bytes.Split returns an empty line AFTER the final "\n"
	// so we drop that one
	if len(lines[len(lines)-1]) > 0 {
		return lines
	}
	return lines[:len(lines)-1]
}

func convertStringToBytes(s string) []byte {
	return []byte(s)
}

func thisIsTheFirstFailure(pgmdata *PgmData) bool {
	return pgmdata.Counts["fail"] == 0
}

func takeNoteOfFirstFailure(pgmdata *PgmData, parts []string, testName string) {
	pgmdata.Firstfailedtest.Fname = parts[0]
	pgmdata.Firstfailedtest.Lineno = parts[1]
	pgmdata.Firstfailedtest.Tname = testName
}

func removeUnneededFAILPrefix(output string) []string {
	parts := strings.Split(strings.TrimSpace(output), ":")
	if strings.Contains(parts[0], "FAIL") {
		// so we take the sublist and continue our work
		parts = parts[1:]
	}
	return parts
}

func addToQuickFixList(pgmdata *PgmData, args []string, parts []string, jlo JLObject) {
	var QfItem PDQfDict
	// Now we can build/fill the QuickFix List
	// QfItem.Filename = PackageDir + "/" + parts[0]
	QfItem.Filename = args[1] + "/" + parts[0]
	QfItem.Lnum, _ = strconv.Atoi(parts[1])
	QfItem.Col = 1
	QfItem.Vcol = 1
	QfItem.Pattern = jlo.Test
	QfItem.Text = strings.Join(parts[2:], ":")
	pgmdata.QfList = append(pgmdata.QfList, QfItem)
}

func finalActionWasPass(action string) bool {
	return action == "pass"
}

func finalActionWasFail(action string) bool {
	return action == "fail"
}

func weHaveHadMoreThanOnePass(passes int) bool {
	return passes > 1
}

func weHaveHadMoreThanOneFail(fails int) bool {
	return fails > 1
}

func adjustOutSuperfluousFinalPass() {
	if finalActionWasPass(jlo.Action) {
		if weHaveHadMoreThanOnePass(PD.Counts["pass"]) {
			PD.Counts["pass"]--
		}
	}
}

func adjustOutSuperfluousFinalFail() {
	if finalActionWasFail(jlo.Action) {
		if weHaveHadMoreThanOneFail(PD.Counts["fail"]) {
			PD.Counts["fail"]--
		}
	}
}

func adjustOutSuperfluousFinalResult() {
	adjustOutSuperfluousFinalPass()
	adjustOutSuperfluousFinalFail()
}

func checkErrorCandidates(pd *PgmData) bool {
	var ErrorCandidates = GTPerrors{
		{Name: "NoTestFiles", Regex: regexNoTestFiles, Message: "In package: " + PackageDirFromVim + ", [No Test Files]", Color: "yellow"},
		{Name: "NoTestsToRun", Regex: regexNoTestsToRun, Message: "In package: " + PackageDirFromVim + ", [Test Files, but No Tests to Run]", Color: "yellow"},
		{Name: "BuildFailed", Regex: regexBuildFailed, Message: "In package: " + PackageDirFromVim + ", [Build Failed]", Color: "yellow"},
		{Name: "Panic", Regex: regexPanic, Message: "In package: " + PackageDirFromVim + ", [Received a Panic]", Color: "yellow"},
	}
	for _, rx := range ErrorCandidates {
		if CheckRegx(rx.Regex, jlo.Output) {
			pd.Perrors = append(pd.Perrors, rx)
			return true
		}
	}
	return false
}

func checkTestCoverage(pd *PgmData) {
	if CheckRegx(regexTestCoverage, jlo.Output) {
		// Remove trailing '\n'
		pd.Info.TestCoverage = strings.TrimSuffix(jlo.Output, "\n")
		// Strip away everything but the percent coverage string ("57.8%", for example)
		pd.Info.TestCoverage = strings.Replace(pd.Info.TestCoverage, "coverage: ", "", 1)
		pd.Info.TestCoverage = strings.Replace(pd.Info.TestCoverage, " of statements", "", 1)
	}
}

func checkForFAILs(pd *PgmData, jlo, prevJlo JLObject) {
	// one of the surest fail indicators is an output about a "_test.go" file
	if CheckRegx(regexFailorTestFile, jlo.Output) {
		parts := removeUnneededFAILPrefix(jlo.Output)
		if thisIsTheFirstFailure(pd) {
			takeNoteOfFirstFailure(pd, parts, prevJlo.Test)
		}
		addToQuickFixList(pd, os.Args, parts, jlo)
	}
}
