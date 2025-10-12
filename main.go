package main

import (
	"fmt"
	"log"

	"github.com/anugrahsputra/quran-api/config"
	"github.com/anugrahsputra/quran-api/router"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	config.ConfigureLogger()
}

func main() {
	cfg := config.LoadConfig()
	router := router.SetupRoute()

	if err := router.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
