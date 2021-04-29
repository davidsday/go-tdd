package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	// regexAvgComplexity  = regexp.MustCompile(`Average: \d{1,2}\.\d{1,2}`)
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
var jlo JLObject
var prevJlo JLObject

// PackageDir is where the current package lives
var PackageDir string

func main() {
	// oneSpace := " "
	commaSpace := ", "

	// Remove ./StdeErr.txt if one still lingers
	// and any lingering JSON logs...
	// So we are guaranteed any present after
	// this run are created by this run
	// os.Remove("./StdErr.txt")
	// os.Remove("./goTestParser_log.json")
	// os.Remove("./gotestlog.json")

	commandLine := "go test -v -json -cover " + os.Args[1]
	// New structs are initialized empty (false, 0, "", [], {} etc)
	// A few struct members need to have different initializations
	// So we take care of that here
	// We will assume we are receiving valid JSON, until we find
	// an invalid JSON Line Object
	PD.Perror.Validjson = true
	// Vim/Neovim knows how many screen columns it has
	// and passes that knowledge to us via os.Args[2]
	// so we can tailor our messages to fit on one screen line
	cols, _ := strconv.Atoi(os.Args[2])
	PD.Barmessage.Columns = cols

	// General test run info is in PD.Info
	PD.Info.Host, _ = os.Hostname()
	PD.Info.GtpIssuedCmd = commandLine
	PD.Info.Begintime = time.Now().Format(time.RFC3339Nano)
	// PD.Info.Endtime is set just before finishing up, down below
	PD.Info.User = os.Getenv("USER")
	// time.Now().Format(time.RFC3339Nano)
	PD.Info.GtpRcvdArgs = os.Args

	// If os.Args[2] == "--" {
	//    open stdin and read from it until EOF
	// } and our stdin becomes our variable stdout, here
	// Might have to reconsider our naming, eh???
	stdout, stderr, _ := Shellout(commandLine)
	if len(stderr) > 0 {
		msg := ""
		PD.Perror.MsgStderr = true
		PD.Barmessage.Color = "yellow"
		if len(stderr) > 100 {
			msg = stderr[:90]
		} else {
			msg = stderr
		}
		PD.Barmessage.Message = "STDERR: " + strings.ReplaceAll(msg, "\n", "|")
		if len(stderr) > PD.Barmessage.Columns-26 {
			path := PackageDir + "/StdErr.txt"
			err := os.WriteFile(path, []byte(stderr), 0664)
			if err != nil {
				log.Fatal("Error writing pkgfile/StdErr.txt")
			}
			// PD.Barmessage.Message = PD.Barmessage.Message[0 : PD.Barmessage.Columns-26]
			PD.Barmessage.Message += commaSpace + "[See pkgdir/StdErr.txt]"
		}
	} else {
		// stdout & stderr are strings, we need []byte
		lines := bytes.Split([]byte(stdout), []byte("\n"))

		for _, jsonLine := range lines[:len(lines)-1] {

			if len(bytes.TrimSpace(jsonLine)) == 0 {
				continue
			}

			// Ensure we're getting valid JSON
			if !json.Valid(jsonLine) {
				PD.Perror.Validjson = false
				break
			} else {
				// Convert line of JSON text to JSON line object (Go struct in this case)
				err := json.Unmarshal(jsonLine, &jlo)
				if err != nil {
					log.Fatal("Error Unmarshalling jsonLine ")
				}
			}

			PackageDir = jlo.Package

			if jlo.Action == "run" {
				PD.Counts.Runs++
			} else if jlo.Action == "continue" {
				PD.Counts.Continues++
			} else if jlo.Action == "pause" {
				PD.Counts.Pauses++
			} else if jlo.Action == "skip" {
				PD.Counts.Skips++
			} else if jlo.Action == "pass" {
				PD.Counts.Passes++
			} else if jlo.Action == "fail" {
				PD.Counts.Fails++
			}
			var err error
			var doBreak bool

			PD, doBreak, err = HandleOutputLines(PD, jlo, prevJlo, PackageDir)
			if err != nil {
				os.Exit(1)
			}
			if doBreak {
				break
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
		if jlo.Action == "pass" {
			if PD.Counts.Passes > 1 {
				PD.Counts.Passes--
			}
		}
		if jlo.Action == "fail" {
			if PD.Counts.Fails > 1 {
				PD.Counts.Fails--
			}
		}
		// Now we cycle through our PD.Error flags and create a
		// yellow bar and  message if appropriate

		if !PD.Perror.Validjson {
			PD.Barmessage.Color = "yellow"
			PD.Barmessage.Message = "In package: " + PackageDir + ", [Found Invalid JSON]"
		} else if PD.Perror.Notestfiles {
			PD.Barmessage.Color = "yellow"
			PD.Barmessage.Message = "In package: " + PackageDir + ", [No Test Files]"
		} else if PD.Perror.Noteststorun {
			PD.Barmessage.Color = "yellow"
			PD.Barmessage.Message = "In package: " + PackageDir + ", [Test Files, but No Tests to Run]"
		} else if PD.Perror.Buildfailed {
			PD.Barmessage.Color = "yellow"
			PD.Barmessage.Message = "In package: " + PackageDir + ", [Build Failed]"
		} else if PD.Perror.RcvPanic {
			PD.Barmessage.Color = "yellow"
			PD.Barmessage.Message = "In package: " + PackageDir + ", [Received a Panic]"
		} else {
			// No errors above so if we have fails or skips we load the quickfixlist
			// and select "red" as our color bar color
			if PD.Counts.Fails > 0 {
				PD.Barmessage.Color = "red"
			} else if PD.Counts.Skips > 0 {
				PD.Barmessage.Color = "yellow"
			} else {
				PD.Barmessage.Color = "green"
				// Since we only show avg cyclomatic complexity on green bars,
				var err error
				PD.Info.AvgComplexity, err = getAvgCyclomaticComplexity(PackageDir)
				if err != nil {
					log.Fatalf("%s, exiting", err)
				}
			}
			// func BuildBarMessage(
			// runs int,
			// skips int,
			// fails int,
			// passes int,
			// elapsed float32,
			// fname string,
			// lineno int) string {}

			PD.Barmessage.Message = BuildBarMessage(
				PD.Counts.Runs,
				PD.Counts.Skips,
				PD.Counts.Fails,
				PD.Counts.Passes,
				PD.Elapsed,
				PD.Firstfailedtest.Fname,
				PD.Firstfailedtest.Lineno,
				PD.Info.TestCoverage,
				PD.Info.AvgComplexity,
			)
		}
	}

	// Endtime for PD.Info
	PD.Info.Endtime = time.Now().Format(time.RFC3339Nano)
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

// function to perform marshalling
func marshallTR(pgmdata PgmData) {

	// data, err := json.MarshalIndent(pgmdata, "", "    ")
	data, _ := json.Marshal(pgmdata)

	_, err := os.Stdout.Write(data)
	if err != nil {
		log.Fatal("Error writing to Stdout")
	}

	// err = os.WriteFile("./goTestParserLog.json", data, 0664)
	// if err != nil {
	//	log.Fatal("Error writing to ./goTestParserLog.json")
	// }

} // end_marshallTR

// HandleOutputLines does the regular expression checking and
// to discern what is happening
func HandleOutputLines(pgmdata PgmData, jlo JLObject, prevJlo JLObject,
	PackageDir string) (PgmData, bool, error) {
	var tDict PDQfDict
	var err error = nil
	var parts []string
	var text string
	doBreak := false
	pgmdata.Counts.Outputs++

	if CheckRegx(regexPanic, jlo.Output) {
		pgmdata.Perror.RcvPanic = true
		doBreak = true
		return pgmdata, doBreak, err
	}

	if CheckRegx(regexTestCoverage, jlo.Output) {
		// Remove trailing '\n'
		pgmdata.Info.TestCoverage = strings.TrimSuffix(jlo.Output, "\n")
		// Strip away everything but the percent coverage string ("57.8%", for example)
		pgmdata.Info.TestCoverage = strings.Replace(pgmdata.Info.TestCoverage, "coverage: ", "", 1)
		pgmdata.Info.TestCoverage = strings.Replace(pgmdata.Info.TestCoverage, " of statements", "", 1)
	}

	if CheckRegx(regexNoTestFiles, jlo.Output) {
		pgmdata.Perror.Notestfiles = true
		doBreak = true
		return pgmdata, doBreak, err
	}

	if CheckRegx(regexNoTestsToRun, jlo.Output) {
		pgmdata.Perror.Noteststorun = true
		doBreak = true
		return pgmdata, doBreak, err
	}

	if CheckRegx(regexBuildFailed, jlo.Output) {
		pgmdata.Perror.Buildfailed = true
		doBreak = true
		return pgmdata, doBreak, err
	}
	if CheckRegx(regexFailorTestFile, jlo.Output) {
		parts = strings.Split(strings.TrimSpace(jlo.Output), ":")
		if strings.Contains(parts[0], "FAIL:") {
			// then the 1st element is " FAIL:"
			// so we take the sublist and continue our work
			parts = parts[1:]
		}
		if pgmdata.Counts.Fails == 0 {
			pgmdata.Firstfailedtest.Fname = parts[0]
			pgmdata.Firstfailedtest.Lineno = parts[1]
			pgmdata.Firstfailedtest.Tname = prevJlo.Test
		}
		if len(parts) > 2 {
			text = strings.Join(parts[2:], ":")
		} else {
			text = "xxx"
		}
		// pgmdata.Counts.Fails++
		// Now we can build/fill the QuickFix List
		// tDict.Filename = PackageDir + "/" + parts[0]
		tDict.Filename = os.Args[1] + "/" + parts[0]
		// tDict.Filename = parts[0]
		tDict.Lnum, _ = strconv.Atoi(parts[1])
		tDict.Col = 1
		tDict.Vcol = 1
		tDict.Pattern = jlo.Test
		tDict.Text = text
		pgmdata.Qflist = append(pgmdata.Qflist, tDict)
		// Should already be false, since that is how it was initialized
		doBreak = false
	}

	err = nil
	return pgmdata, doBreak, err
}

// BuildBarMessage dynamically creates the message for passed,
// failed, and skipped tests as appropriate
//
// Given the relevent counters, the elapsed time, a possible 1st error
// filename and line number, return the completed message
func BuildBarMessage(runs int, skips int, fails int, passes int, elapsed PDElapsed, fname string, lineno string, coverage string, complexity string) string {
	oneSpace := " "
	commaSpace := ", "
	barmessage := strconv.Itoa(runs) + oneSpace + "Run"
	barmessage += commaSpace + strconv.Itoa(passes) + oneSpace + "Passed"
	if skips > 0 {
		barmessage += commaSpace + strconv.Itoa(skips) + oneSpace + "Skipped"
	}
	if fails > 0 {
		barmessage += commaSpace + strconv.Itoa(fails) + oneSpace + "Failed"
		barmessage += commaSpace + "1st in" + oneSpace + fname
		barmessage += commaSpace + "on line" + oneSpace + lineno
	}
	if skips == 0 && fails == 0 && len(coverage) > 0 {
		barmessage += commaSpace + "Test Coverage:" + oneSpace + coverage
		barmessage += commaSpace + "Average Complexity:" + oneSpace + complexity
	}
	barmessage += commaSpace + "in" + oneSpace + strconv.FormatFloat(float64(elapsed), 'f', 3, 32) + "s"
	return barmessage
}

// CheckRegx check for a match described by compiled regx with candidate.
// Returns true if theres a match, false otherwise
func CheckRegx(regx *regexp.Regexp, candidate string) bool {
	match := regx.FindString(candidate)
	return match != ""
}

func getAvgCyclomaticComplexity(path string) (string, error) {
	oneSpace := " "
	avgCmplxCmdLine := "gocyclo -avg -ignore 'vendor|_test.go'" + oneSpace + path + oneSpace + "| grep 'Average: ' | awk '{print $2}'"
	sout, _, err := Shellout(avgCmplxCmdLine)
	sout = strings.TrimSuffix(sout, "\n")
	if err != nil {
		myerr := errors.New("error getting cyclomatic complexity")
		sout = ""
		return sout, myerr
	}
	return sout, err
}
