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

func usage() {
	log.Fatal(`Subcommands:

balance
balance $asset (e.g. balance BTC)
ticker
ticker $market (e.g. ticker BTCUSDT)
alert $market $comparator $value (e.g. alert BTCUSDT ">" 56000)
	`)
}

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

func alert(client *binance.Client, futuresClient *futures.Client) {
	if len(os.Args) < 5 {
		usage()
	}
	var (
		doneC         = make(chan struct{})
		aggTradeCount = 0
		symbol        = os.Args[2]
		comparator    = os.Args[3]
		valueStr      = os.Args[4]
		value, err    = strconv.ParseFloat(valueStr, 64)
	)
	if err != nil {
		log.Fatal(err)
	}
	wsAggTradeHandler := func(event *binance.WsAggTradeEvent) {
		aggTradeCount++
		if aggTradeCount%100 == 1 {
			log.Printf("Attempt %v: %v %v %v where %v = %v...\n",
				aggTradeCount,
				symbol,
				comparator,
				valueStr,
				symbol,
				event.Price,
			)
		}
		price, err := strconv.ParseFloat(event.Price, 64)
		if err != nil {
			log.Fatal(err)
		}
		switch comparator {
		case ">":
			if price > value {
				doneC <- struct{}{}
			}
		case "<":
			if price < value {
				doneC <- struct{}{}
			}
		case ">=":
			if price >= value {
				doneC <- struct{}{}
			}
		case "<=":
			if price <= value {
				doneC <- struct{}{}
			}
		default:
			usage()
		}
	}
	errHandler := func(err error) {
		log.Fatal(err)
	}
	_, _, err = binance.WsAggTradeServe(symbol, wsAggTradeHandler, errHandler)
	if err != nil {
		log.Fatal(err)
	}
	<-doneC
	log.Printf("Condition %v %v %v reached!",
		symbol,
		comparator,
		valueStr,
	)
}

func main() {
	var (
		apiKey        = os.Getenv("BINANCE_RO_API_KEY")
		secretKey     = os.Getenv("BINANCE_RO_SECRET")
		client        = binance.NewClient(apiKey, secretKey)
		futuresClient = binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
	)

	if apiKey == "" || secretKey == "" {
		log.Println("Please set BINANCE_RO_API_KEY & BINANCE_RO_SECRET. I can't do any account-related queries.")
	}

	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "balance":
		balance(client, futuresClient)
	case "ticker":
		ticker(client, futuresClient)
	case "alert":
		alert(client, futuresClient)
	default:
		usage()
	}
}
