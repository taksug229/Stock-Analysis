package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"main/api"
	"main/models"
	"main/utils"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// symbol := "AAPL"
	tickers := utils.GetTicker("data/company_tickers.json")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// startDate := os.Getenv("START_DATE")
	// endDate := os.Getenv("END_DATE")
	// interval := os.Getenv("INTERVAL")
	// q, err := api.GetQuoteFromYahoo(symbol, startDate, endDate, interval)
	// if err != nil {
	// 	fmt.Println("Error fetching data:", err)
	// 	return
	// }

	// filename := fmt.Sprintf(
	// 	"data/"+"%s-%s-%s-%s.csv",
	// 	symbol,
	// 	startDate,
	// 	endDate,
	// 	interval,
	// )
	// q.WriteCSV(filename)
	cy := 2022
	netCashData := api.FetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/NetCashProvidedByUsedInOperatingActivities/USD/CY2022.json")
	propertyExpData := api.FetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/PaymentsToAcquirePropertyPlantAndEquipment/USD/CY2022.json")
	sharesOutData := api.FetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/WeightedAverageNumberOfSharesOutstandingBasic/shares/CY2022.json")
	// Combine data
	var (
		combinedData models.CombinedData
		netcash      interface{}
		propertyexp  interface{}
		shares       interface{}
		startdate    string
		enddate      string
	)
	count := 0
	for _, data := range tickers {

		cik := data.CIK
		startdate, enddate = utils.GetCYDates(netCashData, cik)
		netcash = utils.GetFinancialData(netCashData, cik)
		propertyexp = utils.GetFinancialData(propertyExpData, cik)
		shares = utils.GetFinancialData(sharesOutData, cik)
		var sharesfloat float64
		switch v := shares.(type) {
		case int:
			sharesfloat = float64(v)
		case float64:
			sharesfloat = v
		default:
			// log.Printf("Skipping %s due to oustanding shares issue", data.Ticker)
			continue
		}
		if sharesfloat < 1000 {
			sharesfloat = sharesfloat * 1_000_000
		}
		if netcash == 0 || propertyexp == 0 || shares == 0 {
			// log.Printf("Skipping %s due to no netcash, protertyexp, or shares outstanding.", data.Ticker)
			continue
		}
		combinedData.CY = append(combinedData.CY, cy)
		combinedData.StartDate = append(combinedData.StartDate, startdate)
		combinedData.EndDate = append(combinedData.EndDate, enddate)
		combinedData.Ticker = append(combinedData.Ticker, data.Ticker)
		combinedData.CIK = append(combinedData.CIK, cik)
		combinedData.EntityName = append(combinedData.EntityName, data.Title)
		combinedData.NetCash = append(combinedData.NetCash, netcash)
		combinedData.PropertyExp = append(combinedData.PropertyExp, propertyexp)
		combinedData.Shares = append(combinedData.Shares, sharesfloat)

		count++ // TODO Remove during final
		if count >= 5 {
			break
		}
	}

	// Write to CSV
	file, err := os.Create("data/company_data.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
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
			utils.FormatInterface(combinedData.NetCash[i]),
			utils.FormatInterface(combinedData.PropertyExp[i]),
			utils.FormatInterface(combinedData.Shares[i]),
		}
		writer.Write(row)
	}

	fmt.Println("CSV file created successfully!")
}
