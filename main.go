package main

import (
	"database/sql"
	"encoding/csv"
	
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	
)

// type Stock struct {
//     Code  string
//     Name  string
//     Open  float64
//     High  float64
//     Low   float64
//     Close float64
// }


// func main() {
//     // Connect to the MySQL database
//     // db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/your_database")
//     db, err := sql.Open("mysql", "ODBC:0@tcp(localhost:3306)/app_database")
//     if err != nil {
//         panic(err)
//     }
//     defer db.Close()

//     // Read and parse the CSV file
//     stocks, err := readCSV("EQ190124.CSV")
//     if err != nil {
//         panic(err)
//     }

//     // Insert data into the database
//     insertStocks(db, stocks)

//     db, err = sql.Open("mysql", "username:password@/your_database")
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer db.Close()

//     router := mux.NewRouter()
//     router.HandleFunc("/stocks/top", getTopStocks).Methods("GET")
//     router.HandleFunc("/stocks/search", findStocksByName).Methods("GET")
//     // router.HandleFunc("/stocks/history/{code}", getStockHistory).Methods("GET")
//     router.HandleFunc("/favourites", addFavouriteStock).Methods("POST")
//     router.HandleFunc("/favourites", getFavouriteStocks).Methods("GET")
//     router.HandleFunc("/favourites/{code}", removeFavouriteStock).Methods("DELETE")

//     log.Println("Server started on :8080")
//     http.ListenAndServe(":8080", router)
// }

func readCSV(filePath string) ([]Stock, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    _, err = reader.Read() // Assuming first row is headers
    if err != nil {
        return nil, err
    }

    var stocks []Stock
    for {
        record, err := reader.Read()
        if err != nil {
            break // EOF or error
        }

        open, _ := strconv.ParseFloat(record[2], 64)
        high, _ := strconv.ParseFloat(record[3], 64)
        low, _ := strconv.ParseFloat(record[4], 64)
        close, _ := strconv.ParseFloat(record[5], 64)

        stocks = append(stocks, Stock{
            Code:  record[0],
            Name:  record[1],
            Open:  open,
            High:  high,
            Low:   low,
            Close: close,
        })
    }

    return stocks, nil
}

func insertStocks(db *sql.DB, stocks []Stock) {
    for _, stock := range stocks {
        _, err := db.Exec("INSERT INTO stocks (code, name, open, high, low, close) VALUES (?, ?, ?, ?, ?, ?)",
            stock.Code, stock.Name, stock.Open, stock.High, stock.Low, stock.Close)

        if err != nil {
            fmt.Printf("Error inserting record: %v\n", err)
        }
    }
}

