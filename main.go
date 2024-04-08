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

func enterSymbols() []string {
	var n int
	fmt.Println("Enter number of stocks")
	_, err := fmt.Scan(&n)
	if err != nil {
		log.Println("Error reading input:", err)
		return nil
	}

	tickers := make([]string, n)
	fmt.Println("Enter stocks")
	for i := 0; i < n; i++ {
		_, err := fmt.Scan(&tickers[i])
		if err != nil {
			log.Println("Error reading input:", err)
			return nil
		}
	}
	return tickers
}

func viewStocks() {
	tickers := enterSymbols()

	stocks := scrapeStocks(tickers)

	fmt.Println("Stocks:")
	for _, stock := range stocks {
		fmt.Printf("Company: %s, Price: %s, Change: %s\n", stock.company, stock.price, stock.change)
	}
}

func downloadStocks() {
	tickers := enterSymbols()

	stocks := scrapeStocks(tickers)

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

	var ticker string
	fmt.Println("Enter a symbol")
	_, err := fmt.Scan(&ticker)
	if err != nil {
		log.Fatalln("Error during input", err)
		return
	}

	var timeRange string
	fmt.Println("Enter time range (e.g., '1min', '5min', '15min', '30min', '60min', 'daily', 'weekly', 'monthly'): ")
	_, err = fmt.Scan(&timeRange)
	if err != nil {
		log.Fatalln("Error during input", err)
		return
	}

	historicalData, err := fetchHistoricalData(ticker, timeRange)
	if err != nil {
		log.Println("Error fetching historical data:", err)
		return
	}

	file, err := os.Create("historical_data.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Date", "Open", "High", "Low", "Close", "Volume"}
	writer.Write(headers)

	for _, data := range historicalData {
		record := []string{data.Date, fmt.Sprintf("%.2f", data.Open), fmt.Sprintf("%.2f", data.High), fmt.Sprintf("%.2f", data.Low), fmt.Sprintf("%.2f", data.Close), fmt.Sprintf("%d", data.Volume)}
		writer.Write(record)
	}

	fmt.Println("Historical data downloaded to historical_data.csv file successfully")
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

    
    var dataResponse map[string]interface{}
    err = json.Unmarshal(resp.Body(), &dataResponse)
    if err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }

  
    fmt.Println("Response Body:", string(resp.Body()))

    

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
