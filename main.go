package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change string
}

type HistoricalData struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}


const alphaVantageAPIKey = "JY7A63216ER5NGBK"

func main() {
	fmt.Println("Welcome to Stock Viewer!")
	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. View Stocks")
		fmt.Println("2. Download Stocks Data")
		fmt.Println("3. View Historical Data")
		fmt.Println("4. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			log.Fatalln("Failed to read choice:", err)
		}

		switch choice {
		case 1:
			viewStocks()
		case 2:
			downloadStocks()
		case 3:
			viewHistoricalData()
		case 4:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func viewStocks() {
	ticker := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"VMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	stocks := scrapeStocks(ticker)

	fmt.Println("Stocks:")
	for _, stock := range stocks {
		fmt.Printf("Company: %s, Price: %s, Change: %s\n", stock.company, stock.price, stock.change)
	}
}

func downloadStocks() {
	ticker := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"VMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	stocks := scrapeStocks(ticker)

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Failed to create output csv file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Company", "Price", "Change"}
	writer.Write(headers)
	for _, stock := range stocks {
		record := []string{stock.company, stock.price, stock.change}
		writer.Write(record)
	}

	fmt.Println("Stocks data downloaded to stocks.csv file successfully")
}

func viewHistoricalData() {
	fmt.Println("Viewing Historical Data...")

	symbol := "AAPL"     // Example symbol (Apple Inc.)
	timeRange := "1month" // Example time range

	historicalData, err := fetchHistoricalData(symbol, timeRange)
	if err != nil {
		log.Println("Error fetching historical data:", err)
		return
	}

	fmt.Println("Historical Data:")
	for _, data := range historicalData {
		fmt.Printf("Date: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %d\n",
			data.Date, data.Open, data.High, data.Low, data.Close, data.Volume)
	}
}

func fetchHistoricalData(symbol, timeRange string) ([]HistoricalData, error) {
    resp, err := resty.New().R().
        SetQueryParam("function", "TIME_SERIES_DAILY").
        SetQueryParam("symbol", symbol).
        SetQueryParam("outputsize", "compact").
        SetQueryParam("apikey", alphaVantageAPIKey).
        Get("https://www.alphavantage.co/query")
    if err != nil {
        return nil, fmt.Errorf("failed to fetch data: %v", err)
    }

    // Parse response body
    var dataResponse map[string]interface{}
    err = json.Unmarshal(resp.Body(), &dataResponse)
    if err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }

    // Print response body
    fmt.Println("Response Body:", string(resp.Body()))

    // Implement parsing based on response body structure

    return nil, nil
}

func scrapeStocks(tickers []string) []Stock {
	stocks := []Stock{}

	c := colly.NewCollector()
	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{}

		stock.company = e.ChildText("h1")
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")

		stocks = append(stocks, stock)
	})

	for _, t := range tickers {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}

	return stocks
}
