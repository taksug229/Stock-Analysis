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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sqlFile := "backend/sql/get_available_tickers.sql"
	tmpl := template.Must(template.ParseFiles("frontend/templates/index.html"))
	availableTickers, _ := gcp.GetAvailableTickers(sqlFile)
	err := tmpl.Execute(w, availableTickers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetGoodStocksHandler(w http.ResponseWriter, r *http.Request) {
	sqlFile := "backend/sql/get_good_stocks.sql"
	tmpl := template.Must(template.ParseFiles("frontend/templates/index.html"))
	availableTickers, _ := gcp.GetAvailableTickers(sqlFile)
	err := tmpl.Execute(w, availableTickers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	tmpl, err := template.ParseFiles("frontend/templates/ticker_info.html")
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
