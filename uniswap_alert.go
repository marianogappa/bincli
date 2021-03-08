package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	uniswapAliases = map[string]string{
		"APYS": "0xf7413489c474ca4399eee604716c72879eea3615",
		"SYNC": "0xb6ff96b8a8d214544ca0dbc9b33f7ad6503efd32",
	}
)

type tokenDayData struct {
	PriceUSD string `json:"priceUSD"`
}

type tokenDayDatas struct {
	TokenDayDatas []tokenDayData `json:"tokenDayDatas"`
}

type response struct {
	Data tokenDayDatas `json:"data"`
}

func isUniswapConditionMet(contractAddress string, comparator string, target float64) (float64, bool) {
	url := "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"

	var jsonStr = []byte(`{"query": "{tokenDayDatas(first: 1, orderBy: date, orderDirection: desc, where: { token: \"` + contractAddress + `\"}) {priceUSD } }"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseData := response{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Fatal(err)
	}
	strPrice := responseData.Data.TokenDayDatas[0].PriceUSD
	price, err := strconv.ParseFloat(strPrice, 64)
	if err != nil {
		log.Fatal(err)
	}

	switch comparator {
	case ">":
		if price > target {
			return price, true
		}
	case "<":
		if price < target {
			return price, true
		}
	case ">=":
		if price >= target {
			return price, true
		}
	case "<=":
		if price <= target {
			return price, true
		}
	default:
		usage()

	}
	return price, false
}

func uniswapAlert() {
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
	if uniswapAliases[symbol] != "" {
		contractAddress = uniswapAliases[symbol]
	}

	for {
		price, isConditionMet := isUniswapConditionMet(contractAddress, comparator, target)
		log.Printf("%v %v %v where %v = %v...\n",
			symbol,
			comparator,
			targetStr,
			symbol,
			price,
		)
		if isConditionMet {
			break
		}
		time.Sleep(10 * time.Second)
	}
	log.Println("Condition met!")
}
