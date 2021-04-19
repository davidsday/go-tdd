package main

type PgmData struct {
	Info            PD_Info            `json:"info"`
	Counts          PD_Counts          `json:"counts"`
	Firstfailedtest PD_FirstFailedTest `json:"firstfailedtest"`
	Elapsed         PD_Elapsed         `json:"elapsed"`
	Perror          PD_Perror          `json:"error"`
	Qflist          PD_QfList          `json:"qflist"`
	Barmessage      PD_BarMessage      `json:"barmessage"`
}

type PD_Info struct {
	Host        string `json:"host"`
	User        string `json:"user"`
	Begintime   string `json:"begintime"`
	Endtime     string `json:"endtime"`
	Commandline string `json:"commandline"`
}

type PD_Counts struct {
	Runs      int `json:"runs"`
	Pauses    int `json:"pauses"`
	Continues int `json:"continues"`
	Skips     int `json:"skips"`
	Passes    int `json:"passes"`
	Fails     int `json:"fails"`
	Outputs   int `json:"outputs"`
}

type PD_FirstFailedTest struct {
	Fname  string `json:"fname"`
	Tname  string `json:"tname"`
	Lineno string `json:"lineno"`
}

type PD_Elapsed float64

type PD_Perror struct {
	Validjson   bool `json:"validjson"`
	Notestfiles bool `json:"notestfiles"`
	Rcv_panic   bool `json:"panic"`
	Buildfailed bool `json:"buildfailed"`
	Msg_stderr  bool `json:"msg_stderr"`
}

type PD_QfList []PD_QfDict

type PD_BarMessage struct {
	Color   string `json:"color"`
	Message string `json:"message"`
}

type PD_QfDict struct {
	Filename string `json:"filename"`
	Lnum     int    `json:"lnum"`
	Col      int    `json:"col"`
	Vcol     int    `json:"vcol"`
	Pattern  string `json:"pattern"`
	Text     string `json:"text"`
}
