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
	symbol := "MSFT"
	tickers := utils.GetTicker("data/company_tickers.json")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	startDate := os.Getenv("START_DATE")
	endDate := os.Getenv("END_DATE")
	interval := os.Getenv("INTERVAL")
	q, err := api.GetQuoteFromYahoo(symbol, startDate, endDate, interval)
	if err != nil {
		fmt.Println("Error fetching data:", err)
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
	cy := 2022

	combinedData := api.GetCYCombinedData(tickers, cy)
	saveFileName := "data/company_data.csv"
	utils.SaveCYCombinedData(combinedData, saveFileName)
	fmt.Println("CSV file created successfully!")
}
