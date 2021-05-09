package main

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

func (j *JLObject) getOutput() string {
	return j.Output
}

func (j *JLObject) getAction() string {
	return j.Action
}

func (j *JLObject) getTest() string {
	return j.Test
}

func (j *JLObject) getPackage() string {
	return j.Package
}

func (j *JLObject) getElapsed() float32 {
	return j.Elapsed
}
