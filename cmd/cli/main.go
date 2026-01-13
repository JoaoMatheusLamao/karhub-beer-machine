package main

import (
	"log"

	"karhub-beer-machine/cmd/cli/root"
)

func main() {
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
