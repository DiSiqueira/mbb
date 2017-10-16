package main

import (
	"fmt"
	"github.com/beldur/kraken-go-api-client"
	"strconv"
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
	}
)

func NewKraken(key, secret string) *kraken {
	return &kraken{
		api:    krakenapi.New(key, secret),
		values: map[string]float32{},
	}
}

func (b *kraken) Balance() (float32, float32, error) {
	balance, err := b.api.Balance()
	if err != nil {
		return 0, 0, err
	}

	b.values["EUR"] = balance.ZEUR
	b.values["XBT"] = balance.ZEUR

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
	_, _, err = b.Balance()
	return err
}

func (b *kraken) order(direction string, volume float32) error {

	fmt.Sprintf("Volume: %.8f \n", volume)

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
	return b.sell(b.values["XBT"])
}

func (b *kraken) Buy() error {
	price, err := strconv.ParseFloat(b.lastTicker.Close[0], 10)
	if err != nil {
		return err
	}

	return b.buy(b.values["EUR"] / float32(price*1.001))
}

func (b *kraken) Ticker() (*krakenapi.PairTickerInfo, error) {
	ticker, err := b.api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		return nil, err
	}

	b.lastTicker = &ticker.XXBTZEUR

	return &ticker.XXBTZEUR, nil
}
