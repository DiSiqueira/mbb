package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beldur/kraken-go-api-client"
	"time"
)

type (
	Balance interface {
		Balance() (float32, error)
	}

	Buyer interface {
		Buy() error
	}

	Seller interface {
		Sell() error
	}

	Market interface {
		Ticker() (*krakenapi.PairTickerInfo, error)
	}

	kraken struct {
		api        *krakenapi.KrakenApi
		values     map[string]float32
		lastTicker *krakenapi.PairTickerInfo
		lastBuy    float32
	}
)

func NewKraken(key, secret string) *kraken {
	k := &kraken{
		api:    krakenapi.New(key, secret),
		values: map[string]float32{},
	}

	k.Balance()
	k.Ticker()

	return k
}

func (b *kraken) Balance() (float32, float32, error) {
	balance, err := b.api.Balance()
	if err != nil {
		return 0, 0, err
	}

	b.values["EUR"] = balance.ZEUR
	b.values["XBT"] = balance.XXBT

	fmt.Printf("XBT: %.8f\n", balance.XXBT)
	fmt.Printf("EUR: %.8f\n", balance.ZEUR)

	return balance.ZEUR, balance.XXBT, nil
}

func (b *kraken) sell(volume float32) error {
	return b.forceOrder("sell", volume)
}

func (b *kraken) buy(volume float32) error {
	return b.forceOrder("buy", volume)
}

func (b *kraken) forceOrder(direction string, volume float32) error {
	err := fmt.Errorf("teste")
	for err != nil {
		err = b.order(direction, volume)
	}

	eurStart, btcStart := b.values["EUR"], b.values["XBT"]

	eur, btc, err := b.Balance()
	if err != nil {
		return err
	}

	for eur == eurStart && btc == btcStart {
		time.Sleep(5 * time.Second)
		eur, btc, err = b.Balance()
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *kraken) order(direction string, volume float32) error {
	fmt.Printf("ORDER: %s - %.8f \n", strings.ToUpper(direction), volume)

	_, err := b.api.AddOrder(
		"XXBTZEUR",
		direction,
		"market",
		fmt.Sprintf("%.8f", volume),
		map[string]string{
			"trading_agreement": "agree",
		},
	)
	if err != nil {
		if err.Error() == "Could not execute request! #7 ([EOrder:Insufficient funds])" {
			return nil
		}
		return err
	}

	return nil
}

func (b *kraken) Sell() error {
	fmt.Printf("Selling: %f \n", b.values["XBT"])
	if b.values["XBT"] <= 0.0001 {
		return nil
	}

	return b.sell(b.values["XBT"])
}

func (b *kraken) Buy() error {
	fmt.Printf("Buying: %f \n", b.values["EUR"])
	p, err := strconv.ParseFloat(b.lastTicker.Close[0], 32)
	if err != nil {
		return err
	}
	price := float32(p)

	btc := b.values["EUR"] / price * 1.005
	if btc <= 0.0001 {
		return nil
	}

	b.lastBuy = price

	return b.buy(btc)
}

func (b *kraken) ticker() *krakenapi.PairTickerInfo {
	ticker, err := b.api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		return b.ticker()
	}

	return &ticker.XXBTZEUR
}

func (b *kraken) Ticker() (*krakenapi.PairTickerInfo, error) {
	b.lastTicker = b.ticker()

	return b.lastTicker, nil
}
