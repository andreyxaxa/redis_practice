package main

import (
	"log"
	"os"
	"red_test/internal/app/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(os.Getenv("POSTGRES_ADDR"), os.Getenv("REDIS_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
}
