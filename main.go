package main

import (
	"fmt"
	"log"
	"main/api"
	"main/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cy := 2022
	tickers := utils.GetTicker("data/company_tickers.json")
	combinedData := api.GetCYCombinedData(tickers, cy)
	saveFileNameFinancial := "data/financial_data.csv"
	utils.SaveCYCombinedData(combinedData, saveFileNameFinancial)

	symbol := "MSFT"
	startDate := os.Getenv("START_DATE")
	endDate := os.Getenv("END_DATE")
	interval := os.Getenv("INTERVAL")
	q, err := api.GetQuoteFromYahoo(symbol, startDate, endDate, interval)
	if err != nil {
		log.Println("Error fetching data:", err)
		return
	}
	filename := fmt.Sprintf(
		"data/"+"%s-%s-%s-%s.csv",
		symbol,
		startDate,
		endDate,
		interval,
	)
	q.WriteCSV(filename)
	fmt.Println("CSV file created successfully!")
}
