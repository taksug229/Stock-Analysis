package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"main/backend/api"
	"main/backend/gcp"
	"main/backend/models"
	"main/backend/utils"

	"github.com/gorilla/mux"
)

func GetGoodStocksHandler(w http.ResponseWriter, r *http.Request) {
	// paramStr := strings.Split(r.URL.Path, "/")
	sqlFile := "backend/sql/get_good_stocks.sql"
	err := gcp.PrintQueryResults(sqlFile, w)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
}

func GetLiveStockData(w http.ResponseWriter, r *http.Request) {
	ticker := mux.Vars(r)["id"]
	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	todayFormatted := today.Format("2006-01-02")
	yesterdayFormatted := yesterday.Format("2006-01-02")

	q, _ := api.GetQuoteFromYahoo(ticker, yesterdayFormatted, todayFormatted, "daily")
	stockprice := q.Close[0]
	stockpricerounded := fmt.Sprintf("%.2f", stockprice)
	intrinsicval, marketcap, predictedStockPrice := gcp.GetStockInfo(ticker, stockpricerounded)

	predictedStockrounded := fmt.Sprintf("%.2f", predictedStockPrice)
	var rec string
	if ((intrinsicval * 0.7) > marketcap) && (predictedStockPrice > stockprice) {
		rec = "Buy!"
	} else {
		rec = "Don't Buy"
	}

	marketcapShorten := utils.ShortenLargeNumbers(marketcap)
	intrinsicvalShorten := utils.ShortenLargeNumbers(intrinsicval)

	data := models.LiveStockData{
		Ticker:              ticker,
		CurrentStockPrice:   stockpricerounded,
		PredictedStockPrice: predictedStockrounded,
		MarketCap:           marketcapShorten,
		IntrinsicValue:      intrinsicvalShorten,
		Recommendation:      rec,
	}
	// Parse the template file
	tmpl, err := template.ParseFiles("frontend/templates/template.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("frontend/templates/index.html"))
	availableTickers, _ := gcp.GetAvailableTickers()
	// var tickers = []models.AvailableTicker{
	// 	{ID: "AAPL"},
	// 	{ID: "GOOGL"},
	// 	{ID: "AMZN"},
	// }
	err := tmpl.Execute(w, availableTickers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
