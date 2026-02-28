package main

import (
	"log"

	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/internal/blkhell"
)

func main() {
	cfg := config.Init()
	cli := blkhell.NewCli(cfg)

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
