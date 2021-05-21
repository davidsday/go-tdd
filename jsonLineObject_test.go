package main

import (
	"encoding/json"
	"testing"
)

//TestgetOutput() ....
func TestGetOutput(t *testing.T) {
	jlo := JLObject{}
	jsonline := []byte(`{"Time":"2021-05-08T22:19:07.615912067-04:00","Action":"output","Package":"github.com/davidsday/go-tdd","Output":"PASS\n"}`)
	err := json.Unmarshal(jsonline, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	got := jlo.getOutput()
	want := "PASS\n"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestgetPackage() ....
func TestGetPackage(t *testing.T) {
	jlo := JLObject{}
	jsonline := []byte(`{"Time":"2021-05-08T22:19:07.615912067-04:00","Action":"output","Package":"github.com/davidsday/go-tdd","Output":"PASS\n"}`)
	err := json.Unmarshal(jsonline, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	got := jlo.getPackage()
	want := "github.com/davidsday/go-tdd"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestgetAction() ....
func TestGetAction(t *testing.T) {
	jlo := JLObject{}
	jsonline := []byte(`{"Time":"2021-05-08T22:19:07.615912067-04:00","Action":"output","Package":"github.com/davidsday/go-tdd","Output":"PASS\n"}`)
	err := json.Unmarshal(jsonline, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	got := jlo.getAction()
	want := "output"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestgetElapsed() ....
func TestGetElapsed(t *testing.T) {
	jlo := JLObject{}
	jsonline := []byte(`{"Time":"2021-05-08T22:19:07.616169473-04:00","Action":"pass","Package":"github.com/davidsday/go-tdd","Elapsed":0.005}`)
	err := json.Unmarshal(jsonline, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	got := jlo.getElapsed()
	want := float32(0.005)
	if got != want {
		t.Errorf("got '%f' want '%f'", got, want)
	}
}

//TestgetTest() ....
func TestGetTest(t *testing.T) {
	jlo := JLObject{}
	jsonline := []byte(`{"Time":"2021-05-08T22:19:07.615458883-04:00","Action":"pass","Package":"github.com/davidsday/go-tdd","Test":"TestGetAverageCyclomaticComplexity_no_go_files","Elapsed":0}`)
	err := json.Unmarshal(jsonline, &jlo)
	chkErr(err, "Error Unmarshaling jsonLine")
	got := jlo.getTest()
	want := "TestGetAverageCyclomaticComplexity_no_go_files"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
