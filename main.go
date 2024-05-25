package main

import (
	"log"

	"main/api"
	"main/gcp"
)

func main() {
	uniqueTickers := api.SaveFinancialData()
	uniqueTickers["VOO"] = struct{}{}
	api.SaveQuoteFromYahoo(uniqueTickers)
	gcp.UploadToGCSToBigQuery()
	gcp.CreateMLTable()
	log.Println("Completed!")
}
