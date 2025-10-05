package main

import (
	"fmt"

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

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
