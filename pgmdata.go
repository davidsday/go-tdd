package main

import (
	"os"
	"regexp"
	"strconv"
	"time"
)

type PgmData struct {
	Info            PDInfo            `json:"info"`
	Counts          PDCounts          `json:"counts"`
	Firstfailedtest PDFirstFailedTest `json:"firstfailedtest"`
	Elapsed         PDElapsed         `json:"elapsed"`
	Perrors         GTPerrors         `json:"errors"`
	QfList          PDQfList          `json:"qflist"`
	Barmessage      PDBarMessage      `json:"barmessage"`
}

func (p *PgmData) setBarMessage() {
	if len(PD.Perrors) > 0 {
		PD.Barmessage.Color = PD.Perrors[0].Color
		PD.Barmessage.Message = PD.Perrors[0].Message
	} else {
		if p.Counts["fail"] > 0 {
			p.Barmessage.Color = "red"
		} else if p.Counts["skip"] > 0 {
			p.Barmessage.Color = "yellow"
		} else {
			p.Barmessage.Color = "green"
			// Since we only show avg cyclomatic complexity on green bars,
			// only run it for green bars
			p.Info.AvgComplexity = getAvgCyclomaticComplexity(PackageDirFromVim)
		}

		barmessage := runMsg(p.Counts["run"])
		barmessage += passMsg(p.Counts["pass"])
		barmessage += skipMsg(p.Counts["skip"])
		barmessage += failMsg(p.Counts["fail"], p.Firstfailedtest.Fname, p.Firstfailedtest.Lineno)
		barmessage += metricsMsg(p.Counts["skip"], p.Counts["fail"], p.Info.TestCoverage, p.Info.AvgComplexity)
		barmessage += elapsedMsg(p.Elapsed)
		p.Barmessage.Message = barmessage
	}
}

type PDInfo struct {
	Host          string   `json:"host"`
	User          string   `json:"user"`
	Begintime     string   `json:"begintime"`
	Endtime       string   `json:"endtime"`
	GtpIssuedCmd  string   `json:"gtp_issued_cmd"`
	GtpRcvdArgs   []string `json:"gtp_rcvd_args"`
	TestCoverage  string   `json:"test_coverage"`
	AvgComplexity string   `json:"avg_complexity"`
}

func (p *PgmData) initializePgmData(commandLine string) {

	// New structs are initialized empty (false, 0, "", [], {} etc)
	// A few struct members need to have different initializations
	// So we take care of that here
	// We will assume we are receiving valid JSON, until we find
	// an invalid JSON Line Object
	p.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 0, "fail": 0, "output": 0}

	// Vim/Neovim knows how many screen columns it has
	// and passes that knowledge to us via os.Args[2]
	// so we can tailor our messages to fit on one screen line
	p.Barmessage.Columns, _ = strconv.Atoi(os.Args[2])

	// General info is held in PD.Info
	p.Info.Host, _ = os.Hostname()
	p.Info.GtpIssuedCmd = commandLine
	p.Info.Begintime = time.Now().Format(time.RFC3339Nano)
	// PD.Info.Endtime is set just before finishing up, down below
	p.Info.User = os.Getenv("USER")
	// goTestParser is started by vim
	// these are the args it received
	p.Info.GtpRcvdArgs = os.Args
}

type PDCounts map[string]int

type PDFirstFailedTest struct {
	Fname  string `json:"fname"`
	Tname  string `json:"tname"`
	Lineno string `json:"lineno"`
}

type PDElapsed float64

type GTPerror struct {
	Name    string         `json:"name"`
	Regex   *regexp.Regexp `json:"regexp"`
	Message string         `json:"message"`
	Color   string         `json:"color"`
}

type GTPerrors []GTPerror

type PDPerror struct {
	Validjson    bool `json:"validjson"`
	Notestfiles  bool `json:"notestfiles"`
	Noteststorun bool `json:"noteststorun"`
	RcvPanic     bool `json:"panic"`
	Buildfailed  bool `json:"buildfailed"`
	MsgStderr    bool `json:"msg_stderr"`
}

type PDQfList []PDQfDict

type PDBarMessage struct {
	Color   string `json:"color"`
	Message string `json:"message"`
	Columns int    `json:"columns"`
}

type PDQfDict struct {
	Filename string `json:"filename"`
	Lnum     int    `json:"lnum"`
	Col      int    `json:"col"`
	Vcol     int    `json:"vcol"`
	Pattern  string `json:"pattern"`
	Text     string `json:"text"`
}
