package main

import (
	"log"
	"os"

	"github.com/adshao/go-binance/v2"
)

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
		balance(client)
	case "chartBalance":
		chartBalance(client)
	case "ticker":
		ticker(client, futuresClient)
	case "alert":
		alert(client, futuresClient)
	case "uniswapAlert":
		uniswapAlert()
	case "honeyswapAlert":
		honeyswapAlert()
	case "bitforexAlert":
		bitforexAlert()
	case "bitmaxAlert":
		bitmaxAlert()
	default:
		usage()
	}
}
