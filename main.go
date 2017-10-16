package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	key := os.Getenv("KRAKEN_KEY")
	secret := os.Getenv("KRAKEN_SECRET")

	api := NewKraken(key, secret)

	eur, xbt, err := api.Balance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("XBT: %.8f\n", xbt)
	fmt.Printf("EUR: %.8f\n", eur)

	ticket, err := api.Ticker()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Price: %.8f\n", ticket.Close[0])

	if xbt > 0 {
		err = api.Sell()
		if err != nil {
			log.Fatal(err)
		}
	}

	err = api.Buy()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("finished")
}
