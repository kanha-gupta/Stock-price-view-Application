# Stock Price View Application

## Description
This application is designed to interact with the BSE to download, 
process, and store stock data. 
It provides a RESTful API for accessing this data, 
including features like viewing top stocks, searching stocks by name, 
and managing a list of favourite stocks.

## Features

The application offers a range of functionalities to manage and access stock data:

- **CSV File Download:** Downloads the CSV file in ZIP format from a given URL.
- **CSV Extraction and Reading:** Extracts and reads the CSV file contained within the ZIP.
- **Data Storage:** Stores the data in a MySQL database.
- **Historical Data Support:** Enables fetching of data from the last 50 days.
- **Favourites Management:**
   - Addition of stocks to favourites.
   - Deletion of stocks from favourites.
   - Get information of stocks in favourites.
- **Stock Search:** Allows searching for stocks by name.
- **Stock History Access:** Facilitates access to the historical data of stocks.
- **API Command Reference:** For detailed API commands, [click here](#stock-api-commands).


# How to Setup

This guide outlines the steps to set up the necessary database and tables for the stock management API.
## 1. Docker (Build from source code)
### Steps:
1. **Run the following commands to set up**
    ```
   docker-compose build --no-cache
   docker compose up
   ```

## 2. Local environment

### Requirements : 

1. GoLang (version 1.21.2 used in the project)
2. MySQL

### Steps:
1. **Create a Database (MySQL)**

   Use the following MySQL command to create a new database:

   ```sql
   CREATE DATABASE app_database;
   ```

2. **Create Two Tables Named `stocks` and `favourites` (MySQL)**

    - To create the `stocks` table, use the following command:

      ```sql
      CREATE TABLE stocks (
      id INT AUTO_INCREMENT PRIMARY KEY,
      code VARCHAR(10),
      name VARCHAR(255),
      open DECIMAL(10, 2),
      high DECIMAL(10, 2),
      low DECIMAL(10, 2),
      close DECIMAL(10, 2)
      );
      ```

    - To create the `favourites` table, use the following command:

      ```sql
      CREATE TABLE favourites (
      id INT AUTO_INCREMENT PRIMARY KEY,
      code VARCHAR(255) NOT NULL
      );
      ```

3. **Run `go run main.go [URL]` to Start the Server**

    - Example : go run main.go https://www.bseindia.com/download/BhavCopy/Equity/EQ_ISINCODE_250124.zip
     - After setting up the database and tables, run the `main.go` file to start the server. Ensure you have put your preferred URL in the configuration.


# Stock API Commands

This document provides a detailed description of various `curl` commands used to interact with a local stock management API running on `http://localhost:8080`.

### 1. Get Top Stocks
```bash
curl http://localhost:8080/stocks/top
```
**Description:** Retrieves a list of top-performing stocks. This command sends a GET request to the `/stocks/top` endpoint, which returns information about the stocks that are currently leading in the market.

### 2. Search for a Specific Stock
```bash
curl http://localhost:8080/stocks/search?name=<stock_name>
```
**Description:** Searches for a specific stock by name. Replace `<stock_name>` with the actual name of the stock you are interested in. This command sends a GET request to the `/stocks/search` endpoint with a query parameter to specify the stock name.

### 3. Get Stock History
```bash
curl http://localhost:8080/stocks/history/500206
```
**Description:** Retrieves the historical data of a specific stock. The last part of the command (`500206`) is the stock code. Replace this with the code of the stock whose history you want to access. This command sends a GET request to the `/stocks/history/{stock_code}` endpoint.

### 4. Add Stock to Favourites
```bash
curl -X POST http://localhost:8080/favourites -H "Content-Type: application/json" -d "{"Code": "504988"}"
```
**Description:** Adds a stock to your list of favourites. Replace `504988` with the stock code you wish to add. This command sends a POST request to the `/favourites` endpoint with a JSON payload containing the stock code.

### 5. Get Specific Favourite Stock
```bash
curl -X GET http://localhost:8080/favourites/504988
```
**Description:** Retrieves information about a specific favourite stock. Replace `504988` with the stock code you want to inquire about. This command sends a GET request to the `/favourites/{stock_code}` endpoint.

### 6. Remove Stock from Favourites
```bash
curl -X DELETE http://localhost:8080/favourites/504988
```
**Description:** Removes a stock from your favourites list. Replace `504988` with the stock code you wish to remove. This command sends a DELETE request to the `/favourites/{stock_code}` endpoint.

---

**Note:** Ensure that your local server is running on `http://localhost:8080` before executing these commands. Replace stock codes and names with actual values as required.


