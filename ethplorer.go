package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type ethplorerPrice struct {
	Rate            float64 `json:"rate"`
	Diff            float64 `json:"diff"`
	Diff7d          float64 `json:"diff7d"`
	Ts              float64 `json:"ts"`
	MarketCapUsd    float64 `json:"marketCapUsd"`
	AvailableSupply float64 `json:"availableSupply"`
	Volume24h       float64 `json:"volume24h"`
	Diff30d         float64 `json:"diff30d"`
	VolDiff1        float64 `json:"volDiff1"`
	VolDiff7        float64 `json:"volDiff7"`
	VolDiff30       float64 `json:"volDiff30"`
	Currency        string  `json:"currency"`
}

type ethplorerETH struct {
	Balance float64        `json:"balance"`
	Price   ethplorerPrice `json:"price"`
}

type ethplorerTokenInfo struct {
	Address           string          `json:"address"`
	Decimals          json.RawMessage `json:"decimals"`
	Name              string          `json:"name"`
	Symbol            string          `json:"symbol"`
	TotalSupply       string          `json:"totalSupply"`
	LastUpdated       int             `json:"lastUpdated"`
	IssuancesCount    int             `json:"issuancesCount"`
	HoldersCount      int             `json:"holdersCount"`
	Website           string          `json:"website"`
	Telegram          string          `json:"telegram"`
	Twitter           string          `json:"twitter"`
	Image             string          `json:"image"`
	Coingecko         string          `json:"coingecko"`
	EthTransfersCount int             `json:"ethTransfersCount"`
	Price             json.RawMessage `json:"price"`
}

type ethplorerToken struct {
	TokenInfo ethplorerTokenInfo `json:"tokenInfo"`
	Balance   float64            `json:"balance"`
	TotalIn   float64            `json:"totalIn"`
	TotalOut  float64            `json:"totalOut"`
}

type ethplorerResponse struct {
	Address  string           `json:"address"`
	Eth      ethplorerETH     `json:"ETH"`
	Tokens   []ethplorerToken `json:"tokens"`
	CountTxs int              `json:"countTxs"`
}

func (r ethplorerResponse) usdBalances() map[string]float64 {
	usdBalances := map[string]float64{}
	for _, token := range r.Tokens {
		decimalsStr := string(token.TokenInfo.Decimals)
		if len(decimalsStr) > 2 && decimalsStr[0] == '"' && decimalsStr[len(decimalsStr)-1] == '"' {
			decimalsStr = decimalsStr[1 : len(decimalsStr)-1]
		}
		decimals, err := strconv.Atoi(decimalsStr)
		if err != nil {
			log.Println(err)
			continue
		}
		price := ethplorerPrice{}
		err = json.Unmarshal(token.TokenInfo.Price, &price)
		if err != nil {
			continue
		}
		usdBalances[token.TokenInfo.Symbol] = token.Balance * price.Rate / math.Pow(10, float64(decimals))
	}
	usdBalances["ETH"] = r.Eth.Balance * r.Eth.Price.Rate
	return usdBalances
}

func requestEthplorer(address string) (ethplorerResponse, error) {
	url := fmt.Sprintf("https://api.ethplorer.io/getAddressInfo/%v?apiKey=freekey", address)

	resp, err := http.Get(url)
	if err != nil {
		return ethplorerResponse{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseData := ethplorerResponse{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return ethplorerResponse{}, err
	}
	return responseData, nil
}

func ethBalance() {
	if len(os.Args) < 3 {
		usage()
	}
	var (
		address = os.Args[2]
	)

	resp, err := requestEthplorer(address)
	if err != nil {
		log.Println(err)
		usage()
	}

	for symbol, usdBalance := range resp.usdBalances() {
		fmt.Printf("%v\t%v\n", symbol, usdBalance)
	}
}
