package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

func ticker(client *binance.Client, futuresClient *futures.Client) {
	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	switch len(os.Args) {
	case 2:
		bs, err := json.Marshal(prices)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bs))
	case 3:
		for _, price := range prices {
			if price.Symbol == strings.ToUpper(os.Args[2]) {
				fmt.Println(price.Price)
				break
			}
		}
	default:
		usage()
	}
}
