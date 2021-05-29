package main

import (
	"os"
	"strings"
)

func exampleError(output string) bool {
	return CheckRegx(regexExampleFail, output)
}

func findExampleFunc(pluginDir, exampleFuncDecl, path string) (filename, linenum, testname string) {
	oneSpace := " "
	curDir, _ := os.Getwd()
	os.Chdir(path)
	cmdLine := pluginDir + "/bin/ag  --vimgrep --ignore " + "testdata" + oneSpace + exampleFuncDecl + oneSpace + path + `/*.go`
	// cmdLine := pluginDir + "/bin/ag  --vimgrep --ignore testdata" + oneSpace + exampleFuncDecl + oneSpace + path + `/*.go`
	out, _, err := Shellout(cmdLine)
	os.Chdir(curDir)
	chkErr(err, "Error in ag searching for an example func declaration")
	return splitExampleFuncSearchResults(out)
}

func splitExampleFuncSearchResults(result string) (filename, linenum, testname string) {
	trimmed := strings.TrimSuffix(result, "() {")
	split := splitOnColons(trimmed)
	split[3] = strings.TrimPrefix(split[3], "func ")
	return split[0], split[1], split[3]
}
