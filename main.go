package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexPanic        = regexp.MustCompile(`^panic:`)
	regexNoTestsToRun = regexp.MustCompile(`no tests to run`)
	regexNoTestFiles  = regexp.MustCompile(`\[no test files\]`)
	regexBuildFailed  = regexp.MustCompile(`\[build failed\]`)
	regexTestFileRef  = regexp.MustCompile(`_test.go`)
	//"coverage: 76.7% of statements\n"}
	regexTestCoverage = regexp.MustCompile(`^coverage: \d{1,2}\.\d\% of statements`)
	regexNil          = &regexp.Regexp{}
)

func main() {

	// Results has all the data from go test
	// It has the methods it needs to build the BarMessage
	// It lives in results.go
	var Results GtpResults

	// Initialize map of Counts in Results
	Results.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 0, "fail": 0, "output": 0}

	// Barmessage includes QfList. They are populated by the methods
	// in Results.  They don't "do" anything except hold
	// the data Vim will need, and they also marshal themselves
	// into JSON and send it to Vim via stdout for display
	// BarMessage lives in barMessage.go
	var Barmessage BarMessage

	// jlo & JLO -> JSON Line Object
	// go test -json spits these out, one at a time, separated by newlines
	// These objects live in jsonLineObject.go
	var jlo JLObject
	// prevJlo gets populated at the bottom of the for loop in
	// case we need to look back at the previous object (line)
	var prevJlo JLObject

	// PackageDirFromVim is where the current package lives
	// We get it from Vim as os.Args[1]
	var PackageDirFromVim string
	// Gocyclo likes to receive lists of paths to search
	// We don't have any, but to avoid mucking with gocyclo internals
	// we create an empty list and append PackageDirFromVim to it splitOnSemiColons
	// gocyclo can be happy
	var PackageDirsToSearch []string

	// We get quidance from Vim about where go test and gocyclo
	// should search, there is really only one dir from Vim,
	// but gocyclo wants a list of dirs, so we create an empty
	// list and append the dir we got from Vim to it so
	// gocyclo will be happy
	PackageDirFromVim = os.Args[1]
	PackageDirsToSearch = append(PackageDirsToSearch, PackageDirFromVim)
	// Vim tells us how many columns it has available for messages via the
	// third command line argument
	Results.VimColumns, _ = strconv.Atoi(os.Args[2])
	// initialize Barmessage's QuickFixList to an empty QuickFixList
	// or it will be null instead of [] when marshaled to JSON.
	Barmessage.QuickFixList = GtpQfList{}

	commandLine := "go test -v -json -cover " + PackageDirFromVim

	stdout, stderr, _ := Shellout(commandLine)
	if rcvdMsgOnStdErr(stderr) {
		doStdErrMsg(stderr, &Results, PackageDirFromVim, &Barmessage)
	} else {
		jsonLines := splitIntoLines(stdout)

		for _, jsonLine := range jsonLines {
			// Ensure we're getting valid JSON
			if !json.Valid(convertStringToBytes(jsonLine)) {
				buildAndAppendAnErrorForInvalidJSON(&Results)
				break
			}

			jlo.unmarshal(jsonLine)

			PackageDirFromJlo := jlo.getPackage()
			Results.incCount(jlo.getAction())

			var err error
			var doBreak bool

			if jlo.getAction() == "output" {
				doBreak, err = HandleOutputLines(&Results, jlo, prevJlo, PackageDirFromJlo, &Barmessage)
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
		Results.Summary.setElapsed(GtpElapsed(jlo.getElapsed()))

		// We've completed the for loop,
		// The last emitted line (JSON Line Object) announces
		// if the run as a whole was a pass or fail.  It does
		// not represent a test.  So it throws off our counts
		// by one.
		Results.Counts["pass"], Results.Counts["fail"] =
			adjustOutSuperfluousFinalResult(jlo.getAction(), &Results)
		// Now we check for PD.Errors and create a
		// yellow bar and  message if appropriate
		Results.buildBarMessage(&Barmessage, PackageDirsToSearch)
	} //Endif

	// Turn our Results object into JSON and send it to stdout
	Barmessage.marshalToStdOut()
	// BarMessage.writeStdErrMsgToDisk()

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

// HandleOutputLines does much of our work, very similarly
// to how we would search and grep go test -v output.
// To do a good job reporting to the user, we still have
// to grep through normal go test -v type outputs.
// go test -json emits these in jlo.Output fields we handle
// this task here
func HandleOutputLines(Results *GtpResults, jlo JLObject, prevJlo JLObject,
	PackageDirFromVim string, Barmessage *BarMessage) (bool, error) {

	var err error = nil
	doBreak := false

	Results.incCount("output")

	doBreak = checkErrorCandidates(Results, jlo.getOutput(), PackageDirFromVim)
	if doBreak {
		return doBreak, err
	}

	if hasTestCoverage(jlo.getOutput()) {
		Results.Summary.setCoverage(jlo.getOutput())
	}

	if hasTestFileReferences(jlo.getOutput()) {
		list := splitOnSemiColons(jlo.getOutput())
		list = removeUnneededFAILPrefix(list)
		if thisIsTheFirstFailure(Results) {
			takeNoteOfFirstFailure(Results, list, prevJlo.getTest())
		}
		qfItem := buildQuickFixItem(os.Args, list, jlo)
		Barmessage.QuickFixList.Add(qfItem)
	}
	return doBreak, err
} // End HandleOutputLines()

// CheckRegx check for a match described by compiled regx with candidate.
// Returns true if theres a match, false otherwise
func CheckRegx(regx *regexp.Regexp, candidate string) bool {
	match := regx.FindString(candidate)
	return match != ""
}

func rcvdMsgOnStdErr(stderror string) bool {
	return len(stderror) > 0
}

func doStdErrMsg(stderr string, Results *GtpResults, PackageDir string, Barmessage *BarMessage) {
	oneSpace := " "
	msg := stderr
	stdErrMsgPrefix := "STDERR:"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	Barmessage.Color = "yellow"
	if stdErrMsgTooLongForOneLine(stderr, stdErrMsgPrefix, stdErrMsgTrailer, Results.VimColumns) {
		writeStdErrMsgToDisk(stderr, PackageDir)
		Barmessage.Message = buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgTrailer, msg, Results.VimColumns)
	} else {
		Barmessage.Message = stdErrMsgPrefix + oneSpace + strings.ReplaceAll(msg, "\n", "|")
		Barmessage.Message = strings.TrimSuffix(Barmessage.Message, "|") + stdErrMsgTrailer
	}
	gtperror := GtpError{Name: "StdErrError", Regex: regexNil, Message: Barmessage.Message, Color: "yellow"}
	Results.Errors = append(Results.Errors, gtperror)
}

func buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgTrailer, msg string, cols int) string {
	oneSpace := " "
	commaSpace := ", "
	retMsg := stdErrMsgPrefix + oneSpace + strings.ReplaceAll(msg, "\n", "|")
	retMsg = strings.TrimSuffix(retMsg, "|")
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

func buildAndAppendAnErrorForInvalidJSON(Results *GtpResults) {
	Results.Errors = append(Results.Errors,
		GtpError{
			Name:    "InvalidJSON",
			Regex:   regexNil,
			Message: "[Invalid JSON]",
			Color:   "yellow",
		})
}
func splitIntoLines(s string) []string {
	lines := strings.Split(s, "\n")
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

func thisIsTheFirstFailure(Results *GtpResults) bool {
	return Results.Counts["fail"] == 0
}

func takeNoteOfFirstFailure(Results *GtpResults, parts []string, testName string) {
	Results.FirstFail.Fname = parts[0]
	Results.FirstFail.Lineno = parts[1]
	Results.FirstFail.Tname = testName
}

func removeUnneededFAILPrefix(list []string) []string {
	if strings.Contains(list[0], "FAIL") {
		// so we take the sublist and continue our work
		list = list[1:]
	}
	return list
}

func splitOnSemiColons(output string) []string {
	return strings.Split(strings.TrimSpace(output), ":")
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

func adjustOutSuperfluousFinalPass(action string, passCount int) int {
	if finalActionWasPass(action) {
		if weHaveHadMoreThanOnePass(passCount) {
			passCount--
		}
	}
	return passCount
}

func adjustOutSuperfluousFinalFail(action string, failCount int) int {
	if finalActionWasFail(action) {
		if weHaveHadMoreThanOneFail(failCount) {
			failCount--
		}
	}
	return failCount
}

func adjustOutSuperfluousFinalResult(action string, Results *GtpResults) (int, int) {
	passCount := adjustOutSuperfluousFinalPass(action, Results.getCount("pass"))
	failCount := adjustOutSuperfluousFinalFail(action, Results.getCount("fail"))
	return passCount, failCount
}

func checkErrorCandidates(Results *GtpResults, output string, PackageDirFromVim string) bool {
	var ErrorCandidates = GtpErrors{
		{Name: "NoTestFiles", Regex: regexNoTestFiles, Message: "In package: " + PackageDirFromVim + ", [No Test Files]", Color: "yellow"},
		{Name: "NoTestsToRun", Regex: regexNoTestsToRun, Message: "In package: " + PackageDirFromVim + ", [Test Files, but No Tests to Run]", Color: "yellow"},
		{Name: "BuildFailed", Regex: regexBuildFailed, Message: "In package: " + PackageDirFromVim + ", [Build Failed]", Color: "yellow"},
		{Name: "Panic", Regex: regexPanic, Message: "In package: " + PackageDirFromVim + ", [Received a Panic]", Color: "yellow"},
	}
	for _, rx := range ErrorCandidates {
		if CheckRegx(rx.Regex, output) {
			Results.Errors.Add(rx)
			// Results.Errors = append(Results.Errors, rx)
			return true
		}
	}
	return false
}

func hasTestCoverage(output string) bool {
	return CheckRegx(regexTestCoverage, output)
}

func hasTestFileReferences(output string) bool {
	// one of the surest fail indicators is an output about a "_test.go" file
	return CheckRegx(regexTestFileRef, output)
}
