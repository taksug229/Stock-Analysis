package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"main/backend/models"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

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

func GetFinancialKPIData(data models.FinancialData, cik int) float64 {
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
	default:
		return ""
	}
}

func MaxOfFloats(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}

	maxValue := values[0]
	for _, value := range values[1:] {
		maxValue = math.Max(maxValue, value)
	}
	return maxValue
}

func SaveCYCombinedData(combinedData models.CombinedData, saveFileName string) {
	_, file_err := os.Stat(saveFileName)
	if file_err != nil && !os.IsNotExist(file_err) {
		return
	}

	file, err := os.OpenFile(saveFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if os.IsNotExist(file_err) {
		writer.Write([]string{"CY", "StartDate", "EndDate", "Ticker", "CIK", "EntityName", "Revenue", "NetCash", "PropertyExp", "Shares", "CashAsset", "Investments", "Securities"})
	}
	for i := 0; i < len(combinedData.CY); i++ {
		row := []string{
			strconv.Itoa(combinedData.CY[i]),
			combinedData.StartDate[i],
			combinedData.EndDate[i],
			combinedData.Ticker[i],
			strconv.Itoa(combinedData.CIK[i]),
			combinedData.EntityName[i],
			FormatInterface(combinedData.Revenue[i]),
			FormatInterface(combinedData.NetCash[i]),
			FormatInterface(combinedData.PropertyExp[i]),
			FormatInterface(combinedData.Shares[i]),
			FormatInterface(combinedData.CashAsset[i]),
			FormatInterface(combinedData.Investments[i]),
			FormatInterface(combinedData.Securities[i]),
		}
		writer.Write(row)
	}
}

// AuthMiddleware checks for authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication logic
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func ShortenLargeNumbers(value float64) string {
	var (
		valueShort  string
		valueString string
	)
	if value >= 1_000_000_000_000 {
		valueShort = fmt.Sprintf("%.2f", value/1_000_000_000_000)
		valueString = valueShort + "T"
	} else if value >= 1_000_000_000 {
		valueShort = fmt.Sprintf("%.2f", value/1_000_000_000)
		valueString = valueShort + "B"
	} else {
		valueShort = fmt.Sprintf("%.2f", value/1_000_000)
		valueString = valueShort + "M"
	}
	return valueString
}
