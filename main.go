package main

import (
	"main/backend/router"
	// "time"
)

func main() {
	// uniqueTickers := api.SaveFinancialData()
	// api.SaveQuoteFromYahoo(uniqueTickers)
	// gcp.UploadToGCSToBigQuery()
	// gcp.CreateMLTable()
	// gcp.CreateTrainTestTable()
	// gcp.CreateModel()
	router.NewRouter()
}
