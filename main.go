package main

import (
	"fmt"
	"log"
	"main/api"
	"main/models"
	"main/utils"
)

func main() {
	// symbol := "AAPL"
	tickers := utils.GetTicker("data/company_tickers.json")
	// 	err := godotenv.Load()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	startDate := os.Getenv("START_DATE")
	// 	endDate := os.Getenv("END_DATE")
	// 	interval := os.Getenv("INTERVAL")
	// 	q, err := api.GetQuoteFromYahoo(symbol, startDate, endDate, interval)
	// 	if err != nil {
	// 		fmt.Println("Error fetching data:", err)
	// 		return
	// 	}

	// 	filename := fmt.Sprintf(
	// 		"data/"+"%s-%s-%s-%s.csv",
	// 		symbol,
	// 		startDate,
	// 		endDate,
	// 		interval,
	// 	)
	// 	q.WriteCSV(filename)
	cy := 2022
	netCashData := api.FetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/NetCashProvidedByUsedInOperatingActivities/USD/CY2022.json")
	propertyExpData := api.FetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/PaymentsToAcquirePropertyPlantAndEquipment/USD/CY2022.json")
	outstandingData := api.FetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/WeightedAverageNumberOfSharesOutstandingBasic/shares/CY2022.json")
	// Combine data
	var (
		combinedData models.CombinedData
		netcash      interface{}
		propertyexp  interface{}
		outstanding  interface{}
		startdate    string
		enddate      string
	)
	count := 0
	for _, data := range tickers {
		cik := data.CIK
		startdate, enddate = utils.GetCYDates(netCashData, cik)
		netcash = utils.GetFinancialData(netCashData, cik)
		propertyexp = utils.GetFinancialData(propertyExpData, cik)
		outstanding = utils.GetFinancialData(outstandingData, cik)
		var outstandingfloat float64
		switch v := outstanding.(type) {
		case int:
			outstandingfloat = float64(v)
		case float64:
			outstandingfloat = v
		default:
			log.Printf("Skipping %s due to oustanding shares issue", data.Ticker)
			continue
		}
		if outstandingfloat < 1000 {
			outstandingfloat = outstandingfloat * 1_000_000
		}
		if netcash == 0 || propertyexp == 0 || outstanding == 0 {
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
		combinedData.Outstanding = append(combinedData.Outstanding, outstandingfloat)

		count++
		if count >= 5 {
			break
		}
	}
	fmt.Println(combinedData)

	// // Write to CSV
	// file, err := os.Create("data/company_data.csv")
	// if err != nil {
	// 	fmt.Println("Error creating CSV file:", err)
	// 	return
	// }
	// defer file.Close()

	// writer := csv.NewWriter(file)
	// defer writer.Flush()

	// // Write header
	// writer.Write([]string{"CIK", "EntityName", "NetCash", "PropertyExp", "Outstanding"})
	// for _, item := range netCashData.Data {
	// 	fmt.Printf("Accn: %s, Entity Name: %s, Val: %d\n", item.Accn, item.EntityName, item.Val)
	// }
	// // Write data
	// for _, data := range combinedData {
	// 	writer.Write([]string{
	// 		fmt.Sprintf("%d", data.CIK),
	// 		data.EntityName,
	// 		fmt.Sprintf("%d", data.NetCash),
	// 		fmt.Sprintf("%d", data.PropertyExp),
	// 		fmt.Sprintf("%d", data.Outstanding),
	// 	})
	// }

	// fmt.Println("CSV file created successfully!")
}
