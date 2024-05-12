package main

import (
	"fmt"
	"log"
	"main/api"
	"main/utils"
	"os"
	"strconv"

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
	for year := startYear; year <= endYear; year++ {
		combinedData := api.GetCYCombinedData(tickers, year)
		saveFileNameFinancial := "data/financial_data.csv"
		utils.SaveCYCombinedData(combinedData, saveFileNameFinancial)
		for _, ticker := range combinedData.Ticker {
			uniqueTickers[ticker] = struct{}{}
		}
	}
	numUniqueTickers := len(uniqueTickers)
	log.Println("Unique Tickers: ", numUniqueTickers)
	startDate := os.Getenv("START_DATE")
	endDate := os.Getenv("END_DATE")
	intervals := []string{"daily", "weekly", "monthly"}
	for _, interval := range intervals {
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
