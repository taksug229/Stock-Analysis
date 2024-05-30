package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"main/backend/models"
	"main/backend/utils"
)

func FetchData(apiURL string) (models.FinancialData, error) {
	utils.LoadEnv()
	Timeout := os.Getenv("TIMEOUT")
	Header := os.Getenv("HEADER")
	var financialData models.FinancialData

	TimeoutInt, err := strconv.Atoi(Timeout)
	if err != nil {
		log.Println("Error:", err)
		log.Println(apiURL)
		log.Fatal("Error:", err)
	}
	ClientTimeout := time.Duration(TimeoutInt) * time.Second
	client := &http.Client{
		Timeout: ClientTimeout,
	}
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error reading response body:", err)
		log.Println(apiURL)
		return financialData, err
	}
	request.Header.Set("User-Agent", Header)
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error sending request:", err)
		log.Println(apiURL)
		return financialData, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		log.Println(apiURL)
		return financialData, err
	}

	if err := json.Unmarshal(body, &financialData); err != nil {
		// log.Println("Error decoding JSON:", err)
		// log.Println(apiURL)
		return financialData, err
	}
	return financialData, err
}

func GetCYCombinedData(tickers []models.Ticker, cy int) (models.CombinedData, error) {
	revenueUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/Revenues/USD/CY%d.json", cy)
	revenueContractUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/RevenueFromContractWithCustomerExcludingAssessedTax/USD/CY%d.json", cy)
	salesRevenueUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/SalesRevenueNet/USD/CY%d.json", cy)
	netCashUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/NetCashProvidedByUsedInOperatingActivities/USD/CY%d.json", cy)
	propertyExpUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/PaymentsToAcquirePropertyPlantAndEquipment/USD/CY%d.json", cy)
	propertyExpAssetsUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/PaymentsToAcquireProductiveAssets/USD/CY%d.json", cy)
	sharesOutUrl := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/WeightedAverageNumberOfSharesOutstandingBasic/shares/CY%d.json", cy)

	curCashQ4Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/CashAndCashEquivalentsAtCarryingValue/USD/CY%dQ4I.json", cy)
	curCashQ3Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/CashAndCashEquivalentsAtCarryingValue/USD/CY%dQ3I.json", cy)
	curCashQ2Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/CashAndCashEquivalentsAtCarryingValue/USD/CY%dQ2I.json", cy)
	curCashQ1Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/CashAndCashEquivalentsAtCarryingValue/USD/CY%dQ1I.json", cy)
	investQ4Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/AvailableForSaleSecuritiesCurrent/USD/CY%dQ4I.json", cy)
	investQ3Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/AvailableForSaleSecuritiesCurrent/USD/CY%dQ3I.json", cy)
	investQ2Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/AvailableForSaleSecuritiesCurrent/USD/CY%dQ2I.json", cy)
	investQ1Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/AvailableForSaleSecuritiesCurrent/USD/CY%dQ1I.json", cy)
	investmentQ4Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/Investments/USD/CY%dQ4I.json", cy)
	investmentQ3Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/Investments/USD/CY%dQ3I.json", cy)
	investmentQ2Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/Investments/USD/CY%dQ2I.json", cy)
	investmentQ1Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/Investments/USD/CY%dQ1I.json", cy)
	securityQ4Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/MarketableSecuritiesCurrent/USD/CY%dQ4I.json", cy)
	securityQ3Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/MarketableSecuritiesCurrent/USD/CY%dQ3I.json", cy)
	securityQ2Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/MarketableSecuritiesCurrent/USD/CY%dQ2I.json", cy)
	securityQ1Url := fmt.Sprintf("https://data.sec.gov/api/xbrl/frames/us-gaap/MarketableSecuritiesCurrent/USD/CY%dQ1I.json", cy)

	var (
		combinedData     models.CombinedData
		revenue          float64
		revenuecontract  float64
		salesrevenue     float64
		revenuecombined  float64
		netcash          float64
		propertyexp      float64
		propertyexpasset float64
		propertycombined float64
		shares           float64
		cashasset        float64
		invest           float64
		investment       float64
		investcombined   float64
		security         float64
		startdate        string
		enddate          string
	)
	revenueData, _ := FetchData(revenueUrl)
	revenueContractData, _ := FetchData(revenueContractUrl)
	salesRevenueData, _ := FetchData(salesRevenueUrl)

	netCashData, err := FetchData(netCashUrl)
	if err != nil {
		return combinedData, err
	}
	propertyExpData, _ := FetchData(propertyExpUrl)
	propertyExpAssetData, _ := FetchData(propertyExpAssetsUrl)

	sharesOutData, err := FetchData(sharesOutUrl)
	if err != nil {
		return combinedData, err
	}

	time.Sleep(time.Second)
	curCashDataQ4, _ := FetchData(curCashQ4Url)
	curCashDataQ3, _ := FetchData(curCashQ3Url)
	curCashDataQ2, _ := FetchData(curCashQ2Url)
	curCashDataQ1, _ := FetchData(curCashQ1Url)
	investDataQ4, _ := FetchData(investQ4Url)
	investDataQ3, _ := FetchData(investQ3Url)
	investDataQ2, _ := FetchData(investQ2Url)
	investDataQ1, _ := FetchData(investQ1Url)
	time.Sleep(time.Second)
	investmentDataQ4, _ := FetchData(investmentQ4Url)
	investmentDataQ3, _ := FetchData(investmentQ3Url)
	investmentDataQ2, _ := FetchData(investmentQ2Url)
	investmentDataQ1, _ := FetchData(investmentQ1Url)
	securityDataQ4, _ := FetchData(securityQ4Url)
	securityDataQ3, _ := FetchData(securityQ3Url)
	securityDataQ2, _ := FetchData(securityQ2Url)
	securityDataQ1, _ := FetchData(securityQ1Url)

	curCashSlice := []models.FinancialData{curCashDataQ4, curCashDataQ3, curCashDataQ2, curCashDataQ1}
	investSlice := []models.FinancialData{investDataQ4, investDataQ3, investDataQ2, investDataQ1}
	investmentSlice := []models.FinancialData{investmentDataQ4, investmentDataQ3, investmentDataQ2, investmentDataQ1}
	securitySlice := []models.FinancialData{securityDataQ4, securityDataQ3, securityDataQ2, securityDataQ1}

	utils.LoadEnv()
	maxTickers := os.Getenv("MAXTICKERS")
	skipIteration := true
	maxTickersInt, err := strconv.Atoi(maxTickers)
	if err != nil {
		skipIteration = false
	}
	count := 0
	for _, data := range tickers {
		cik := data.CIK
		startdate, enddate = utils.GetCYDates(netCashData, cik)
		revenue = utils.GetFinancialKPIData(revenueData, cik)
		revenuecontract = utils.GetFinancialKPIData(revenueContractData, cik)
		salesrevenue = utils.GetFinancialKPIData(salesRevenueData, cik)
		revenuecombined = math.Max(revenue, math.Max(revenuecontract, salesrevenue))
		netcash = utils.GetFinancialKPIData(netCashData, cik)
		propertyexp = math.Abs(utils.GetFinancialKPIData(propertyExpData, cik))
		propertyexpasset = math.Abs(utils.GetFinancialKPIData(propertyExpAssetData, cik))
		propertycombined = math.Max(propertyexp, propertyexpasset)

		shares = utils.GetFinancialKPIData(sharesOutData, cik)
		// if netcash == 0 || propertyexp == 0 || shares == 0 {
		// 	continue
		// }
		cashasset = getAssetFromQuarters(curCashSlice, cik)
		invest = getAssetFromQuarters(investSlice, cik)
		// if cy <= 2022 {
		// 	invest = getAssetFromQuarters(investSlice, cik)
		// } else {
		// 	invest = 0
		// }
		investment = getAssetFromQuarters(investmentSlice, cik)
		investcombined = math.Max(invest, investment)
		security = getAssetFromQuarters(securitySlice, cik)

		// var propertyexpfloat, sharesfloat float64
		// switch v := propertyexp.(type) {
		// case int:
		// 	propertyexpfloat = math.Abs(float64(v))
		// case float64:
		// 	propertyexpfloat = math.Abs(v)
		// default:
		// 	continue
		// }
		// switch v := shares.(type) {
		// case int:
		// 	sharesfloat = float64(v)
		// case float64:
		// 	sharesfloat = v
		// default:
		// 	continue
		// }
		if shares < 1000 {
			shares = shares * 1_000_000
		}
		shares = math.Round(shares)
		combinedData.CY = append(combinedData.CY, cy)
		combinedData.StartDate = append(combinedData.StartDate, startdate)
		combinedData.EndDate = append(combinedData.EndDate, enddate)
		combinedData.Ticker = append(combinedData.Ticker, data.Ticker)
		combinedData.CIK = append(combinedData.CIK, cik)
		combinedData.EntityName = append(combinedData.EntityName, data.Title)
		combinedData.Revenue = append(combinedData.Revenue, revenuecombined)
		combinedData.NetCash = append(combinedData.NetCash, netcash)
		combinedData.PropertyExp = append(combinedData.PropertyExp, propertycombined)
		combinedData.Shares = append(combinedData.Shares, shares)
		combinedData.CashAsset = append(combinedData.CashAsset, cashasset)

		combinedData.Investments = append(combinedData.Investments, investcombined)
		combinedData.Securities = append(combinedData.Securities, security)

		count++
		if skipIteration && count >= maxTickersInt {
			break
		}
	}
	return combinedData, nil
}

func getAssetFromQuarters(assetSlice []models.FinancialData, cik int) float64 {
	// assetfloat := 0.0
	asset := 0.0
	for _, assetData := range assetSlice {
		asset = utils.GetFinancialKPIData(assetData, cik)
		// switch v := asset.(type) {
		// case int:
		// 	assetfloat = float64(v)
		// case float64:
		// 	assetfloat = v
		// default:
		// 	continue
		// }
		if asset > 0 {
			break
		}
	}
	return asset
}

func SaveFinancialData() map[string]struct{} {
	utils.LoadEnv()
	startYearStr := os.Getenv("START_YEAR")
	endYearStr := os.Getenv("END_YEAR")
	finDataFile := os.Getenv("FINANCIAL_DATA_FILE")
	startYear, err := strconv.Atoi(startYearStr)
	if err != nil {
		log.Fatal("Error:", err)
	}
	endYear, err := strconv.Atoi(endYearStr)
	if err != nil {
		log.Fatal("Error:", err)
	}
	tickers := utils.GetTicker("backend/data/company_tickers.json")
	uniqueTickers := make(map[string]struct{})
	iterationCounter := 0
	log.Println("Getting financial data")
	for year := startYear; year <= endYear; year++ {
		combinedData, err := GetCYCombinedData(tickers, year)
		iterationCounter++
		if iterationCounter%10 == 0 {
			time.Sleep(5 * time.Second)
		}
		if err != nil {
			log.Println("Failed: ", year)
			continue
		}
		utils.SaveCYCombinedData(combinedData, finDataFile)
		for _, ticker := range combinedData.Ticker {
			uniqueTickers[ticker] = struct{}{}
		}
		log.Println("Success: ", year)
		time.Sleep(time.Second)
	}
	log.Println("Saved financial data:", finDataFile)
	numUniqueTickers := len(uniqueTickers)
	log.Println("Unique Tickers:", numUniqueTickers)
	return uniqueTickers
}

func GetQuoteFromYahoo(symbol, startDate, endDate, period string) (models.Quote, error) {
	var resp *http.Response
	utils.LoadEnv()
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

func SaveQuoteFromYahoo(uniqueTickers map[string]struct{}) {
	utils.LoadEnv()
	startDate := os.Getenv("START_DATE")
	endDate := os.Getenv("END_DATE")
	intervals := strings.Split(os.Getenv("INTERVALS"), ",")
	var saveFileNameStock string
	for _, interval := range intervals {
		log.Println("Getting stock data for interval:", interval)
		count := 0
		for symbol := range uniqueTickers {
			q, err := GetQuoteFromYahoo(symbol, startDate, endDate, interval)
			if err != nil {
				log.Println("Error fetching data:", err)
				return
			}
			saveFileNameStock = fmt.Sprintf(
				"backend/data/"+"stock_price_%s.csv",
				interval,
			)
			q.WriteCSV(saveFileNameStock)

			count++
			if count%100 == 0 {
				log.Println("Ticker count:", count)
				time.Sleep(5 * time.Second)
			}
		}
		log.Println("Saved stock data:", saveFileNameStock)
	}
}
