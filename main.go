package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// "main/api"
	"main/gcp"

	"github.com/gorilla/mux"
)

func main() {
	// uniqueTickers := api.SaveFinancialData()
	// uniqueTickers["VOO"] = struct{}{}
	// api.SaveQuoteFromYahoo(uniqueTickers)
	// gcp.UploadToGCSToBigQuery()
	// gcp.CreateMLTable()

	router := mux.NewRouter()
	router.HandleFunc("/goodstocks", GetGoodStocksHandler).Methods("GET")
	// router.HandleFunc("/resources/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "/resources/%s!", mux.Vars(r)["id"])
	// })
	// router.HandleFunc("/resources/{id:[0-9]+}/values", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "/resources/values: %s", mux.Vars(r)["id"])
	// })
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	log.Println("Completed!")
}

func GetGoodStocksHandler(w http.ResponseWriter, r *http.Request) {
	// paramStr := strings.Split(r.URL.Path, "/")
	sqlFile := "sql/get_good_stocks.sql"
	err := gcp.PrintQueryResults(sqlFile, w)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
}
