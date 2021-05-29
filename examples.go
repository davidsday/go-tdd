package main

import (
	"fmt"
)

func exampleError(output string) bool {
	return CheckRegx(regexExampleFail, output)
}

func findExampleFunc(pluginDir, exampleFuncDecl, path, ignore string) (filename, linenum, testname string) {
	cmdLine := fmt.Sprintf("%s/bin/ag --vimgrep -G '.*_test.go' --ignore '%s' '%s' %s", pluginDir, ignore, exampleFuncDecl, path)
	// log.Printf("In findExampleFunc, cmdLine: %s\n", cmdLine)
	out, _, err := Shellout(cmdLine)
	chkErr(err, "Error in ag searching for an example func declaration")
	split := splitOnColons(out)
	return split[0], split[1], split[3]
}
