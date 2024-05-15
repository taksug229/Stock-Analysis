package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"main/models"
	"main/utils"
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
	if period == "daily" {
		interval = "1d"
	} else if period == "weekly" {
		interval = "1wk"
	} else if period == "monthly" {
		interval = "1mo"
	} else {
		log.Fatal("period must be either 'daily', 'weekly', or 'monthly'")
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

func FetchData(apiURL string) (models.FinancialData, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	Timeout := os.Getenv("TIMEOUT")
	Header := os.Getenv("HEADER")
	var financialData models.FinancialData

	TimeoutInt, err := strconv.Atoi(Timeout)
	if err != nil {
		log.Fatal("Error:", err)
	}
	ClientTimeout := time.Duration(TimeoutInt) * time.Second
	client := &http.Client{
		Timeout: ClientTimeout,
	}
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error reading response body:", err)
		return financialData, err
	}
	request.Header.Set("User-Agent", Header)
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return financialData, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return financialData, err
	}

	if err := json.Unmarshal(body, &financialData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return financialData, err
	}
	return financialData, err
}

func GetCYCombinedData(tickers []models.Ticker, cy int) (models.CombinedData, error) {
	netCashUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/NetCashProvidedByUsedInOperatingActivities/USD/CY%d.json", cy)
	propertyExpUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/PaymentsToAcquirePropertyPlantAndEquipment/USD/CY%d.json", cy)
	sharesOutUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/WeightedAverageNumberOfSharesOutstandingBasic/shares/CY%d.json", cy)
	var (
		combinedData models.CombinedData
		netcash      interface{}
		propertyexp  interface{}
		shares       interface{}
		startdate    string
		enddate      string
	)
	netCashData, err := FetchData(netCashUrl)
	if err != nil {
		return combinedData, err
	}
	propertyExpData, err := FetchData(propertyExpUrl)
	if err != nil {
		return combinedData, err
	}
	sharesOutData, err := FetchData(sharesOutUrl)
	if err != nil {
		return combinedData, err
	}

	count := 0
	for _, data := range tickers {
		cik := data.CIK
		startdate, enddate = utils.GetCYDates(netCashData, cik)
		netcash = utils.GetFinancialData(netCashData, cik)
		propertyexp = utils.GetFinancialData(propertyExpData, cik)
		shares = utils.GetFinancialData(sharesOutData, cik)
		var sharesfloat float64
		switch v := shares.(type) {
		case int:
			sharesfloat = float64(v)
		case float64:
			sharesfloat = v
		default:
			// log.Printf("Skipping %s due to oustanding shares issue", data.Ticker)
			continue
		}
		if sharesfloat < 1000 {
			sharesfloat = sharesfloat * 1_000_000
		}
		if netcash == 0 || propertyexp == 0 || shares == 0 {
			// log.Printf("Skipping %s due to no netcash, protertyexp, or shares outstanding.", data.Ticker)
			continue
		}
		combinedData.CY = append(combinedData.CY, cy)
		combinedData.StartDate = append(combinedData.StartDate, startdate)
		combinedData.EndDate = append(combinedData.EndDate, enddate)
		combinedData.Ticker = append(combinedData.Ticker, data.Ticker)
		combinedData.CIK = append(combinedData.CIK, cik)
		combinedData.EntityName = append(combinedData.EntityName, data.Title)
		combinedData.NetCash = append(combinedData.NetCash, netcash)
		combinedData.PropertyExp = append(combinedData.PropertyExp, propertyexp)
		combinedData.Shares = append(combinedData.Shares, sharesfloat)

		// TODO Remove during final
		count++
		if count >= 5 {
			break
		}
	}
	return combinedData, nil
}
