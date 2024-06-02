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

	_, file_err := os.Stat(filename)
	if file_err != nil && !os.IsNotExist(file_err) {
		return file_err
	}

	// Open file with append mode, create if not exists
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header if the file doesn't exist
	if os.IsNotExist(file_err) {
		header := "ticker,datetime,open,high,low,close,volume\n"
		_, err := file.WriteString(header)
		if err != nil {
			return err
		}
	}
	csv := q.CSV()
	_, err = file.WriteString(csv)
	if err != nil {
		return err
	}

	return nil
}

func (q Quote) CSV() string {
	precision := getPrecision(q.Symbol)
	var buffer bytes.Buffer
	for bar := range q.Close {
		str := fmt.Sprintf("%s,%s,%.*f,%.*f,%.*f,%.*f,%.*f\n", q.Symbol, q.Date[bar].Format("2006-01-02 15:04"),
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
	Revenue     []float64
	NetCash     []float64
	PropertyExp []float64
	Shares      []float64
	CashAsset   []float64
	Investments []float64
	Securities  []float64
}

type LiveStockData struct {
	Ticker              string
	CurrentStockPrice   string
	PredictedStockPrice string
	MarketCap           string
	IntrinsicValue      string
	Recommendation      string
}

type AvailableTicker struct {
	ID string
}
