package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func setupLogging() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("go-tdd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Printf("%s\n\n", "Logging initiated...")
}

func chkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("Error: %v, %s", err, msg)
	}
}

// Shellout - run a command, capturing stdout, stderr, and errors
func Shellout(command string) (string, string, error) {
	// Force POSIX compliant shell for predictability
	// var ShellToUse = "/bin/sh"
	shellPath, _ := exec.LookPath("sh")
	var ShellToUse = shellPath
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
} //end_Shellout()

func readFile(fPath string) string {
	data, err := ioutil.ReadFile(fPath)
	emsg := fmt.Sprintf("Error in readFile(%s)", fPath)
	chkErr(err, emsg)
	return string(data)
}
