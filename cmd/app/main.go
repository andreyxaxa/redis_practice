package main

import (
	"log"
	"os"
	"red_test/internal/app/server"
	redisstorage "red_test/internal/app/storage/redis_storage"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	storage := redisstorage.New(os.Getenv("REDIS_ADDR"))

	srv := server.New(storage)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
}
