package main

import (
	"fmt"
	"strings"
)

func exampleError(output string) bool {
	return CheckRegx(regexExampleFail, output)
}

func findExampleFunc(pluginDir, exampleFuncDecl, path string) (filename, linenum, testname string) {
	// curDir, _ := os.Getwd()
	// os.Chdir(path)
	cmdLine := fmt.Sprintf("%s/bin/ag --vimgrep -G '.*_test.go' --ignore '%s' '%s' %s", pluginDir, `testdata|vendor`, exampleFuncDecl, `/home/dave/sw/go/go-tdd`)
	// cmdLine := fmt.Sprintf("%s/bin/ag --vimgrep -G '.*_test.go' --ignore '%s' '%s' %s", pluginDir, `testdata|vendor`, exampleFuncDecl, path)
	// log.Printf("In findExampleFunc, cmdLine: %s\n", cmdLine)
	out, _, err := Shellout(cmdLine)
	// os.Chdir(curDir)
	chkErr(err, "Error in ag searching for an example func declaration")
	return splitExampleFuncSearchResults(out)
}

func splitExampleFuncSearchResults(result string) (filename, linenum, testname string) {
	trimmed := strings.TrimSuffix(result, "() {")
	split := splitOnColons(trimmed)
	split[3] = strings.TrimPrefix(split[3], "func ")
	return split[0], split[1], split[3]
}
