package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adshao/go-binance/v2"
)

func balance(client *binance.Client) {
	balances := mustRequestAccountBalances(client)
	ticker := mustRequestTicker(client)

	switch len(os.Args) {
	case 2:
		balances := mustCalculateAllBalances(balances, ticker, false)
		bs, err := json.Marshal(balances)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bs))
	case 3:
		balance := mustCalculateBalance(strings.ToUpper(os.Args[2]), balances, ticker)
		bs, err := json.Marshal(balance)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bs))
	default:
		usage()
	}
}
