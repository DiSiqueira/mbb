package main

import (
	"fmt"
	"github.com/disiqueira/mbb/pkg/config"
	"github.com/disiqueira/mbb/pkg/trader"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

func main() {
	configs := config.NewSpecification()
	if err := envconfig.Process("mbb", configs); err != nil {
		panic(err.Error())
	}

	api, err := trader.NewExchange(configs)
	if err != nil {
		panic(err.Error())
	}

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

func sell(api trader.Exchange) error {
	smallerPrice, err := api.Ticker()
	if err != nil {
		return err
	}
	wait := true
	for wait {
		lastPrice, err := api.Ticker()
		if err != nil {
			return err
		}

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

func buy(api trader.Exchange) error {
	biggerPrice, err := api.Ticker()
	if err != nil {
		return err
	}
	hold := true
	for hold {
		lastPrice, err := api.Ticker()
		if err != nil {
			return err
		}

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
