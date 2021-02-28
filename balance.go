package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

func balance(client *binance.Client, futuresClient *futures.Client) {
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	switch len(os.Args) {
	case 2:
		bs, err := json.Marshal(account.Balances)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bs))
	case 3:
		for _, balance := range account.Balances {
			if balance.Asset == strings.ToUpper(os.Args[2]) {
				free, err := strconv.ParseFloat(balance.Free, 64)
				if err != nil {
					log.Fatal(err)
				}
				locked, err := strconv.ParseFloat(balance.Locked, 64)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(free + locked)
			}
		}
	}
}
