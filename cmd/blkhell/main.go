package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vague2k/blkhell/internal/blkhell"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	if err := blkhell.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
