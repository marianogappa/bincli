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
	honeyswapAliases = map[string]string{
		"AGVE": "0x3a97704a1b25f08aa230ae53b352e2e72ef52843",
	}
)

func isHoneyswapConditionMet(contractAddress string, comparator string, target float64) (float64, bool) {
	url := "https://api.thegraph.com/subgraphs/name/1hive/uniswap-v2"

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

func honeyswapAlert() {
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
	if honeyswapAliases[symbol] != "" {
		contractAddress = honeyswapAliases[symbol]
	}

	for {
		price, isConditionMet := isHoneyswapConditionMet(contractAddress, comparator, target)
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
