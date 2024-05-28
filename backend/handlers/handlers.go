package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"main/backend/gcp"
	"main/backend/models"

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
	intrinsic_val, shares := gcp.GetTickerFinancials(ticker)

	data := models.LiveStockData{
		Ticker:              ticker,
		CurrentStockPrice:   float64(shares),
		PredictedStockPrice: 1600.00,
		MarketValue:         1_000_000_000.00,
		IntrinsicValue:      intrinsic_val,
		Recommendation:      "Buy",
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
