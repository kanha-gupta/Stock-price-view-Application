package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kanha-gupta/stockapp/database-structure"
	"net/http"
)

func getTopStocks(w http.ResponseWriter, r *http.Request) { //working fine
	rows, err := db.Query("SELECT id, code, name, open, high, low, close FROM stocks ORDER BY close DESC LIMIT 10")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stocks []database_structure.Stock
	for rows.Next() {
		var s database_structure.Stock
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
	var stocks []database_structure.Stock
	for rows.Next() {
		var s database_structure.Stock
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

	var history []database_structure.Stock
	for rows.Next() {
		var h database_structure.Stock
		if err := rows.Scan(&h.ID, &h.Open, &h.High, &h.Low, &h.Close); err != nil {
			http.Error(w, "Error scanning database rows: "+err.Error(), http.StatusInternalServerError)
			return
		}
		history = append(history, h)
	}

	json.NewEncoder(w).Encode(history)
}

func addFavouriteStock(w http.ResponseWriter, r *http.Request) {
	var s database_structure.Stock
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
	stockCode := vars["code"]

	// SQL statement to delete the stock from favourites.
	_, err := db.Exec("DELETE FROM favourites WHERE code = ?", stockCode)
	if err != nil {
		// Handle the error, send an appropriate response.
		http.Error(w, "Error removing stock from favourites: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Sending a success response.
	fmt.Fprintln(w, "Stock removed from favourites")
}

func getFavouriteStocks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT s.id, s.code, s.name, s.open, s.high, s.low, s.close 
        FROM stocks s 
        JOIN favourites f ON s.code = f.code
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var favourites []database_structure.Stock
	for rows.Next() {
		var s database_structure.Stock
		if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Open, &s.High, &s.Low, &s.Close); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		favourites = append(favourites, s)
	}

	json.NewEncoder(w).Encode(favourites)
}
