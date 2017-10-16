package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	key := os.Getenv("KRAKEN_KEY")
	secret := os.Getenv("KRAKEN_SECRET")

	api := NewKraken(key, secret)

	if err := start(api); err != nil {
		log.Fatal(err)
	}

	biggerPrice := api.lastBuy
	hold := true
	for hold {
		ticker, err := api.Ticker()
		if err != nil {
			log.Fatal(err)
		}

		p, err := strconv.ParseFloat(ticker.Close[0], 32)
		if err != nil {
			log.Fatal(err)
		}
		lastPrice := float32(p)
		diff := lastPrice / biggerPrice

		if diff < .99 {
			hold = false
		}

		if lastPrice > biggerPrice {
			biggerPrice = lastPrice
		}

		fmt.Printf("Last Price: %f Bigger Price: %f Difference: %f \n", lastPrice, biggerPrice, diff)
		time.Sleep(time.Second * 25)
	}

	if err := api.Sell(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("finished")
}

func start(api *kraken) error {
	if err := api.Sell(); err != nil {
		return err
	}

	if err := api.Buy(); err != nil {
		return err
	}

	return nil
}
