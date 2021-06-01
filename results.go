package main

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/fzipp/gocyclo"
)

//============================================================================
// Here are data types, structs and maps and methods
//============================================================================

type GtpCounts map[string]int

type GtpArgs struct {
	PackageDir    string `json:"package_dir"`
	ScreenColumns string `json:"screen_columns"`
	GocycloIgnore string `json:"gocyclo_ignore"`
	GoTddDebug    bool   `json:"go_tdd_debug"`
	PluginDir     string `json:"plugin_dir"`
	Timeout       string `json:"timeout"`
}

func (a *GtpArgs) getScreenColumns() int {
	cols, err := strconv.Atoi(a.ScreenColumns)
	chkErr(err, "Error converting results.Args.ScreenColumns (string) to int")
	return cols
}

func (a *GtpArgs) setScreenColumns(cols int) {
	a.ScreenColumns = strconv.Itoa(cols)
}

type GtpResults struct {
	Summary   GtpSummary
	Counts    GtpCounts
	Errors    GtpErrors
	FirstFail GtpFirstFail
	Args      GtpArgs
	// VimColumns    int
	GocycloIgnore string
}

func newResults() GtpResults {
	r := new(GtpResults)
	// Initialize map of Counts in Results
	r.Counts = map[string]int{"run": 0, "pause": 0, "continue": 0, "skip": 0, "pass": 0, "fail": 0, "output": 0}
	r.GocycloIgnore = `vendor|testdata`
	return *r
}

type GtpFirstFail struct {
	Fname  string `json:"fname"`
	Tname  string `json:"tname"`
	Lineno string `json:"lineno"`
}

func (f *GtpFirstFail) setFname(fname string) {
	f.Fname = fname
}

func (f *GtpFirstFail) setTname(tname string) {
	f.Tname = tname
}

func (f *GtpFirstFail) setLineno(lineno string) {
	f.Lineno = lineno
}

func (f *GtpFirstFail) getFname() string {
	return f.Fname
}

// func (f *GtpFirstFail) getTname() string {
//	return f.Tname
// }

func (f *GtpFirstFail) getLineno() string {
	return f.Lineno
}

type GtpError struct {
	Name    string         `json:"name"`
	Regex   *regexp.Regexp `json:"regexp"`
	Message string         `json:"message"`
	Color   string         `json:"color"`
}
type GtpErrors []GtpError

type GtpCoverage string
type GtpComplexity string
type GtpElapsed float64

// We build the Bar Message here
func (r *GtpResults) buildBarMessage(bm *BarMessage, PackageDirsToSearch []string) {
	if len(r.Errors) > 0 {
		// each error has its own message and color (right now, all yellow)
		bm.setColor(r.Errors[0].Color)
		bm.setMessage(r.Errors[0].Message)
		// we display one message at a time and we've already go one
		return
	}
	if r.getCount("fail") > 0 {
		// any fails -> red
		bm.setColor("red")
	} else if r.getCount("skip") > 0 {
		// no fails but skips -> yellow
		bm.setColor("yellow")
	} else {
		// all passed! -> green
		bm.setColor("green")
		// Since we only show avg cyclomatic complexity on green bars,
		// only run it for green bars
		r.Summary.setComplexity(PackageDirsToSearch, r.GocycloIgnore)
	}

	// build the message based on how we did ...
	msg := runMsg(r.getCount("run"))
	msg += passMsg(r.getCount("pass"))
	msg += skipMsg(r.getCount("skip"))
	msg += failMsg(r.getCount("fail"), r.FirstFail.getFname(), r.FirstFail.getLineno())
	msg += metricsMsg(r.getCount("skip"), r.getCount("fail"), string(r.Summary.getCoverage()), string(r.Summary.getComplexity()))
	msg += elapsedMsg(r.Summary.getElapsed())
	bm.setMessage(msg)
}

func passMsg(passes int) string {
	oneSpace := " "
	commaSpace := ", "
	return commaSpace + strconv.Itoa(passes) + oneSpace + "Passed"
}

func runMsg(runs int) string {
	oneSpace := " "
	return strconv.Itoa(runs) + oneSpace + "Run"
}

func elapsedMsg(elapsed GtpElapsed) string {
	oneSpace := " "
	commaSpace := ", "
	msg := commaSpace + "in" + oneSpace + strconv.FormatFloat(float64(elapsed), 'f', 3, 32) + "s"
	return msg
}

func metricsMsg(skips, fails int, coverage, complexity string) string {
	oneSpace := " "
	commaSpace := ", "
	if skips == 0 && fails == 0 && len(coverage) > 0 {
		msg := commaSpace + "Test Coverage:" + oneSpace + coverage
		msg += commaSpace + "Average Complexity:" + oneSpace + complexity
		return msg
	}
	return ""
}

func failMsg(fails int, fname, lineno string) string {
	if fails > 0 {
		oneSpace := " "
		commaSpace := ", "
		msg := commaSpace + strconv.Itoa(fails) + oneSpace + "Failed"
		msg += commaSpace + "1st in" + oneSpace + path.Base(fname)
		msg += commaSpace + "on line" + oneSpace + lineno
		return msg
	}
	return ""
}

func skipMsg(skips int) string {
	oneSpace := " "
	commaSpace := ", "
	if skips > 0 {
		return commaSpace + strconv.Itoa(skips) + oneSpace + "Skipped"
	}
	return ""
}

func (r *GtpResults) incCount(key string) {
	r.Counts[key]++
}

func (r *GtpResults) decCount(key string) {
	r.Counts[key]--
}

func (r *GtpResults) getCount(key string) int {
	return r.Counts[key]
}

type GtpSummary struct {
	Coverage   GtpCoverage
	Complexity GtpComplexity
	Elapsed    GtpElapsed
}

//"coverage: 76.7% of statements\n"}
func (s *GtpSummary) setCoverage(coverage string) {
	// Here we are expecting coverage to be a jlo.Output
	// that has already been verified to have coverage
	// and we do a little digging to extract the coverage
	// Remove trailing '\n'
	cov := strings.TrimSuffix(coverage, "\n")
	// Strip away everything but the percent coverage string ("57.8%", for example)
	cov = strings.Replace(cov, "coverage: ", "", 1)
	cov = strings.Replace(cov, " of statements", "", 1)
	s.Coverage = GtpCoverage(cov)
}

func (s *GtpSummary) setComplexity(paths []string, ignore string) {
	allStats := gocyclo.Analyze(paths, regexp.MustCompile(ignore))
	s.Complexity = GtpComplexity(fmt.Sprintf("%.3g", allStats.AverageComplexity()))
}

func (s *GtpSummary) setElapsed(elapsed GtpElapsed) {
	s.Elapsed = elapsed
}

func (s *GtpSummary) getCoverage() GtpCoverage {
	return s.Coverage
}

func (s *GtpSummary) getComplexity() GtpComplexity {
	return s.Complexity
}

func (s *GtpSummary) getElapsed() GtpElapsed {
	return s.Elapsed
}

func (e *GtpError) getMessage() string {
	return e.Message
}

func (e *GtpError) getColor() string {
	return e.Color
}

func (e *GtpErrors) Add(errorItem GtpError) {
	*e = append(*e, errorItem)
}
