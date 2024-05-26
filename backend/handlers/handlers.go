package handlers

import (
	"fmt"
	"net/http"

	"main/backend/gcp"
)

func GetGoodStocksHandler(w http.ResponseWriter, r *http.Request) {
	// paramStr := strings.Split(r.URL.Path, "/")
	sqlFile := "backend/sql/get_good_stocks.sql"
	err := gcp.PrintQueryResults(sqlFile, w)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
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
