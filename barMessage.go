package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	// data, err := json.MarshalIndent(*b, "", "    ")
	data, _ := json.Marshal(*b)
	_, err := os.Stdout.Write(data)
	chkErr(err, "Error writing to Stdout in BarMessage.marshalToStdOut()")
}

func (b *BarMessage) marshalToDisk() {
	data, _ := json.Marshal(*b)
	path := filepath.Join(".", "go-tdd_log.json")
	err := os.WriteFile(path, data, 0664)
	emsg := fmt.Sprintf("Error writing to '%s', in marshalToDisk()", path)
	chkErr(err, emsg)
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

// searchDir
// filename
// linenum
// pattern
// text

// Now we can build/fill the QuickFix Item
func buildQuickFixItem(searchDir, filename, linenum, pattern, text string) GtpQfItem {
	QfItem := GtpQfItem{}
	if searchDir != "" {
		if !strings.HasSuffix(searchDir, string(filepath.Separator)) {
			QfItem.Filename = searchDir + string(filepath.Separator)
		}
	}
	QfItem.Filename += filename
	QfItem.Lnum, _ = strconv.Atoi(linenum)
	QfItem.Col = 1
	QfItem.Vcol = 1
	QfItem.Pattern = pattern
	QfItem.Text = text
	return QfItem
}

func newBarMessage() BarMessage {
	b := new(BarMessage)
	b.QuickFixList = GtpQfList{}
	return *b
}
