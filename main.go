package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/repository"
	"github.com/anugrahsputra/go-quran-api/router"
	"github.com/anugrahsputra/go-quran-api/service"
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
	reindex := flag.Bool("reindex", false, "Re-index the Quran data")
	flag.Parse()

	cfg := config.LoadConfig()

	if *reindex {
		fmt.Println("Indexing Quran data...")
		quranRepo := repository.NewQuranRepository(cfg)
		searchRepo, err := repository.NewSearchRepository()
		if err != nil {
			log.Fatalf("failed to create search repository: %v", err)
		}
		searchService := service.NewSearchService(quranRepo, searchRepo)
		if err := searchService.IndexQuran(); err != nil {
			log.Fatalf("failed to index quran data: %v", err)
		}
		fmt.Println("Indexing complete.")
		return
	}

	r := router.SetupRoute()

	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
