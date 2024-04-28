package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"main/api"
)

func main() {
	symbol := "AAPL"

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
		"%s-%s-%s-%s.csv",
		symbol,
		startDate,
		endDate,
		interval,
	)
	q.WriteCSV(filename)
}
