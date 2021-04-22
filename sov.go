package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func sovTicker() {
	if len(os.Args) < 2 {
		usage()
	}
	price, err := requestSovTicker()
	if err != nil {
		log.Printf("Error getting ticker price for SOV (because %v)", err)
		usage()
	}
	fmt.Println(price)
}

func requestSovTicker() (float64, error) {
	url := fmt.Sprintf("https://coins.green")

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	strBody = strings.Replace(strBody, "\n", "", -1)

	rx := regexp.MustCompile(`<span data-sov-usd-price>(.+?)<\/span>`)
	sm := rx.FindStringSubmatch(strBody)
	if len(sm) < 2 {
		return 0, errors.New("Couldn't find SOV value on coins.green. They probably changed their website!")
	}

	strPrice := sm[1]
	price, err := strconv.ParseFloat(strPrice, 64)
	return price, err
}

func isSovConditionMet(comparator string, target float64) (float64, bool, error) {
	price, err := requestSovTicker()
	if err != nil {
		return price, false, err
	}
	return isConditionMet(price, comparator, target)
}

func sovAlert() {
	if len(os.Args) < 4 {
		usage()
	}
	var (
		comparator  = os.Args[2]
		targetStr   = os.Args[3]
		target, err = strconv.ParseFloat(targetStr, 64)
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		price, isConditionMet, err := isSovConditionMet(comparator, target)
		if err != nil {
			log.Printf("Error getting ticker price for SOV (because %v)", err)
		} else {
			log.Printf("SOV %v %v where SOV = %v...\n",
				comparator,
				targetStr,
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
