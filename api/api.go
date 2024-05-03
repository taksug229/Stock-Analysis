package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"main/models"
)

func GetQuoteFromYahoo(symbol, startDate, endDate, period string) (models.Quote, error) {
	var resp *http.Response
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	from := ParseDateString(startDate)
	to := ParseDateString(endDate)

	Timeout := os.Getenv("TIMEOUT")
	TimeoutInt, err := strconv.Atoi(Timeout)
	if err != nil {
		log.Fatal("Error:", err)
	}
	ClientTimeout := time.Duration(TimeoutInt) * time.Second
	client := &http.Client{
		Timeout: ClientTimeout,
	}
	initReq, err := http.NewRequest("GET", "https://finance.yahoo.com", nil)
	if err != nil {
		return models.NewQuote("", 0), err
	}
	initReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; U; Linux i686) Gecko/20071127 Firefox/2.0.0.11")
	client.Do(initReq)
	var interval string
	if period == "Daily" {
		interval = "1d"
	} else if period == "Weekly" {
		interval = "1wk"
	} else if period == "Monthly" {
		interval = "1mo"
	} else {
		log.Fatal("period must be either 'Daily', 'Weekly', or 'Monthly'")
	}
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=%s&events=history&corsDomain=finance.yahoo.com",
		symbol,
		from.Unix(),
		to.Unix(),
		interval,
	)
	resp, err = client.Get(url)
	if err != nil {
		log.Printf("symbol '%s' not found between '%s' and '%s'\n", symbol, startDate, endDate)
		return models.NewQuote("", 0), err
	}
	defer resp.Body.Close()
	var csvdata [][]string
	reader := csv.NewReader(resp.Body)
	csvdata, err = reader.ReadAll()
	if err != nil {
		log.Printf("bad data for symbol '%s'\n", symbol)
		return models.NewQuote("", 0), err
	}

	// numrows := len(csvdata) - 1
	// quote := models.NewQuote(symbol, numrows)
	var dateList []time.Time
	var openList []float64
	var highList []float64
	var lowList []float64
	var closeList []float64
	var volumeList []float64
	for row := 1; row < len(csvdata); row++ {
		// Parse row of data
		d, _ := time.Parse("2006-01-02", csvdata[row][0])
		o, _ := strconv.ParseFloat(csvdata[row][1], 64)
		h, _ := strconv.ParseFloat(csvdata[row][2], 64)
		l, _ := strconv.ParseFloat(csvdata[row][3], 64)
		c, _ := strconv.ParseFloat(csvdata[row][4], 64)
		a, _ := strconv.ParseFloat(csvdata[row][5], 64)
		v, _ := strconv.ParseFloat(csvdata[row][6], 64)
		if o == 0 && h == 0 && l == 0 && c == 0 && a == 0 && v == 0 {
			continue
		}
		dateList = append(dateList, d)
		openList = append(openList, o)
		highList = append(highList, h)
		lowList = append(lowList, l)
		closeList = append(closeList, c)
		volumeList = append(volumeList, v)
	}
	quote := models.Quote{
		Symbol: symbol,
		Date:   dateList,
		Open:   openList,
		High:   highList,
		Low:    lowList,
		Close:  closeList,
		Volume: volumeList,
	}
	return quote, nil
}

func ParseDateString(dt string) time.Time {
	if dt == "" {
		return time.Now()
	}
	t, _ := time.Parse("2006-01-02 15:04", dt+"0000-01-01 00:00"[len(dt):])
	return t
}
