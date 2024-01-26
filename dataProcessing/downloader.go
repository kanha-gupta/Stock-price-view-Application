package dataProcessing

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadPackage(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error while downloading: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error: HTTP Status is %s", resp.Status)
	}

	filename := filepath.Base(url)
	out, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("Error while creating file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error while copying data to file: %v", err)
	}
	fmt.Println("File downloaded successfully as", filename)
	return filename, nil
}

func ExtractZip(destDir, zipFileName string) (string, error) {
	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		return "", fmt.Errorf("Error opening zip file: %v", err)
	}
	defer r.Close()

	for _, file := range r.File {
		fpath := filepath.Join(destDir, file.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return "", fmt.Errorf("Invalid file path: %s", fpath)
		}
		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", fmt.Errorf("Error opening file for writing: %v", err)
		}
		rc, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("Error opening file in zip: %v", err)
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return "", fmt.Errorf("Error copying data: %v", err)
		}
		return file.Name, nil
	}
	return "", nil
}
