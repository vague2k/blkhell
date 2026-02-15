package main

import (
	"log"

	"github.com/vague2k/blkhell/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
