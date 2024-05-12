package models

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

type Ticker struct {
	CIK    int    `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

type Quote struct {
	Symbol string      `json:"symbol"`
	Date   []time.Time `json:"date"`
	Open   []float64   `json:"open"`
	High   []float64   `json:"high"`
	Low    []float64   `json:"low"`
	Close  []float64   `json:"close"`
	Volume []float64   `json:"volume"`
}

func NewQuote(symbol string, bars int) Quote {
	return Quote{
		Symbol: symbol,
		Date:   make([]time.Time, bars),
		Open:   make([]float64, bars),
		High:   make([]float64, bars),
		Low:    make([]float64, bars),
		Close:  make([]float64, bars),
		Volume: make([]float64, bars),
	}
}

func (q Quote) WriteCSV(filename string) error {
	if filename == "" {
		if q.Symbol != "" {
			filename = q.Symbol + ".csv"
		} else {
			filename = "quote.csv"
		}
	}
	csv := q.CSV()
	return os.WriteFile(filename, []byte(csv), 0644)
}

func (q Quote) CSV() string {
	precision := getPrecision(q.Symbol)
	var buffer bytes.Buffer
	buffer.WriteString("datetime,open,high,low,close,volume\n")
	for bar := range q.Close {
		str := fmt.Sprintf("%s,%.*f,%.*f,%.*f,%.*f,%.*f\n", q.Date[bar].Format("2006-01-02 15:04"),
			precision, q.Open[bar], precision, q.High[bar], precision, q.Low[bar], precision, q.Close[bar], precision, q.Volume[bar])
		buffer.WriteString(str)
	}
	return buffer.String()
}

func getPrecision(symbol string) int {
	var precision int
	precision = 2
	if strings.Contains(strings.ToUpper(symbol), "BTC") ||
		strings.Contains(strings.ToUpper(symbol), "ETH") ||
		strings.Contains(strings.ToUpper(symbol), "USD") {
		precision = 8
	}
	return precision
}

type FinancialData struct {
	Taxonomy    string `json:"taxonomy"`
	Tag         string `json:"tag"`
	Ccp         string `json:"ccp"`
	Uom         string `json:"uom"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Pts         int64  `json:"pts"`
	Data        []struct {
		Accn       string  `json:"accn"`
		CIK        int     `json:"cik"`
		EntityName string  `json:"entityName"`
		Loc        string  `json:"loc"`
		Start      string  `json:"start"`
		End        string  `json:"end"`
		Val        float64 `json:"val"`
	} `json:"data"`
}

type CombinedData struct {
	CY          []int
	StartDate   []string
	EndDate     []string
	Ticker      []string
	CIK         []int
	EntityName  []string
	NetCash     []interface{}
	PropertyExp []interface{}
	Shares      []interface{}
}
