package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
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
var PackageDir string

func main() {

	// results has all the results we collect from go test
	// to help us decide how to present the results to the user
	// It has the methods it needs to build the BarMessage
	// It lives in results.go
	// We have built a func, newResults(), which creates, initializes
	// and returns a new results to us
	results := newResults()

	// barMessage holds the message and background color Vim will
	// display plus the quickfix list, QfList. All of these are populated
	// by the methods in Results.   A barMessage can marshal itself
	// into JSON and send it to stdout for Vim to display
	// or to disk for logging purposes, or into a []bytes.
	// BarMessage lives in barMessage.go

	barMessage := newBarMessage()

	// Instead of calling go-tdd with a fairly long list of command line
	// arguments, put the args into a structure and marshal them to JSON,
	// passing just that one JSON encoded string on the command line.
	// We are storing our JSON arguments in the Args structure embedded
	// in our results struct, which calculates our results for us.
	// It needs several of these arguments to do that.
	err := json.Unmarshal([]byte(os.Args[1]), &results.Args)
	chkErr(err, "Error in json.Unmarshal of os.Args[1]")

	// We get quidance from Vim about where go test and gocyclo
	// should search, there is really only one dir from Vim,
	// but gocyclo wants a list of dirs, so we create an empty
	// list and append the dir we got from Vim to it so
	// gocyclo will be happy
	var packageDirsToSearch []string
	packageDirsToSearch = append(packageDirsToSearch, results.Args.PackageDir)

	debug = results.Args.GoTddDebug
	if debug {
		setupLogging()
		log.Printf("debug: '%#v', type: '%T'\n\n", debug, debug)
		log.Printf("os.Args[1] '%v'\n\n", os.Args[1])
		log.Printf("results.Args: '%#v'\n\n", results.Args)
	}

	oneSpace := " "
	goTestTimeout := results.Args.Timeout

	commandLine := "go test -v -json -cover"
	commandLine += oneSpace + packageDirsToSearch[0]
	commandLine += oneSpace + "-timeout"
	commandLine += oneSpace + goTestTimeout

	errString := fmt.Sprintf("error running '%s'\n", commandLine)
	stdout, stderr, err := Shellout(commandLine)
	chkErr(err, errString)

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
			doBreak, err = HandleOutputLines(results, jloSlice, i, PackageDir, results.Args.PluginDir, Barmessage)
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
		if debug {
			log.Printf("About to call findExampleFunc(): results.Args.PluginDir: '%s'\n\n", results.Args.PluginDir)
			log.Printf("In findExampleFunc(): pluginDir: '%s'\n\n", pluginDir)
		}
		filepath, linenum, testname := findExampleFunc(pluginDir, exampleFuncDecl, PackageDir, results.Args.GocycloIgnore)

		text := "Got: '" + jloSlice[i+2].getOutput() + "'" + oneSpace + "Want: '" + jloSlice[i+4].getOutput() + "'"

		if thisIsTheFirstFailure(results) {
			takeNoteOfFirstFailure(filepath, linenum, jloSlice[i-1].getTest(), results)
		}

		// sometimes the paths can make the messages too long to fit on one
		// screen so, just use the filename
		filename := path.Base(filepath)
		qfItem := buildQuickFixItem(PackageDir, filename, linenum, testname, text)
		Barmessage.QuickFixList.Add(qfItem)
		return doBreak, err
	}

	// If a jlo.Output field refers to a _test.go file, there has been a
	// test failure and it is telling us in which file and on which line
	// the failure was triggered
	if hasTestFileReferences(jloSlice[i].getOutput()) {
		oneSpace := " "
		indent := "    "
		secondLine := ""
		list := splitOnColons(jloSlice[i].getOutput())
		filename := list[0]
		linenum := list[1]
		text := strings.Join(list[2:], " | ")
		// check that we are not reaching past the end of jloSlice
		if safeToLookAhead(jloSlice, i, 1) {
			secondLine = jloSlice[i+1].getOutput()
			if strings.HasPrefix(secondLine, indent+indent) {
				text += oneSpace + "|" + oneSpace + strings.TrimSpace(secondLine)
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
	if stdErrMsgTooLongForOneLine(stderr, stdErrMsgPrefix, stdErrMsgSuffix, results.Args.getScreenColumns()) {
		writeStdErrMsgToDisk(stderr, PackageDirsToSearch[0])
		Barmessage.setMessage(buildShortenedBarMessage(stdErrMsgPrefix, stdErrMsgSuffix, msg, results.Args.getScreenColumns()))
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
	path := "./StdErr.txt"
	if len(strings.TrimSpace(pkgdir)) > 0 {
		pkgdir = strings.TrimSuffix(pkgdir, "/")
		path = pkgdir + "/StdErr.txt"
	}
	err := os.WriteFile(path, []byte(stderr), 0664)
	chkErr(err, "error writing "+path)
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

func safeToLookAhead(jloSlice []JLObject, i, incr int) bool {
	length := len(jloSlice)
	return length-1 >= i+incr
}
