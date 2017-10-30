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

	fmt.Println("Buy all in BTC")
	if err := api.Buy(); err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println("Wait to sell all BTC")
		if err := sell(api); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Wait to buy all BTC")
		if err := buy(api); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("finished")
}

func sell(api *kraken) error {
	smallerPrice := api.lastBuy
	wait := true
	for wait {
		ticker, err := api.Ticker()
		if err != nil {
			return err
		}

		p, err := strconv.ParseFloat(ticker.Close[0], 32)
		if err != nil {
			return err
		}
		lastPrice := float32(p)
		diff := lastPrice / smallerPrice

		if diff > 1.01 {
			wait = false
		}

		if lastPrice < smallerPrice {
			smallerPrice = lastPrice
		}

		fmt.Printf("Last Price: %f Bigger Price: %f Difference: %f \n", lastPrice, smallerPrice, diff)
		time.Sleep(time.Second * 4)
	}

	if err := api.Sell(); err != nil {
		return err
	}

	return nil
}

func buy(api *kraken) error {
	biggerPrice := api.lastBuy
	hold := true
	for hold {
		ticker, err := api.Ticker()
		if err != nil {
			return err
		}

		p, err := strconv.ParseFloat(ticker.Close[0], 32)
		if err != nil {
			return err
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
		time.Sleep(time.Second * 4)
	}

	if err := api.Buy(); err != nil {
		return err
	}

	return nil
}
