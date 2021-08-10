package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	sushiswapAliases = map[string]string{
		"DELTA": "0x9ea3b5b4ec044b70375236a281986106457b20ef",
		"RLP":   "0xfcfc434ee5bff924222e084a8876eee74ea7cfba",
		"YGG":   "0x25f8087ead173b73d6e8b84329989a8eea16cf73",
		"CVX":   "0x4e3fbd56cd56c3e72c1403e103b45db9da5b9d2b",
	}
)

func sushiswapTicker() {
	if len(os.Args) < 3 {
		usage()
	}
	symbol := os.Args[2]
	contractAddress := symbol
	if sushiswapAliases[symbol] != "" {
		contractAddress = sushiswapAliases[symbol]
	}

	price, err := requestSushiswapTicker(contractAddress)
	if err != nil {
		log.Printf("Error getting ticker price for %v (because %v)", contractAddress, err)
		usage()
	}
	fmt.Println(price)
}

func requestSushiswapTicker(contractAddress string) (float64, error) {
	url := "https://api.thegraph.com/subgraphs/name/sushiswap/exchange"

	var jsonStr = []byte(`{"query": "{tokenDayDatas(first: 1, orderBy: date, orderDirection: desc, where: { token: \"` + contractAddress + `\"}) {priceUSD } }"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseData := response{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return 0, err
	}
	strPrice := responseData.Data.TokenDayDatas[0].PriceUSD
	price, err := strconv.ParseFloat(strPrice, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func isSushiswapConditionMet(contractAddress string, comparator string, target float64) (float64, bool, error) {
	price, err := requestSushiswapTicker(contractAddress)
	if err != nil {
		return price, false, err
	}
	return isConditionMet(price, comparator, target)
}

func sushiswapAlert() {
	if len(os.Args) < 5 {
		usage()
	}
	var (
		symbol      = os.Args[2]
		comparator  = os.Args[3]
		targetStr   = os.Args[4]
		target, err = strconv.ParseFloat(targetStr, 64)
	)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := symbol
	if sushiswapAliases[symbol] != "" {
		contractAddress = sushiswapAliases[symbol]
	}

	for {
		price, isConditionMet, err := isSushiswapConditionMet(contractAddress, comparator, target)
		if err != nil {
			log.Printf("Error getting ticker price for %v (because %v)", symbol, err)
		} else {
			log.Printf("%v %v %v where %v = %v...\n",
				symbol,
				comparator,
				targetStr,
				symbol,
				price,
			)
		}
		if isConditionMet {
			break
		}
		time.Sleep(10 * time.Second)
	}
	log.Println("Condition met!")
}
