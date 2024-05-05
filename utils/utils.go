package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models"
	"os"
)

func GetTicker(filepath string) map[string]models.Ticker {
	var tickers map[string]models.Ticker
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Error opening JSON file:", err)
		return tickers
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&tickers)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return tickers
	}
	return tickers
}

func GetFinancialData(data models.FinancialData, cik int) interface{} {
	for _, d := range data.Data {
		if d.CIK == cik {
			return d.Val
		}
	}
	return 0
}

func GetCYDates(data models.FinancialData, cik int) (string, string) {
	for _, d := range data.Data {
		if d.CIK == cik {
			return d.Start, d.End
		}
	}
	return "", ""
}
