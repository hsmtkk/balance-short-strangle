package main

import (
	"log"

	"github.com/hsmtkk/balance-short-strangle/command"
)

func main() {
	command := command.Command
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
