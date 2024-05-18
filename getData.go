package main

import (
	"fmt"
	"log"
	"main/api"
	"main/utils"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
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
	tickers := utils.GetTicker("data/company_tickers.json")
	uniqueTickers := make(map[string]struct{})
	iterationCounter := 0
	log.Println("Getting financial data")
	for year := startYear; year <= endYear; year++ {
		combinedData, err := api.GetCYCombinedData(tickers, year)
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
	}
	log.Println("Saved financial data:", finDataFile)
	numUniqueTickers := len(uniqueTickers)
	log.Println("Unique Tickers:", numUniqueTickers)
	uniqueTickers["VOO"] = struct{}{}
	startDate := os.Getenv("START_DATE")
	endDate := os.Getenv("END_DATE")
	intervals := []string{"monthly"} // "daily", "weekly",
	var saveFileNameStock string
	for _, interval := range intervals {
		log.Println("Getting stock data for interval:", interval)
		for symbol := range uniqueTickers {
			q, err := api.GetQuoteFromYahoo(symbol, startDate, endDate, interval)
			if err != nil {
				log.Println("Error fetching data:", err)
				return
			}
			saveFileNameStock = fmt.Sprintf(
				"data/"+"stock_price_%s.csv",
				interval,
			)
			q.WriteCSV(saveFileNameStock)
		}
		log.Println("Saved stock data:", saveFileNameStock)
	}
	log.Println("Completed!")
}
