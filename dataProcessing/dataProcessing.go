package dataProcessing

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/kanha-gupta/stockapp/database-structure"
	"os"
	"strconv"
)

func ReadCSV(filePath string) ([]database_structure.Stock, error) {
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

	var stocks []database_structure.Stock
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		open, _ := strconv.ParseFloat(record[4], 64)
		high, _ := strconv.ParseFloat(record[5], 64)
		low, _ := strconv.ParseFloat(record[6], 64)
		close, _ := strconv.ParseFloat(record[7], 64)
		stocks = append(stocks, database_structure.Stock{
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

func InsertStocks(db *sql.DB, stocks []database_structure.Stock) {
	for _, stock := range stocks {
		_, err := db.Exec("INSERT INTO stocks (code, name, open, high, low, close) VALUES (?, ?, ?, ?, ?, ?)",
			stock.Code, stock.Name, stock.Open, stock.High, stock.Low, stock.Close)

		if err != nil {
			fmt.Printf("Error inserting record: %v\n", err)
		}
	}
}
