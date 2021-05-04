package main

import "regexp"

type PgmData struct {
	Info            PDInfo            `json:"info"`
	Counts          PDCounts          `json:"counts"`
	Firstfailedtest PDFirstFailedTest `json:"firstfailedtest"`
	Elapsed         PDElapsed         `json:"elapsed"`
	Perror          PDPerror          `json:"error"`
	Perrors         GTPerrors         `json:"errors"`
	QfList          PDQfList          `json:"qflist"`
	Barmessage      PDBarMessage      `json:"barmessage"`
}

func (p *PgmData) setBarMessage() {
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
