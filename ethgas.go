package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ethGasResponse struct {
	Code int            `json:"code"`
	Data map[string]int `json:"data"`
}

func requestEthGas() (map[string]int, error) {
	resp, err := http.Get("https://www.gasnow.org/api/v3/gas/price?utm_source=bincli")
	if err != nil {
		return map[string]int{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	responseData := ethGasResponse{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return map[string]int{}, err
	}
	if responseData.Code != 200 {
		return map[string]int{}, fmt.Errorf("gasnow.org returned non-200: %v", responseData.Code)
	}
	return responseData.Data, nil
}

func ethGas() {
	if len(os.Args) < 3 {
		usage()
	}
	var (
		qos = os.Args[2]
	)

	gas, err := requestEthGas()
	if err != nil {
		log.Println(err)
		usage()
	}

	value, ok := gas[qos]
	if !ok {
		log.Printf("Unsupported QoS: %v. Choose one of {rapid|fast|standard|slow}\n", qos)
		usage()
	}

	fmt.Println(value / 1000000000)
}
