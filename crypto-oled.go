package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	baseURL = "https://api.coinmarketcap.com/v1"
	url     string
)

//Coin struct
type Coin struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank,string"`
	PriceUsd         float64 `json:"price_usd,string"`
	PriceBtc         float64 `json:"price_btc,string"`
	Usd24hVolume     float64 `json:"24h_volume_usd,string"`
	MarketCapUsd     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1h  float64 `json:"percent_change_1h,string"`
	PercentChange24h float64 `json:"percent_change_24h,string"`
	PercentChange7d  float64 `json:"percent_change_7d,string"`
	LastUpdated      string  `json:"last_updated"`
}

func main() {

	// Get info about coin
	btcInfo, err := getCoinData("bitcoin")
	if err != nil {
		log.Println(err)
	} else {
		printCoin(btcInfo)
	}

	fmt.Print("\n")

	ethInfo, err := getCoinData("ethereum")
	if err != nil {
		log.Println(err)
	} else {
		printCoin(ethInfo)
	}
}

func parseUnixTime(updatetime string) time.Time {
	i, err := strconv.ParseInt(updatetime, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm
}

func printCoin(coin Coin) {
	fmt.Println("Name: " + coin.Name)
	fmt.Print("Price USD: $")
	fmt.Printf("%.2f\n", coin.PriceUsd)
	fmt.Printf("1 hour Δ: %.2f%% \n", coin.PercentChange1h)
	fmt.Printf("24 hour Δ: %.2f%% \n", coin.PercentChange24h)
	fmt.Printf("7 day Δ: %.2f%% \n", coin.PercentChange7d)
	fmt.Print("Last Updated: ")
	fmt.Println(parseUnixTime(coin.LastUpdated))

	timestamp := time.Now().UTC().Unix()
	cointimestamp, err := strconv.ParseInt(coin.LastUpdated, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	diffstamp := timestamp - cointimestamp
	if diffstamp > 300 {
		fmt.Println("WARNING: Prices are > 5 mins old")
	}
}

// GetCoinData Gets information about a crypto currency.
func getCoinData(coin string) (Coin, error) {
	url = fmt.Sprintf("%s/ticker/%s", baseURL, coin)
	resp, err := makeReq(url)
	if err != nil {
		return Coin{}, err
	}
	var data []Coin
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return Coin{}, err
	}

	return data[0], nil
}

// HTTP Request Helper
func makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := doReq(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// HTTP Client
func doReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}
