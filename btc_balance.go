package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adshao/go-binance/v2"
)

func btcBalance(client *binance.Client) {
	balances := mustRequestAccountBalances(client)
	ticker := mustRequestTicker(client)

	switch len(os.Args) {
	case 2:
		balancesInBtc := mustCalculateAllBalancesInBtc(balances, ticker, false)
		bs, err := json.Marshal(balancesInBtc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bs))
	case 3:
		fmt.Println(mustCalculateBalanceInBtc(strings.ToUpper(os.Args[2]), balances, ticker).balance)
	default:
		usage()
	}
}
