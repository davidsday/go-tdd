package main

import (
	"fmt"
	"log"
	"strings"
)

func exampleError(output string) bool {
	return CheckRegx(regexExampleFail, output)
}

func findExampleFunc(pluginDir, exampleFuncDecl, path, ignore string) (string, string, string) {
	cmdLine := fmt.Sprintf(
		"%s/bin/ag --vimgrep -G '.*_test.go' --ignore '%s' '%s' %s",
		pluginDir,
		ignore,
		exampleFuncDecl,
		path,
	)
	log.Printf("In findExampleFunc, cmdLine: %s\n", cmdLine)
	result, _, err := Shellout(cmdLine)
	chkErr(err, "Error in ag searching for an example func declaration")
	trimmed := strings.TrimSuffix(result, "() {")
	split := splitOnColons(trimmed)
	split[3] = strings.TrimPrefix(split[3], "func ")
	// filename, lineno, testName
	return split[0], split[1], split[3]
}
