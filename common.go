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
bincli chartBalanceBtc > index.html && open index.html
bincli chartBalanceUsdt > index.html && open index.html
bincli chartBalanceDataBtc
bincli chartBalanceDataUsdt
bincli ticker
bincli ticker BTCUSDT
bincli alert BTCUSDT ">" 56000 && cowsay "Reached"
bincli uniswapAlert APYS ">" 0.1 && cowsay "Reached"
bincli honeyswapAlert DAI ">" 0.1 && cowsay "Reached"
bincli bitforexAlert coin-usdt-omi ">" 0.1 && cowsay "Reached"
bincli bitmaxAlert BTC/USDT ">" 56000 && cowsay "Reached"
binci uniswapTicker APYS
binci honeyswapTicker DAI
binci bitforexTicker coin-usdt-omi
binci bitmaxTicker BTC/USDT
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
	BalanceInBTC  float64 `json:"balanceInBTC"`
	BalanceInUSDT float64 `json:"balanceInUSDT"`
	Delta24pct    float64 `json:"delta24pct"`
}

func mustCalculateBalance(asset string, balances map[string]float64, ticker map[string]assetTicker) assetStatus {
	if asset == "USDT" {
		return assetStatus{
			BalanceInBTC:  balances[asset] / ticker["BTCUSDT"].lastPrice,
			BalanceInUSDT: balances[asset],
			Delta24pct:    0.0,
		}
	}
	if asset == "BTC" {
		return assetStatus{
			BalanceInBTC:  balances[asset],
			BalanceInUSDT: balances[asset] * ticker["BTCUSDT"].lastPrice,
			Delta24pct:    ticker["BTCUSDT"].delta24pct,
		}
	}
	if ticker[asset+"BTC"].lastPrice != 0.0 {
		balanceInBTC := balances[asset] * ticker[asset+"BTC"].lastPrice
		balanceInUSDT := balanceInBTC * ticker["BTCUSDT"].lastPrice
		if ticker[asset+"USDT"].lastPrice != 0 {
			balanceInUSDT = balances[asset] * ticker[asset+"USDT"].lastPrice
		}
		return assetStatus{
			BalanceInBTC:  balanceInBTC,
			BalanceInUSDT: balanceInUSDT,
			Delta24pct:    ticker[asset+"BTC"].delta24pct,
		}
	}
	if ticker[asset+"BTC"].lastPrice == 0 && ticker[asset+"BNB"].lastPrice != 0 {
		btcPct := ticker["BNBBTC"].delta24pct/100.0 + 1
		bnbPct := ticker[asset+"BNB"].delta24pct/100.0 + 1
		pct := (btcPct*bnbPct - 1) * 100

		balanceInBTC := balances[asset] * ticker[asset+"BNB"].lastPrice * ticker["BNBBTC"].lastPrice
		balanceInUSDT := balanceInBTC * ticker["BTCUSDT"].lastPrice
		if ticker[asset+"USDT"].lastPrice != 0 {
			balanceInUSDT = balances[asset] * ticker[asset+"USDT"].lastPrice
		}
		return assetStatus{
			BalanceInBTC:  balanceInBTC,
			BalanceInUSDT: balanceInUSDT,
			Delta24pct:    pct,
		}
	}
	if ticker[asset+"BTC"].lastPrice == 0 && ticker[asset+"BNB"].lastPrice == 0 && ticker[asset+"USDT"].lastPrice != 0 {
		balanceInBTC := balances[asset] * ticker[asset+"USDT"].lastPrice / ticker["BTCUSDT"].lastPrice
		balanceInUSDT := balanceInBTC * ticker["BTCUSDT"].lastPrice
		if ticker[asset+"USDT"].lastPrice != 0 {
			balanceInUSDT = balances[asset] * ticker[asset+"USDT"].lastPrice
		}
		return assetStatus{
			BalanceInBTC:  balanceInBTC,
			BalanceInUSDT: balanceInUSDT,
			// TODO
			Delta24pct: 0.0,
		}
	}
	return assetStatus{BalanceInBTC: 0, BalanceInUSDT: 0, Delta24pct: 0}
}

func mustCalculateAllBalances(balances map[string]float64, ticker map[string]assetTicker, dontTrim bool) map[string]assetStatus {
	calculatedBalances := map[string]assetStatus{}
	totalBTC := 0.0
	totalUSDT := 0.0
	for asset := range balances {
		calculatedBalances[asset] = mustCalculateBalance(asset, balances, ticker)
		totalBTC += calculatedBalances[asset].BalanceInBTC
		totalUSDT += calculatedBalances[asset].BalanceInUSDT
		if !dontTrim && calculatedBalances[asset].BalanceInBTC < 0.005 {
			delete(calculatedBalances, asset)
		}
	}
	calculatedBalances["Total"] = assetStatus{BalanceInBTC: totalBTC, BalanceInUSDT: totalUSDT}
	return calculatedBalances
}

func isConditionMet(price float64, comparator string, target float64) (float64, bool, error) {

	switch comparator {
	case ">":
		if price > target {
			return price, true, nil
		}
	case "<":
		if price < target {
			return price, true, nil
		}
	case ">=":
		if price >= target {
			return price, true, nil
		}
	case "<=":
		if price <= target {
			return price, true, nil
		}
	default:
		usage()

	}
	return price, false, nil
}
