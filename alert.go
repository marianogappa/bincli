package main

import (
	"log"
	"os"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

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
