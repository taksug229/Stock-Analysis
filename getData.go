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
		if err != nil {
			log.Println("Failed: ", year)
			continue
		}
		saveFileNameFinancial := "data/financial_data.csv"
		utils.SaveCYCombinedData(combinedData, saveFileNameFinancial)
		for _, ticker := range combinedData.Ticker {
			uniqueTickers[ticker] = struct{}{}
		}
		log.Println("Success: ", year)
		iterationCounter++
		if iterationCounter%10 == 0 {
			time.Sleep(5 * time.Second)
		}
	}
	numUniqueTickers := len(uniqueTickers)
	log.Println("Unique Tickers: ", numUniqueTickers)
	startDate := os.Getenv("START_DATE")
	endDate := os.Getenv("END_DATE")
	intervals := []string{"monthly"} // "daily", "weekly",
	for _, interval := range intervals {
		log.Println("Getting stock data for ", interval)
		for symbol := range uniqueTickers {
			q, err := api.GetQuoteFromYahoo(symbol, startDate, endDate, interval)
			if err != nil {
				log.Println("Error fetching data:", err)
				return
			}
			saveFileNameStock := fmt.Sprintf(
				"data/"+"stock_price-%s-%s-%s.csv",
				interval,
				startDate,
				endDate,
			)
			q.WriteCSV(saveFileNameStock)
		}
	}

	log.Println("CSV file created successfully!")
}
