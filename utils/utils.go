package utils

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"main/models"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetTicker(filepath string) []models.Ticker {
	var tickers map[string]models.Ticker
	var tickers_ordered []models.Ticker
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Error opening JSON file:", err)
		return tickers_ordered
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&tickers)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return tickers_ordered
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	topnstr := os.Getenv("TOPTICKERS")
	topn, err := strconv.Atoi(topnstr)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < topn; i++ {
		s := strconv.Itoa(i)
		tickers_ordered = append(tickers_ordered, tickers[s])
	}
	return tickers_ordered
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

func FormatInterface(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	// Add cases for other types if needed
	default:
		return ""
	}
}

func SaveCYCombinedData(combinedData models.CombinedData, saveFileName string) {
	// Write to CSV
	file, err := os.Create(saveFileName)
	if err != nil {
		log.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"CY", "StartDate", "EndDate", "Ticker", "CIK", "EntityName", "NetCash", "PropertyExp", "Shares"})
	for i := 0; i < len(combinedData.CY); i++ {
		row := []string{
			strconv.Itoa(combinedData.CY[i]),
			combinedData.StartDate[i],
			combinedData.EndDate[i],
			combinedData.Ticker[i],
			strconv.Itoa(combinedData.CIK[i]),
			combinedData.EntityName[i],
			FormatInterface(combinedData.NetCash[i]),
			FormatInterface(combinedData.PropertyExp[i]),
			FormatInterface(combinedData.Shares[i]),
		}
		writer.Write(row)
	}
}
