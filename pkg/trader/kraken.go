package trader

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beldur/kraken-go-api-client"
)

type (
	kraken struct {
		api        *krakenapi.KrakenApi
		values     map[string]float32
		lastTicker *krakenapi.PairTickerInfo
		lastBuy    float32
	}
)

func init() {
	k := kraken{}
	Register("kraken", &k)
}

func (k *kraken) Start(key, secret string) error {
	k.api = krakenapi.New(key, secret)
	k.values = map[string]float32{}

	if _, _, err := k.Balance(); err != nil {
		return err
	}
	if _, err := k.Ticker(); err != nil {
		return err
	}
	return nil
}

func (k *kraken) Balance() (float32, float32, error) {
	balance, err := k.api.Balance()
	if err != nil {
		return 0, 0, err
	}

	k.values["EUR"] = balance.ZEUR
	k.values["XBT"] = balance.XXBT

	fmt.Printf("XBT: %.8f\n", balance.XXBT)
	fmt.Printf("EUR: %.8f\n", balance.ZEUR)

	return balance.ZEUR, balance.XXBT, nil
}

func (k *kraken) sell(volume float32) error {
	return k.forceOrder("sell", volume)
}

func (k *kraken) buy(volume float32) error {
	return k.forceOrder("buy", volume)
}

func (k *kraken) forceOrder(direction string, volume float32) error {
	err := fmt.Errorf("teste")
	for err != nil {
		err = k.order(direction, volume)
	}

	eurStart, btcStart := k.values["EUR"], k.values["XBT"]

	eur, btc, err := k.Balance()
	if err != nil {
		return err
	}

	for eur == eurStart && btc == btcStart {
		time.Sleep(5 * time.Second)
		eur, btc, err = k.Balance()
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *kraken) order(direction string, volume float32) error {
	fmt.Printf("ORDER: %s - %.8f \n", strings.ToUpper(direction), volume)

	_, err := k.api.AddOrder(
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

func (k *kraken) Sell() error {
	fmt.Printf("Selling: %f \n", k.values["XBT"])
	if k.values["XBT"] <= 0.0001 {
		return nil
	}

	return k.sell(k.values["XBT"])
}

func (k *kraken) Buy() error {
	fmt.Printf("Buying: %f \n", k.values["EUR"])
	p, err := strconv.ParseFloat(k.lastTicker.Close[0], 32)
	if err != nil {
		return err
	}
	price := float32(p)

	btc := k.values["EUR"] / price * 1.005
	if btc <= 0.0001 {
		return nil
	}

	k.lastBuy = price

	return k.buy(btc)
}

func (k *kraken) ticker() *krakenapi.PairTickerInfo {
	ticker, err := k.api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		return k.ticker()
	}

	return &ticker.XXBTZEUR
}

func (k *kraken) Ticker() (*krakenapi.PairTickerInfo, error) {
	k.lastTicker = k.ticker()

	return k.lastTicker, nil
}
