package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexPanic        = regexp.MustCompile(`^(panic: |fatal error: ).*(\[recovered\]){0,1}`)
	regexTimeOutPanic = regexp.MustCompile(`^(panic: |fatal error: )test timed out .*(\[recovered\]){0,1}`)
	regexNoTestsToRun = regexp.MustCompile(`no tests to run`)
	regexNoTestFiles  = regexp.MustCompile(`\[no test files\]`)
	regexBuildFailed  = regexp.MustCompile(`\[build failed\]`)
	regexTestFileRef  = regexp.MustCompile(`_test.go`)
	regexExampleFail  = regexp.MustCompile(`--- FAIL: Example.*\(\d{1,3}\.\d{1,3}s\)`)
	// regexSubTestPass  = regexp.MustCompile(`(    )\+PASS: `)
	// regexSubTestFail  = regexp.MustCompile(`(    )\+FAIL: `)

	//"coverage: 76.7% of statements\n"}
	regexTestCoverage = regexp.MustCompile(`^coverage: \d{1,3}\.\d{0,1}\% of statements`)
	// "coverage: [no statements]\n"
	regexCoverageNoStmts = regexp.MustCompile(`^coverage: \[no statements\]`)
	regexNil             = &regexp.Regexp{}
)
var debug bool
var pluginDir string
var PackageDir string

func main() {

	// if runtime.GOOS = 'windows' {
	// just thinking about portability.....
	//}
	// user, _ := user.Current()
	// User := user.Username
	// HomeDir := user.HomeDir

	// results has all the results we collect from go test
	// to help us decide how to present the results to the user
	// It has the methods it needs to build the BarMessage
	// It lives in results.go
	// We have built a func, newResults(), which creates, initializes
	// and returns a new results to us
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
	// Gocyclo accepts a regex as a request to ignore paths which
	// match the regex.  There is a vim global g:gocyclo_ignore which
	// defaults to 'vendor|testdata' which the user may set to his/here
	// preference and gocyclo will ignore those matching directories
	// vendor is where go projects keep package dependencies and testdata
	// is where we keep our, well, testdata, dirs and files that our tests
	// need
	if len(os.Args) > 3 {
		results.GocycloIgnore = os.Args[3]
	}

	// Turn Debug off
	debug = false
	// The user may also request some debugging logging via
	// this argument
	debug = setDebug(os.Args)
	if debug {
		setupLogging()
	}

	pluginDir = os.Args[5]

	commandLine := "go test -v -json -cover " + packageDirsToSearch[0]
	stdout, stderr, _ := Shellout(commandLine)

	if rcvdMsgOnStdErr(stderr) {
		processStdErr(stderr, &results, packageDirsToSearch, &barMessage)
	} else {
		processStdOut(stdout, &results, packageDirsToSearch, &barMessage)
	}

	// Turn our Barmessage object into JSON and send it to stdout
	barMessage.marshalToStdOut()
	// and/or save it to disk
	if debug {
		barMessage.marshalToDisk()
	}

} // endmain()

func processStdOut(stdout string, results *GtpResults, PackageDirsToSearch []string, Barmessage *BarMessage) {
	// jlo & JLO -> JSON Line Object
	// go test -json spits these out, one at a time, separated by newlines
	// These objects are defined in jsonLineObject.go
	var jlo JLObject
	var jloSlice []JLObject

	// the var stdout is one long line, separated by newlines
	// split them and convert each to a JLObject and append
	// to jloSlice
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
		// in jloSlice
		jlo.unmarshal(jsonLine)
		jloSlice = append(jloSlice, jlo)
	}

	PackageDir = PackageDirsToSearch[0]
	// loop through the slice
	// for ExampleFunctions, we may have to
	// peak ahead up to four lines (jloObjects)
	// so item := range jloSlice wouldn't work
	for i := 0; i < len(jloSlice); i++ {

		results.incCount(jloSlice[i].getAction())

		var err error
		var doBreak bool

		if jloSlice[i].getAction() == "output" {
			doBreak, err = HandleOutputLines(results, jloSlice, i, PackageDir, pluginDir, Barmessage)
			chkErr(err, "Error in HandleOutputLines()")
			if doBreak {
				break
			}
		}
	} //endfor

	// Make note of the elapsed time, as reported by go test
	results.Summary.setElapsed(GtpElapsed(jloSlice[len(jloSlice)-1].getElapsed()))

	// We've completed the for loop,
	// The last emitted line (JSON Line Object) announces
	// if the run as a whole was a pass or fail.  It does
	// not represent a test, but get counted as one.
	// So it throws off our counts by one.
	// So we fix that here
	results.Counts["pass"], results.Counts["fail"] =
		adjustOutSuperfluousFinalResult(jloSlice[len(jloSlice)-1].getAction(), results)
	// Now we check for results.Errors and create a
	// yellow bar and  message if appropriate
	results.buildBarMessage(Barmessage, PackageDirsToSearch)
}

// Shellout - run a command, capturing stdout, stderr, and errors
func Shellout(command string) (string, string, error) {
	// Force POSIX compliant shell for predictability
	// var ShellToUse = "/bin/sh"
	var ShellToUse = "/bin/sh"
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
func HandleOutputLines(results *GtpResults, jloSlice []JLObject, i int,
	PackageDir, pluginDir string, Barmessage *BarMessage) (bool, error) {

	var err error = nil
	doBreak := false

	results.incCount("output")

	doBreak = checkErrorCandidates(results, jloSlice[i].getOutput(), PackageDir)
	if doBreak {
		return doBreak, err
	}

	if hasTestCoverage(jloSlice[i].getOutput()) {
		results.Summary.setCoverage(jloSlice[i].getOutput())
	}

	// Error in an ExampleFunction()?
	if exampleError(jloSlice[i].getOutput()) {
		oneSpace := " "
		testName := jloSlice[i].getTest()
		exampleFuncDecl := fmt.Sprintf("func +%s\\(\\) +{ *", testName)

		filename, linenum, testname := findExampleFunc(pluginDir, exampleFuncDecl, PackageDir, results.GocycloIgnore)

		text := "Got: '" + jloSlice[i+2].getOutput() + "'" + oneSpace + "Want: '" + jloSlice[i+4].getOutput() + "'"

		if thisIsTheFirstFailure(results) {
			takeNoteOfFirstFailure(filename, linenum, jloSlice[i-1].getTest(), results)
		}

		// sometimes the paths can make the messages too long to fit on one
		// screen so, just use the filename
		filename = path.Base(filename)
		qfItem := buildQuickFixItem(PackageDir, filename, linenum, testname, text)
		Barmessage.QuickFixList.Add(qfItem)
		return doBreak, err
	}

	// If a jlo.Output field refers to a _test.go file, there has been a
	// test failure and it is telling us in which file and on which line
	// the failure was triggered
	if hasTestFileReferences(jloSlice[i].getOutput()) {
		length := len(jloSlice)
		oneSpace := " "
		list := splitOnColons(jloSlice[i].getOutput())
		// This may be obsolete, we will watch and see...
		filename := list[0]
		linenum := list[1]
		text := strings.Join(list[2:], "|")
		if i+1 < length-1 {
			secondLine := jloSlice[i+1].getOutput()
			if strings.HasPrefix(secondLine, "        ") {
				text += "|" + oneSpace + strings.TrimSpace(secondLine)
			}
		}
		testname := jloSlice[i-1].getTest()
		if thisIsTheFirstFailure(results) {
			takeNoteOfFirstFailure(filename, linenum, testname, results)
		}
		// sometimes the paths can make the messages too long to fit on one
		// screen so, just use the filename
		filename = path.Base(filename)
		qfItem := buildQuickFixItem(PackageDir, filename, linenum, testname, text)
		Barmessage.QuickFixList.Add(qfItem)
		return doBreak, err
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
	stdErrMsgSuffix := "[See pkgdir/StdErr.txt]"
	Barmessage.setColor("yellow")
	if stdErrMsgTooLongForOneLine(stderr, stdErrMsgPrefix, stdErrMsgSuffix, results.VimColumns) {
		writeStdErrMsgToDisk(stderr, PackageDirsToSearch[0])
		Barmessage.setMessage(buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgSuffix, msg, results.VimColumns))
	} else {
		Barmessage.setMessage(stdErrMsgPrefix + oneSpace + strings.ReplaceAll(msg, "\n", "|"))
		Barmessage.setMessage(strings.TrimSuffix(Barmessage.Message, "|"))
	}
	gtperror := GtpError{Name: "StdErrError", Regex: regexNil, Message: Barmessage.Message, Color: "yellow"}
	results.Errors = append(results.Errors, gtperror)
}

func buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgSuffix, msg string, cols int) string {
	oneSpace := " "
	commaSpace := ", "
	tmsg := strings.Split(msg, "\n")
	// So far I haven't found much use for the first lines which start with #
	// So we skip them to conserve space on the barMessage
	if strings.HasPrefix(tmsg[0], "#") {
		tmsg = tmsg[1:]
	}
	retMsg := strings.Join(tmsg, "|")
	retMsg = stdErrMsgPrefix + oneSpace + retMsg
	retMsg = strings.TrimSuffix(retMsg, "|")
	retMsg = retMsg[0 : cols-(len(stdErrMsgPrefix)+len(stdErrMsgSuffix))]
	retMsg += commaSpace + stdErrMsgSuffix
	return retMsg
}

func stdErrMsgTooLongForOneLine(stderr, stdErrMsgPrefix, stdErrMsgSuffix string, cols int) bool {
	oneSpace := " "
	return (len(stderr) > (cols - (len(stdErrMsgSuffix) + len(stdErrMsgPrefix) + len(oneSpace))))
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

func takeNoteOfFirstFailure(filename, linenum, testName string, results *GtpResults) {
	results.FirstFail.setFname(filename)
	results.FirstFail.setLineno(linenum)
	results.FirstFail.setTname(testName)
}

func splitOnColons(output string) []string {
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
		// order for these next two are important, the last one would also match the first
		{Name: "Panic", Regex: regexTimeOutPanic, Message: "In package: " + PackageDir + ", [Received a Test Time Out Panic]", Color: "yellow"},
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
	return CheckRegx(regexTestCoverage, output) || CheckRegx(regexCoverageNoStmts, output)
}

func hasTestFileReferences(output string) bool {
	// one of the surest fail indicators is an output about a "_test.go" file
	return CheckRegx(regexTestFileRef, output)
}

func setDebug(args []string) bool {
	debug := false
	if len(args) > 4 {
		if args[4] == "true" {
			debug = true
		}
	}
	return debug
}

func setupLogging() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("go-tdd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("Logging initiated.")
}
