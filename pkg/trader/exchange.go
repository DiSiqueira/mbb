package trader

import (
	"errors"
	"fmt"
	"github.com/beldur/kraken-go-api-client"
	"sync"
)

var exchangeMu sync.RWMutex
var exchanges = make(map[string]Exchange)

type (
	Exchange interface {
		Start(key, secret string) error
		Buyer
		Seller
		Balance
		Market
	}

	Balance interface {
		Balance() (float32, float32, error)
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

	Config interface {
		Key() string
		Secret() string
		Exchange() string
	}
)

func NewExchange(configs Config) (Exchange, error) {
	exchange, ok := exchanges[configs.Exchange()]
	if !ok {
		return nil, fmt.Errorf("exchange %s not registred", configs.Exchange())
	}
	return exchange, exchange.Start(configs.Key(), configs.Secret())
}

func Register(name string, exchange Exchange) error {
	exchangeMu.Lock()
	defer exchangeMu.Unlock()
	if exchange == nil {
		return errors.New("register exchange is nil")
	}
	if _, dup := exchanges[name]; dup {
		return fmt.Errorf("exchange %s already registred", name)
	}
	exchanges[name] = exchange
	return nil
}
