package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ftxResponseResult struct {
	Last float64 `json:"last"`
}

type ftxResponse struct {
	Result ftxResponseResult `json:"result"`
}

func ftxTicker() {
	if len(os.Args) < 3 {
		usage()
	}
	price, err := requestFtxTicker(os.Args[2])
	if err != nil {
		log.Printf("Error getting ticker price for %v (because %v)", os.Args[2], err)
		usage()
	}
	fmt.Println(price)
}

func requestFtxTicker(coin string) (float64, error) {
	url := fmt.Sprintf("https://ftx.com./api/markets/%v", coin)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseData := ftxResponse{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return 0, err
	}
	price := responseData.Result.Last
	if err != nil {
		return 0, err
	}
	return price, nil
}

func isFtxConditionMet(contractAddress string, comparator string, target float64) (float64, bool, error) {
	price, err := requestFtxTicker(contractAddress)
	if err != nil {
		return price, false, err
	}
	return isConditionMet(price, comparator, target)
}

func ftxAlert() {
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
		price, isConditionMet, err := isFtxConditionMet(contractAddress, comparator, target)
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
