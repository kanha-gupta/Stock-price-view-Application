package main

import (
	"archive/zip"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	api "github.com/kanha-gupta/stockapp/API"
	"github.com/kanha-gupta/stockapp/dataProcessing"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	url := generateBSEURL(200124)
	downloadFile(url, "C:\\Users\\asus\\Desktop\\backendTaskInitial\\downloaded")

	db, err := sql.Open("mysql", "ODBC:0@tcp(localhost:3306)/app_database")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stocks, err := dataProcessing.ReadCSV("EQ190124.CSV")
	if err != nil {
		panic(err)
	}
	dataProcessing.InsertStocks(db, stocks)
	api.ApiInitialise(db)
}
func generateBSEURL(n int) string {
	formattedDate := n
	return fmt.Sprintf("https://www.bseindia.com/download/BhavCopy/Equity/EQ%s_CSV.ZIP", formattedDate) //200124
}

func downloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractZip extracts a ZIP file to a destination directory
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}
