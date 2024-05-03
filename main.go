package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"

// 	"main/api"
// )

// func main() {
// 	symbol := "AAPL"

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
// }

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// type CompanyData struct {
// 	Taxonomy    string         `json:"taxonomy"`
// 	Tag         string         `json:"tag"`
// 	Ccp         string         `json:"ccp"`
// 	Uom         string         `json:"uom"`
// 	Label       string         `json:"label"`
// 	Description string         `json:"description"`
// 	Pts         int            `json:"pts"`
// 	Data        []CompanyDatum `json:"data"`
// }

// type CompanyDatum struct {
// 	Accn       string `json:"accn"`
// 	CIK        int    `json:"cik"`
// 	EntityName string `json:"entityName"`
// 	Loc        string `json:"loc"`
// 	Start      string `json:"start"`
// 	End        string `json:"end"`
// 	Val        int64  `json:"val"`
// }

// type CombinedData struct {
// 	CIK         int    `json:"cik"`
// 	EntityName  string `json:"entityName"`
// 	NetCash     int64  `json:"netCash"`
// 	PropertyExp int64  `json:"propertyExp"`
// 	Outstanding int64  `json:"outstanding"`
// }

type FinancialData struct {
	Taxonomy    string `json:"taxonomy"`
	Tag         string `json:"tag"`
	Ccp         string `json:"ccp"`
	Uom         string `json:"uom"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Pts         int    `json:"pts"`
	Data        []struct {
		Accn       string `json:"accn"`
		Cik        int    `json:"cik"`
		EntityName string `json:"entityName"`
		Loc        string `json:"loc"`
		Start      string `json:"start"`
		End        string `json:"end"`
		Val        int    `json:"val"`
	} `json:"data"`
}

func main() {
	// Fetch data from APIs
	netCashData := fetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/NetCashProvidedByUsedInOperatingActivities/USD/CY2022.json")
	// propertyExpData := fetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/PaymentsToAcquirePropertyPlantAndEquipment/USD/CY2022.json")
	// outstandingData := fetchData("https://data.sec.gov/api/xbrl/frames/us-gaap/WeightedAverageNumberOfSharesOutstandingBasic/shares/CY2022.json")
	fmt.Print(netCashData)
	// Combine data
	// combinedData := make(map[int]CombinedData)
	// for _, data := range netCashData.Data {
	// 	cik := data.CIK
	// 	combinedData[cik] = CombinedData{
	// 		CIK:         cik,
	// 		EntityName:  data.EntityName,
	// 		NetCash:     data.Val,
	// 		PropertyExp: getPropertyExp(propertyExpData, cik),
	// 		Outstanding: getOutstanding(outstandingData, cik),
	// 	}
	// }

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

func fetchData(apiURL string) FinancialData {
	// Fetch data from API
	// response, err := http.Get(apiURL)
	// if err != nil {
	// 	log.Println("Error fetching data:", err)
	// 	log.Fatal()
	// 	// nil
	// }
	// defer response.Body.Close()
	// Timeout := os.Getenv("TIMEOUT")
	// TimeoutInt, err := strconv.Atoi(Timeout)
	// if err != nil {
	// 	log.Fatal("Error:", err)
	// }
	ClientTimeout := time.Duration(10) * time.Second
	client := &http.Client{
		Timeout: ClientTimeout,
	}
	initReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	initReq.Header.Set("User-Agent", "enteryourname@gmail.com")
	response, err := client.Do(initReq)
	// Read response body
	// response, err := client.Get(apiURL)
	if err != nil {
		fmt.Println("Error sending request:", err)
		// return
	}
	body, _ := io.ReadAll(response.Body)

	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return nil
	// }
	// fmt.Println("Response Body:", body)
	fmt.Println("Response Body:", string(body))
	var financialData FinancialData
	if err := json.Unmarshal(body, &financialData); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	// if err := json.NewDecoder(response.Body).Decode(&financialData); err != nil {
	// 	fmt.Println("Error decoding JSON:", err)
	// 	log.Fatal()
	// }
	// Unmarshal JSON data
	// var data CompanyData
	// json.Unmarshal(body, &data)
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// 	fmt.Println("Error unmarshalling JSON data:", err)
	// 	return nil
	// }
	for _, item := range financialData.Data {
		fmt.Printf("Accn: %s, Entity Name: %s, Val: %d\n", item.Accn, item.EntityName, item.Val)
	}
	// fmt.Println(financialData)
	return financialData
}

// func getPropertyExp(data CompanyData, cik int) int64 {
// 	for _, d := range data.Data {
// 		if d.CIK == cik {
// 			return d.Val
// 		}
// 	}
// 	return 0
// }

// func getOutstanding(data CompanyData, cik int) int64 {
// 	for _, d := range data.Data {
// 		if d.CIK == cik {
// 			return d.Val
// 		}
// 	}
// 	return 0
// }
