package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
var PD PgmData

// jlo & JLO -> JSON Line Object
var jlo JLObject
var prev_jlo JLObject
var PackageDir string

// var testResults TestResults
// var qfLines QuickFixLines
// var qfLine QuickFixLine

func main() {

	commandLine := "go test -v -json " + os.Args[1]
	// New structs are initialized empty (false, 0, "", [], {} etc)
	// A few struct members need to have different initializations
	// So we take care of that here
	PD.Perror.Validjson = true

	// General go test run info is in PD.Info
	PD.Info.Host, _ = os.Hostname()
	PD.Info.Commandline = commandLine
	PD.Info.Begintime = time.Now().Format(time.RFC3339Nano)
	// PD.Info.Endtime is set just before finishing up, down below
	PD.Info.User = os.Getenv("USER")
	// time.Now().Format(time.RFC3339Nano)

	// stdout, stderr, _ := Shellout(commandLine)
	stdout, _, _ := Shellout(commandLine)
	// if len(stderr) > 0 {
	// 	println("Oops!  Got something on stderr...Check it out!")
	// 	os.Exit(1)
	// }
	// stdout & stderr are strings, we need []byte
	lines := bytes.Split([]byte(stdout), []byte("\n"))

	for _, json_line := range lines[:len(lines)-1] {

		// Ensure we're getting valid JSON
		if !json.Valid(json_line) {
			PD.Perror.Validjson = false
			break
		} else {
			// Convert line of JSON text to JSON line object (Go struct in this case)
			json.Unmarshal(json_line, &jlo)
		}

		PackageDir = jlo.Package

		if jlo.Action == "run" {
			PD.Counts.Runs++
		}
		if jlo.Action == "continue" {
			PD.Counts.Continues++
		}
		if jlo.Action == "pause" {
			PD.Counts.Pauses++
		}
		if jlo.Action == "skip" {
			PD.Counts.Skips++
		}
		if jlo.Action == "pass" {
			PD.Counts.Passes++
		}
		if jlo.Action == "fail" {
			PD.Counts.Fails++
		}
		var err error
		var doBreak bool

		PD, doBreak, err = HandleOutputLines(PD, jlo, prev_jlo, PackageDir)
		if err != nil {
			os.Exit(1)
		}
		if doBreak {
			break
		}
		// Bottom of for loop - current JSON Line Object now
		// becomes the Previous JSON Line Object,
		// for look back purposes ...
		prev_jlo = jlo
	} //endfor

	// Make note of the elapsed time
	PD.Elapsed = PD_Elapsed(jlo.Elapsed)

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
	} else if PD.Perror.Rcv_panic {
		PD.Barmessage.Color = "yellow"
		PD.Barmessage.Message = "In package: " + PackageDir + ", [Received a Panic]"
	} else if PD.Perror.Notestfiles {
		PD.Barmessage.Color = "yellow"
		PD.Barmessage.Message = "In package: " + PackageDir + ", [No Test Files]"
	} else if PD.Perror.Buildfailed {
		PD.Barmessage.Color = "yellow"
		PD.Barmessage.Message = "In package: " + PackageDir + ", [Build Failed]"
	} else {
		// No errors above so if we have fails or skips we load the quickfixlist
		// and select "red" as our color bar color
		if PD.Counts.Fails > 0 || PD.Counts.Skips > 0 {
			PD.Barmessage.Color = "red"
		} else {
			PD.Barmessage.Color = "green"
		}
		// func BuildBarMessage(runs int, skips int, fails int, passes int, elapsed float32, fname string, lineno int) string {
		PD.Barmessage.Message = BuildBarMessage(PD.Counts.Runs, PD.Counts.Skips, PD.Counts.Fails, PD.Counts.Passes, PD.Elapsed, PD.Firstfailedtest.Fname, PD.Firstfailedtest.Lineno)
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

	os.Stdout.Write(data)
	os.WriteFile("./gotdd_log.json", data, 0664)
} // end_marshallTR

func HandleOutputLines(pgmdata PgmData, jlo JLObject, prev_jlo JLObject, PackageDir string) (PgmData, bool, error) {
	var tDict PD_QfDict
	var err error = nil
	var parts []string
	var text string
	doBreak := false
	pgmdata.Counts.Outputs++
	// ?   	github.com/zchee/nvim-go/pkg/server	[no test files]
	// var regexNoTestFiles = regexp.MustCompile(`\?.*\[no test files\]`)
	var regexNoTestFiles = regexp.MustCompile(`\?\s*\S*\s*\[no test files\]`)
	var regexBuildFailed = regexp.MustCompile(`\?\s*\S*\s*\[build failed\]`)
	// TODO: I don't yet know an adequate regex for this
	// panic: Log in goroutine after  has completed
	var regexPanic = regexp.MustCompile(`^panic:`)

	var regexFailorTestFile = regexp.MustCompile(`^\s\+FAIL:|_test.go`)

	// Shiny new match for each loop, nothing left over from
	// previous use
	match := regexPanic.FindString(jlo.Output)
	if match != "" {
		// fmt.Println("Found [panic]")
		// fmt.Printf("Found:  %s \n\n", match)
		pgmdata.Perror.Rcv_panic = true
		doBreak = true
		return pgmdata, doBreak, err
	}

	match = regexNoTestFiles.FindString(jlo.Output)
	if match != "" {
		// fmt.Println("WooHooh!!! Found [no test files]")
		// fmt.Printf("Found:  %s \n\n", match)
		pgmdata.Perror.Notestfiles = true
		doBreak = true
		return pgmdata, doBreak, err
	}

	match = regexBuildFailed.FindString(jlo.Output)
	if match != "" {
		// fmt.Println("Found [build failed]")
		// fmt.Printf("Found:  %s \n\n", match)
		pgmdata.Perror.Buildfailed = true
		doBreak = true
		return pgmdata, doBreak, err
	}
	match = regexFailorTestFile.FindString(jlo.Output)
	if match != "" {
		parts = strings.Split(strings.TrimSpace(jlo.Output), ":")
		if strings.Contains(parts[0], "FAIL:") {
			// then the 1st element is " FAIL:"
			// so we take the sublist and continue our work
			parts = parts[1:]
		}
		if pgmdata.Counts.Fails == 0 {
			pgmdata.Firstfailedtest.Fname = parts[0]
			pgmdata.Firstfailedtest.Lineno = parts[1]
			pgmdata.Firstfailedtest.Tname = prev_jlo.Test
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

func BuildBarMessage(runs int, skips int, fails int, passes int, elapsed PD_Elapsed, fname string, lineno string) string {
	barmessage := strconv.Itoa(runs) + " Run, " + strconv.Itoa(passes) + " Passed"
	if skips > 0 {
		barmessage += ", " + strconv.Itoa(skips) + " Skipped"
	}
	if fails > 0 {
		barmessage += ", " + strconv.Itoa(fails) + " Failed, 1st in " + fname + ", on line " + lineno
	}

	barmessage += ", in " + strconv.FormatFloat(float64(elapsed), 'f', 3, 32) + "s"

	return barmessage
}
