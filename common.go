package main

import (
	"context"
	"log"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func usage() {
	log.Fatal(`Examples:

bincli balance
bincli balance BTC
bincli btcBalance
bincli btcBalance ETH
bincli chartBalance > index.html && open index.html
bincli ticker
bincli ticker BTCUSDT
bincli alert BTCUSDT ">" 56000 && cowsay "Reached"
bincli uniswapAlert APYS ">" 0.1 && cowsay "Reached"
bincli honeyswapAlert DAI ">" 0.1 && cowsay "Reached"
bincli bitforexAlert coin-usdt-omi ">" 0.1 && cowsay "Reached"
bincli bitmaxAlert BTC/USDT ">" 56000 && cowsay "Reached"
`)
}

type assetTicker struct {
	asset      string
	lastPrice  float64
	delta24pct float64
}

func mustRequestTicker(client *binance.Client) map[string]assetTicker {
	prices, err := client.NewListPriceChangeStatsService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	ticker := map[string]assetTicker{}
	for _, price := range prices {
		lastPrice, err := strconv.ParseFloat(price.LastPrice, 64)
		if err != nil {
			log.Fatal(err)
		}
		priceChangePercent, err := strconv.ParseFloat(price.PriceChangePercent, 64)
		if err != nil {
			log.Fatal(err)
		}
		ticker[price.Symbol] = assetTicker{asset: price.Symbol, lastPrice: lastPrice, delta24pct: priceChangePercent}
	}
	return ticker
}

func mustRequestAccountBalances(client *binance.Client) map[string]float64 {
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	balances := map[string]float64{}
	for _, assetBalance := range account.Balances {
		free, err := strconv.ParseFloat(assetBalance.Free, 64)
		if err != nil {
			log.Fatal(err)
		}
		locked, err := strconv.ParseFloat(assetBalance.Locked, 64)
		if err != nil {
			log.Fatal(err)
		}
		balances[assetBalance.Asset] = free + locked
	}
	return balances
}

type assetStatus struct {
	balance    float64
	delta24pct float64
}

func mustCalculateBalanceInBtc(asset string, balances map[string]float64, ticker map[string]assetTicker) assetStatus {
	if asset == "USDT" {
		return assetStatus{
			balance:    balances[asset] / ticker["BTCUSDT"].lastPrice,
			delta24pct: 0.0,
		}
	}
	if asset == "BTC" {
		return assetStatus{
			balance:    balances[asset],
			delta24pct: ticker["BTCUSDT"].delta24pct,
		}
	}
	if ticker[asset+"BTC"].lastPrice != 0.0 {
		return assetStatus{
			balance:    balances[asset] * ticker[asset+"BTC"].lastPrice,
			delta24pct: ticker[asset+"BTC"].delta24pct,
		}
	}
	if ticker[asset+"BTC"].lastPrice == 0 && ticker[asset+"BNB"].lastPrice != 0 {
		btcPct := ticker["BNBBTC"].delta24pct/100.0 + 1
		bnbPct := ticker[asset+"BNB"].delta24pct/100.0 + 1
		pct := (btcPct*bnbPct - 1) * 100
		return assetStatus{
			balance:    balances[asset] * ticker[asset+"BNB"].lastPrice * ticker["BNBBTC"].lastPrice,
			delta24pct: pct,
		}
	}
	if ticker[asset+"BTC"].lastPrice == 0 && ticker[asset+"BNB"].lastPrice == 0 && ticker[asset+"USDT"].lastPrice != 0 {
		return assetStatus{
			balance: balances[asset] * ticker[asset+"USDT"].lastPrice / ticker["BTCUSDT"].lastPrice,
			// TODO
			delta24pct: 0.0,
		}
	}
	return assetStatus{balance: 0, delta24pct: 0}
}

func mustCalculateAllBalancesInBtc(balances map[string]float64, ticker map[string]assetTicker, dontTrim bool) map[string]assetStatus {
	balancesInBtc := map[string]assetStatus{}
	for asset := range balances {
		balancesInBtc[asset] = mustCalculateBalanceInBtc(asset, balances, ticker)
		if !dontTrim && balancesInBtc[asset].balance < 0.005 {
			delete(balancesInBtc, asset)
		}
	}
	return balancesInBtc
}
