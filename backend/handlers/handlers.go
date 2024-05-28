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
	intrinsicval, shares := gcp.GetTickerFinancials(ticker)
	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	todayFormatted := today.Format("2006-01-02")
	yesterdayFormatted := yesterday.Format("2006-01-02")

	q, _ := api.GetQuoteFromYahoo(ticker, yesterdayFormatted, todayFormatted, "daily")
	stockprice := q.Close[0]
	stockpricerounded := fmt.Sprintf("%.2f", stockprice)
	marketcap := float64(float64(shares) * stockprice)
	var rec string
	if intrinsicval > marketcap {
		rec = "Buy!"
	} else {
		rec = "Don't Buy"
	}

	marketcapShorten := utils.ShortenLargeNumbers(marketcap)
	intrinsicvalShorten := utils.ShortenLargeNumbers(intrinsicval)

	data := models.LiveStockData{
		Ticker:              ticker,
		CurrentStockPrice:   stockpricerounded,
		PredictedStockPrice: "1600.00",
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

// User Handlers
// func GetUsers(w http.ResponseWriter, r *http.Request) {
//     users := []string{"user1", "user2"}
//     json.NewEncoder(w).Encode(users)
// }

// func GetUser(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]
//     user := map[string]string{"id": id, "name": "John Doe"}
//     json.NewEncoder(w).Encode(user)
// }

// // Product Handlers
// func GetProducts(w http.ResponseWriter, r *http.Request) {
//     products := []string{"product1", "product2"}
//     json.NewEncoder(w).Encode(products)
// }

// func GetProduct(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]
//     product := map[string]string{"id": id, "name": "Product A"}
//     json.NewEncoder(w).Encode(product)
// }
