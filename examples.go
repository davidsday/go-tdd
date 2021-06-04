package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func exampleError(output string) bool {
	return CheckRegx(regexExampleFail, output)
}

func findExampleFunc(pluginDir, exampleFuncDecl, pathToSearch, ignore string) (string, string, string) {
	agExecutable := "ag"
	if runtime.GOOS == "windows" {
		agExecutable = "ag.exe"
	}
	agBinaryPath := filepath.Join(pluginDir, "bin", agExecutable)
	cmdLine := fmt.Sprintf(
		"%s --vimgrep -G '.*_test.go' --ignore '%s' '%s' %s",
		agBinaryPath,
		ignore,
		exampleFuncDecl,
		pathToSearch,
	)
	result, _, err := Shellout(cmdLine)
	chkErr(err, "Error in ag searching for an example func declaration")
	trimmed := strings.TrimSuffix(result, "() {")
	split := splitOnColons(trimmed)
	split[3] = strings.TrimPrefix(split[3], "func ")
	// filename, lineno, testName
	return split[0], split[1], split[3]
}
