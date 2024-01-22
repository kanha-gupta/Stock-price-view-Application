package api

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func ApiInitialise(db1 *sql.DB) {
	db = db1
	router := mux.NewRouter()
	router.HandleFunc("/stocks/top", getTopStocks).Methods("GET")               //working fine curl http://localhost:8080/stocks/top
	router.HandleFunc("/stocks/search", findStocksByName).Methods("GET")        //working fine curl http://localhost:8080/stocks/search?name=<stock_name>
	router.HandleFunc("/stocks/history/{code}", getStockHistory).Methods("GET") //working fine curl http://localhost:8080/stocks/history/500206  (last is stock code)

	router.HandleFunc("/favourites", addFavouriteStock).Methods("POST")             //curl -X POST http://localhost:8080/favourites -H "Content-Type: application/json" -d "{\"Code\": \"504988\"}"
	router.HandleFunc("/favourites", getFavouriteStocks).Methods("GET")             //curl -X GET http://localhost:8080/favourites 504988
	router.HandleFunc("/favourites/{code}", removeFavouriteStock).Methods("DELETE") //curl -X DELETE http://localhost:8080/favourites/504988

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", router)
}
