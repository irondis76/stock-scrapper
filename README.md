# Stock-Scrapper-Viewer

This is a simple command-line application written in Go that allows users to view and download stock information. It provides options to view live stock prices, download stock data to a CSV file, and view historical stock data.

## Installation

1. Make sure you have Go installed on your system.
2. Clone this repository to your local machine.
3. Navigate to the project directory.
4. Run `go build` to build the executable.
5. Run the executable file generated.

## Usage

1. Upon running the executable, you will be presented with a menu.
2. Choose an option by entering the corresponding number and pressing Enter.

## Options

1. **View Stocks**: Displays live stock information for a predefined list of stocks.
2. **Download Stocks Data**: Scrapes live stock information and saves it to a CSV file named `stocks.csv`.
3. **View Historical Data**: Fetches historical stock data for a specified symbol and time range.
4. **Exit**: Exits the application.

## Dependencies

This application uses the following dependencies:

- [github.com/go-resty/resty/v2](https://github.com/go-resty/resty/v2): A simple HTTP and REST client library for Go.
- [github.com/gocolly/colly](https://github.com/gocolly/colly): An elegant and efficient web scraping framework for Go.

### Install the required Go modules:
```
go get github.com/go-resty/resty/v2
```
```
go get github.com/gocolly/colly
```


## API Key

The application uses the Alpha Vantage API to fetch historical stock data. You need to provide your own API key by replacing `alphaVantageAPIKey` constant in the code.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

