package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	api "github.com/kanha-gupta/stockapp/API"
	"github.com/kanha-gupta/stockapp/dataProcessing"
	"os"
)

func main() {
	defaultUrl := "https://www.bseindia.com/download/BhavCopy/Equity/EQ250124_CSV.ZIP"
	var url string
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [URL]")
		url = defaultUrl
	} else {
		url = os.Args[1]
	}
	zipFileName, err := dataProcessing.DownloadPackage(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileName, err := dataProcessing.ExtractZip("downloaded", zipFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := sql.Open("mysql", "ODBC:0@tcp(localhost:3306)/app_database")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	csvFilePath := fmt.Sprintf("downloaded/%s", fileName)
	stocks, err := dataProcessing.ReadCSV(csvFilePath)
	if err != nil {
		panic(err)
	}
	dataProcessing.InsertStocks(db, stocks)
	api.ApiInitialise(db)
}
