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
	regexTestCoverage = regexp.MustCompile(`^coverage: \d{1,3}\.\d{0,1}\% of statements`)
	regexNil          = &regexp.Regexp{}
)

func main() {

	// user, _ := user.Current()
	// User := user.Username
	// HomeDir := user.HomeDir

	// results has all the results we collect from go test
	// to help us decide how to present the results to the user
	// It has the methods it needs to build the BarMessage
	// It lives in results.go
	// We have built a func, newResults(), which creates, initializes
	// and returns the results for us
	results := newResults()

	// barMessage includes QfList. They are populated by the methods
	// in Results.  They don't "do" anything except hold
	// the data Vim will need, and they also marshal themselves
	// into JSON and send it to Vim via stdout for display
	// or to disk for logging purposes
	// BarMessage lives in barMessage.go

	barMessage := newBarMessage()

	// Gocyclo likes to receive lists of paths to search
	// We don't have any, but to avoid mucking with gocyclo internals
	// we create an empty list and append os.Args[1] to it so
	// gocyclo can be happy
	var packageDirsToSearch []string

	// We get quidance from Vim about where go test and gocyclo
	// should search, there is really only one dir from Vim,
	// but gocyclo wants a list of dirs, so we create an empty
	// list and append the dir we got from Vim to it so
	// gocyclo will be happy
	packageDirsToSearch = append(packageDirsToSearch, os.Args[1])
	// Vim tells us how many columns it has available for messages via the
	// third command line argument
	results.VimColumns, _ = strconv.Atoi(os.Args[2])

	commandLine := "go test -v -json -cover " + packageDirsToSearch[0]
	stdout, stderr, _ := Shellout(commandLine)

	if rcvdMsgOnStdErr(stderr) {
		processStdErr(stderr, &results, packageDirsToSearch, &barMessage)
	} else {
		processStdOut(stdout, &results, packageDirsToSearch, &barMessage)
	}

	// Turn our Barmessage object into JSON and send it to stdout
	barMessage.marshalToStdOut()
	// and save it to disk
	// barMessage.marshalToDisk()

} // endmain()

func processStdOut(stdout string, results *GtpResults, PackageDirsToSearch []string, Barmessage *BarMessage) {
	// jlo & JLO -> JSON Line Object
	// go test -json spits these out, one at a time, separated by newlines
	// These objects are defined in jsonLineObject.go
	var jlo JLObject
	// prevJlo gets populated at the bottom of the for loop in
	// case we need to look back at the previous object (line)
	// and we do....
	var prevJlo JLObject

	jsonLines := splitIntoLines(stdout)

	for _, jsonLine := range jsonLines {
		// Ensure we're getting valid JSON
		if !json.Valid(convertStringToBytes(jsonLine)) {
			buildAndAppendAnErrorForInvalidJSON(results)
			break
		}
		// jsonLine -> jsonLineObject, which is a Go struct
		// from here down to the bottom of the for loop,
		// we are dealing with JLObject structs
		jlo.unmarshal(jsonLine)

		PackageDirFromJlo := jlo.getPackage()
		results.incCount(jlo.getAction())

		var err error
		var doBreak bool

		if jlo.getAction() == "output" {
			doBreak, err = HandleOutputLines(results, jlo, prevJlo, PackageDirFromJlo, Barmessage)
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
	results.Summary.setElapsed(GtpElapsed(jlo.getElapsed()))

	// We've completed the for loop,
	// The last emitted line (JSON Line Object) announces
	// if the run as a whole was a pass or fail.  It does
	// not represent a test, but get counted as one.
	// So it throws off our counts by one.
	// So we fix that here
	results.Counts["pass"], results.Counts["fail"] =
		adjustOutSuperfluousFinalResult(jlo.getAction(), results)
	// Now we check for results.Errors and create a
	// yellow bar and  message if appropriate
	results.buildBarMessage(Barmessage, PackageDirsToSearch)
}

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
// go test -json emits these in jlo.Output fields. We handle
// this task here
func HandleOutputLines(results *GtpResults, jlo JLObject, prevJlo JLObject,
	packageDir string, Barmessage *BarMessage) (bool, error) {

	var err error = nil
	doBreak := false

	results.incCount("output")

	doBreak = checkErrorCandidates(results, jlo.getOutput(), packageDir)
	if doBreak {
		return doBreak, err
	}

	if hasTestCoverage(jlo.getOutput()) {
		results.Summary.setCoverage(jlo.getOutput())
	}

	// I a jlo.Output field refers to a _test.go file, there has been a
	// test failure and it is telling us in which file and on which line
	// the failure was triggered
	if hasTestFileReferences(jlo.getOutput()) {
		list := splitOnSemiColons(jlo.getOutput())
		// This may be obsolete, we will watch and see...
		list = removeUnneededFAILPrefix(list)
		if thisIsTheFirstFailure(results) {
			takeNoteOfFirstFailure(results, list, prevJlo.getTest())
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

func processStdErr(stderr string, results *GtpResults, PackageDirsToSearch []string, Barmessage *BarMessage) {
	oneSpace := " "
	msg := stderr
	stdErrMsgPrefix := "STDERR:"
	stdErrMsgTrailer := "[See pkgdir/StdErr.txt]"
	Barmessage.Color = "yellow"
	if stdErrMsgTooLongForOneLine(stderr, stdErrMsgPrefix, stdErrMsgTrailer, results.VimColumns) {
		writeStdErrMsgToDisk(stderr, PackageDirsToSearch[0])
		Barmessage.Message = buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgTrailer, msg, results.VimColumns)
	} else {
		Barmessage.Message = stdErrMsgPrefix + oneSpace + strings.ReplaceAll(msg, "\n", "|")
		Barmessage.Message = strings.TrimSuffix(Barmessage.Message, "|") + stdErrMsgTrailer
	}
	gtperror := GtpError{Name: "StdErrError", Regex: regexNil, Message: Barmessage.Message, Color: "yellow"}
	results.Errors = append(results.Errors, gtperror)
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

func buildAndAppendAnErrorForInvalidJSON(results *GtpResults) {
	results.Errors = append(results.Errors,
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

func thisIsTheFirstFailure(results *GtpResults) bool {
	return results.Counts["fail"] == 0
}

func takeNoteOfFirstFailure(results *GtpResults, parts []string, testName string) {
	results.FirstFail.setFname(parts[0])
	results.FirstFail.setLineno(parts[1])
	results.FirstFail.setTname(testName)
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

func adjustOutSuperfluousFinalResult(action string, results *GtpResults) (int, int) {
	passCount := adjustOutSuperfluousFinalPass(action, results.getCount("pass"))
	failCount := adjustOutSuperfluousFinalFail(action, results.getCount("fail"))
	return passCount, failCount
}

func checkErrorCandidates(results *GtpResults, output string, PackageDir string) bool {
	var ErrorCandidates = GtpErrors{
		{Name: "NoTestFiles", Regex: regexNoTestFiles, Message: "In package: " + PackageDir + ", [No Test Files]", Color: "yellow"},
		{Name: "NoTestsToRun", Regex: regexNoTestsToRun, Message: "In package: " + PackageDir + ", [Test Files, but No Tests to Run]", Color: "yellow"},
		{Name: "BuildFailed", Regex: regexBuildFailed, Message: "In package: " + PackageDir + ", [Build Failed]", Color: "yellow"},
		{Name: "Panic", Regex: regexPanic, Message: "In package: " + PackageDir + ", [Received a Panic]", Color: "yellow"},
	}
	for _, rx := range ErrorCandidates {
		if CheckRegx(rx.Regex, output) {
			results.Errors.Add(rx)
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
