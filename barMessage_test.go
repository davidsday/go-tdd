package main

import (
	"path/filepath"
	"testing"
)

//TestBarMessage_setColor ....
func TestBarMessage_setColor(t *testing.T) {
	barmessage := BarMessage{}
	barmessage.setColor("red")
	got := barmessage.getColor()
	want := "red"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//TestBarMessage_setColor ....
func TestBarMessage_setMessage(t *testing.T) {
	barmessage := BarMessage{}
	barmessage.setMessage("Hello World")
	got := barmessage.getMessage()
	want := "Hello World"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//===========================================================================
// func marshalTR(barmessage Barmessage)
//===========================================================================

// func TestMarshallTR() ....
func Example1() {
	barmessage := BarMessage{}
	barmessage.setColor("yellow")
	barmessage.setMessage("Hello World!")

	barmessage.marshalToStdOut()
	// Output:  {"color":"yellow","message":"Hello World!","quickfixlist":null}
}

// func TestMarshallTR() ....
func Example2() {
	barmessage := BarMessage{}
	barmessage.setColor("yellow")
	barmessage.setMessage("Hello World!")
	qfItem := GtpQfItem{}
	barmessage.QuickFixList.Add(qfItem)
	barmessage.QuickFixList[0].Filename = "Something_test.go"
	barmessage.QuickFixList[0].Lnum = 12
	barmessage.QuickFixList[0].Col = 1
	barmessage.QuickFixList[0].Vcol = 1
	barmessage.QuickFixList[0].Pattern = "thisTest"
	barmessage.QuickFixList[0].Text = "Hello World!"

	barmessage.marshalToStdOut()
	// Output:  {"color":"yellow","message":"Hello World!","quickfixlist":[{"filename":"Something_test.go","lnum":12,"col":1,"vcol":1,"pattern":"thisTest","text":"Hello World!"}]}
}

//TestQfList.Add() ....
func TestQfList_Add(t *testing.T) {
	QfList := GtpQfList{}
	QfItem := GtpQfItem{}
	QfItem.Filename = "Something_test.go"
	QfList.Add(QfItem)
	got := QfList[0].Filename
	want := "Something_test.go"
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

//============================================================================
// func buildQuickFixItem(args []string, parts []string, jlo JLObject) GtpQfItem {
//============================================================================
//TestAddToQuickFixList_length ....
func TestBuildQuickFixItem_filename(t *testing.T) {
	QfItem := GtpQfItem{}
	filename := "filename"
	linenum := "12"
	pattern := "thisTest"
	text := "text"
	pkgDir := "packageDir"
	jlo := JLObject{}
	jlo.Test = "thisTest"
	jlo.Package = "packageDir"

	QfItem = buildQuickFixItem(pkgDir, filename, linenum, pattern, text)
	if QfItem.Filename != "packageDir"+string(filepath.Separator)+"filename" {
		t.Errorf("Filename:  Got: %s, Want: %s\n", QfItem.Filename, "packageDir"+string(filepath.Separator)+"filename")
	}
}

//TestAddToQuickFixList_pattern ....
func TestBuildQuickFixItem_pattern(t *testing.T) {
	QfItem := GtpQfItem{}
	filename := "filename"
	linenum := "12"
	pattern := "thisTest"
	text := "text"
	pkgDir := "packageDir"
	jlo := JLObject{}
	jlo.Test = "thisTest"

	QfItem = buildQuickFixItem(pkgDir, filename, linenum, pattern, text)
	if QfItem.Pattern != "thisTest" {
		t.Errorf("Pattern:  Got: %s, Wanted: %s\n", QfItem.Pattern, "thisTest")
	}
}

//============================================================================
// func QfList.Count()
//============================================================================

//TestQfList.Count ....
func TestQfList_Count_2(t *testing.T) {
	QfList := GtpQfList{}
	QfItem1 := GtpQfItem{}
	QfItem2 := GtpQfItem{}
	QfList.Add(QfItem1)
	QfList.Add(QfItem2)
	got := QfList.Count()
	want := 2
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}

//TestQfList.Count ....
func TestQfList_Count_0(t *testing.T) {
	QfList := GtpQfList{}
	got := QfList.Count()
	want := 0
	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}
