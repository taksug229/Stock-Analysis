package main

import (
	"fmt"
	// "reflect"

	"github.com/markcheno/go-quote"
)

func main() {
	// Replace "AAPL" with the stock symbol you want to fetch data for
	symbol := "AAPL"

	// Fetch historical data for the specified stock symbol
	q, err := quote.NewQuoteFromYahoo(symbol, "2024-04-01", "2024-04-08", quote.Daily, true)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	fmt.Print(q.CSV())
	q.WriteCSV("")
	// v := reflect.ValueOf(q)

	// Print the historical data
	// for i := 0; i < v.NumField(); i++ {
	// 	// Get the field value
	// 	fieldValue := v.Field(i)

	// 	// Get the field type
	// 	fieldType := v.Type().Field(i)
	// 	fmt.Printf("Field Name: %s, Type: %s, Value: %v\n",
	// 		fieldType.Name, fieldType.Type, fieldValue.Interface())
	// }
}
