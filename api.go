package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("mysql", "ODBC:0@tcp(localhost:3306)/app_database")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    router := mux.NewRouter()
    router.HandleFunc("/stocks/top", getTopStocks).Methods("GET") //working fine curl http://localhost:8080/stocks/top
    router.HandleFunc("/stocks/search", findStocksByName).Methods("GET") //working fine curl http://localhost:8080/stocks/search?name=<stock_name>
	router.HandleFunc("/stocks/history/{code}", getStockHistory).Methods("GET") //working fine curl http://localhost:8080/stocks/history/500206  (last is stock code)
    router.HandleFunc("/favourites", addFavouriteStock).Methods("POST") //not working
    router.HandleFunc("/favourites", getFavouriteStocks).Methods("GET") //not working
    router.HandleFunc("/favourites/{code}", removeFavouriteStock).Methods("DELETE")

    log.Println("Server started on :8080")
    http.ListenAndServe(":8080", router)
}

type Stock struct {
	ID    int 
    Code  string
    Name  string
    Open  float64
    High  float64
    Low   float64
    Close float64
}

func getTopStocks(w http.ResponseWriter, r *http.Request) { //working fine
    rows, err := db.Query("SELECT id, code, name, open, high, low, close FROM stocks ORDER BY close DESC LIMIT 10")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var stocks []Stock
    for rows.Next() {
        var s Stock
        if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Open, &s.High, &s.Low, &s.Close); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        stocks = append(stocks, s)
    }

    json.NewEncoder(w).Encode(stocks)
}


func findStocksByName(w http.ResponseWriter, r *http.Request) { //working fine 
    name := r.URL.Query().Get("name")
    rows, err := db.Query("SELECT * FROM stocks WHERE name LIKE ?", "%"+name+"%")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // ... Similar to getTopStocks ...
	var stocks []Stock
    for rows.Next() {
        var s Stock
        if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Open, &s.High, &s.Low, &s.Close); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        stocks = append(stocks, s)
    }

    json.NewEncoder(w).Encode(stocks)
}

func getStockHistory(w http.ResponseWriter, r *http.Request) { //working //since date is not in column, 
	//list stock prices in increasing order of their id, which can serve as a proxy for time if we assume that records were added in chronological order.
    vars := mux.Vars(r)
    stockCode := vars["code"]

    query := "SELECT id, open, high, low, close FROM stocks WHERE code = ? ORDER BY id"
    rows, err := db.Query(query, stockCode)
    if err != nil {
        http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var history []Stock
    for rows.Next() {
        var h Stock
        if err := rows.Scan(&h.ID, &h.Open, &h.High, &h.Low, &h.Close); err != nil {
            http.Error(w, "Error scanning database rows: "+err.Error(), http.StatusInternalServerError)
            return
        }
        history = append(history, h)
    }

    json.NewEncoder(w).Encode(history)
}

func addFavouriteStock(w http.ResponseWriter, r *http.Request) {
    var s Stock
    if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err := db.Exec("INSERT INTO favourites (code) VALUES (?)", s.Code)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Added to favourites")
}

func removeFavouriteStock(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    code := vars["code"]

    _, err := db.Exec("DELETE FROM favourites WHERE code = ?", code)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "Removed from favourites")
}

// func getStockHistory(w http.ResponseWriter, r *http.Request) { 
//     vars := mux.Vars(r)
//     code := vars["code"]

//     rows, err := db.Query("SELECT date, open, high, low, close FROM stock_history WHERE code = ? ORDER BY date DESC", code)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     defer rows.Close()

//     var history []Stock
//     for rows.Next() {
//         var h Stock
//         if err := rows.Scan(&h.Date, &h.Open, &h.High, &h.Low, &h.Close); err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }
//         history = append(history, h)
//     }

//     json.NewEncoder(w).Encode(history)
// }

func getFavouriteStocks(w http.ResponseWriter, r *http.Request) {
    // Example: Fetching favourites for a hardcoded user.
    userID := "exampleUserID" // Replace with actual user ID logic

    rows, err := db.Query("SELECT s.code, s.name, s.open, s.high, s.low, s.close FROM stocks s JOIN favourites f ON s.code = f.code WHERE f.id = ?", userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var favourites []Stock
    for rows.Next() {
        var s Stock
        if err := rows.Scan(&s.Code, &s.Name, &s.Open, &s.High, &s.Low, &s.Close); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        favourites = append(favourites, s)
    }

    json.NewEncoder(w).Encode(favourites)
}




