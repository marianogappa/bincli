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

type bitmaxClose struct {
	Close string `json:"close"`
}

type bitmaxResponse struct {
	Data bitmaxClose `json:"data"`
}

func isBitmaxConditionMet(coin string, comparator string, target float64) (float64, bool, error) {
	url := fmt.Sprintf("https://bitmax.io/api/pro/v1/ticker?symbol=%v", coin)

	resp, err := http.Get(url)
	if err != nil {
		return 0, false, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseData := bitmaxResponse{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return 0, false, err
	}
	priceStr := responseData.Data.Close
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, false, err
	}

	switch comparator {
	case ">":
		if price > target {
			return price, true, nil
		}
	case "<":
		if price < target {
			return price, true, nil
		}
	case ">=":
		if price >= target {
			return price, true, nil
		}
	case "<=":
		if price <= target {
			return price, true, nil
		}
	default:
		usage()

	}
	return price, false, nil
}

func bitmaxAlert() {
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
		price, isConditionMet, err := isBitmaxConditionMet(contractAddress, comparator, target)
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
