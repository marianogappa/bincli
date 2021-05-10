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
	case "ethBalance":
		ethBalance()
	case "chartBalanceBtc":
		chartBalance(client, false)
	case "chartBalanceUsdt":
		chartBalance(client, true)
	case "chartBalanceDataBtc":
		chartBalanceData(client, false)
	case "chartBalanceDataUsdt":
		chartBalanceData(client, true)
	case "ticker":
		ticker(client, futuresClient)
	case "alert":
		alert(client, futuresClient)
	case "uniswapAlert":
		uniswapAlert()
	case "honeyswapAlert":
		honeyswapAlert()
	case "sushiswapAlert":
		sushiswapAlert()
	case "bitforexAlert":
		bitforexAlert()
	case "bitmaxAlert":
		bitmaxAlert()
	case "sovAlert":
		sovAlert()
	case "uniswapTicker":
		uniswapTicker()
	case "honeyswapTicker":
		honeyswapTicker()
	case "sushiswapTicker":
		sushiswapTicker()
	case "bitforexTicker":
		bitforexTicker()
	case "bitmaxTicker":
		bitmaxTicker()
	case "sovTicker":
		sovTicker()
	case "ethGas":
		ethGas()
	default:
		usage()
	}
}
