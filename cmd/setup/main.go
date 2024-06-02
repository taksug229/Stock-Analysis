// cmd/setup/main.go
package main

import (
	"log"
	"main/backend/api"
	"main/backend/gcp"
)

func main() {
	log.Println("Starting environment setup...")

	uniqueTickers := api.SaveFinancialData()
	api.SaveQuoteFromYahoo(uniqueTickers)
	gcp.UploadToGCSToBigQuery()
	gcp.CreateMLTable()
	gcp.CreateTrainTestTable()
	gcp.CreateModel()
	log.Println("Environment setup complete!")
}
