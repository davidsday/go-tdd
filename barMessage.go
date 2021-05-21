package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

//============================================================================
// BarMessage  - bar message
//============================================================================

type BarMessage struct {
	Color        string    `json:"color"`
	Message      string    `json:"message"`
	QuickFixList GtpQfList `json:"quickfixlist"`
}

func (b *BarMessage) setColor(color string) {
	b.Color = color
}

func (b *BarMessage) setMessage(message string) {
	b.Message = message
}

func (b *BarMessage) getColor() string {
	return b.Color
}

func (b *BarMessage) getMessage() string {
	return b.Message
}

func (b *BarMessage) marshalToStdOut() {
	// data, err := json.MarshalIndent(pgmdata, "", "    ")
	data, _ := json.Marshal(*b)
	_, err := os.Stdout.Write(data)
	chkErr(err, "Error writing to Stdout in BarMessage.marshalToStdOut()")
}

func (b *BarMessage) marshalToDisk() {
	// data, err := json.MarshalIndent(pgmdata, "", "    ")
	data, _ := json.Marshal(*b)
	err := os.WriteFile("./go-ttd_log.json", data, 0664)
	chkErr(err, "Error writing to ./go-ttd_log.json, in marshalToStdOut()")
}

func (b *BarMessage) marshalToByteString() []byte {
	// data, err := json.MarshalIndent(pgmdata, "", "    ")
	data, err := json.Marshal(*b)
	chkErr(err, "Error in marshalToByteString()")
	return data
}

//============================================================================
// GtpQfItem  - quickfixitem
//============================================================================

type GtpQfItem struct {
	Filename string `json:"filename"`
	Lnum     int    `json:"lnum"`
	Col      int    `json:"col"`
	Vcol     int    `json:"vcol"`
	Pattern  string `json:"pattern"`
	Text     string `json:"text"`
}

// func (i *GtpQfItem) getFilename() string {
//	return i.Filename
// }

// func (i *GtpQfItem) getLnum() int {
//	return i.Lnum
// }

// func (i *GtpQfItem) getCol() int {
//	return i.Col
// }

// func (i *GtpQfItem) getPattern() string {
//	return i.Pattern
// }

// func (i *GtpQfItem) getText() string {
//	return i.Text
// }

// func (i *GtpQfItem) setFilename(name string) {
//	i.Filename = name
// }

// func (i *GtpQfItem) setCol(col int) {
//	i.Col = col
// }

// func (i *GtpQfItem) setLnum(lnum int) {
//	i.Lnum = lnum
// }

// func (i *GtpQfItem) setPattern(pattern string) {
//	i.Text = pattern
// }

// func (i GtpQfItem) setText(text string) {
//	i.Text = text
// }

//============================================================================
// GtpQfList  - quickfixlist
//============================================================================

type GtpQfList []GtpQfItem

func (q *GtpQfList) Add(item GtpQfItem) {
	*q = append(*q, item)
}

func (q *GtpQfList) Count() int {
	return len(*q)
}

// Now we can build/fill the QuickFix Item
func buildQuickFixItem(args []string, parts []string, jlo JLObject) GtpQfItem {
	QfItem := GtpQfItem{}
	// QfItem.Filename = args[1] + "/" + parts[0]
	QfItem.Filename = jlo.getPackage() + "/" + parts[0]
	QfItem.Lnum, _ = strconv.Atoi(parts[1])
	QfItem.Col = 1
	QfItem.Vcol = 1
	QfItem.Pattern = jlo.getTest()
	QfItem.Text = strings.Join(parts[2:], ":")
	return QfItem
}

func newBarMessage() BarMessage {
	b := new(BarMessage)
	b.QuickFixList = GtpQfList{}
	return *b
}
