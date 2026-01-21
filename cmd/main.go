package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anugrahsputra/go-quran-api/config"
	"github.com/anugrahsputra/go-quran-api/repository"
	"github.com/anugrahsputra/go-quran-api/router"
	"github.com/anugrahsputra/go-quran-api/service"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		if os.Getenv("ENV") != "production" && os.Getenv("ENV") != "prod" {
			fmt.Println("Note: .env file not found (this is OK if using environment variables)")
		}
	}

	config.ConfigureLogger()
}

func main() {
	reindex := flag.Bool("reindex", false, "Re-index the Quran data")
	flag.Parse()

	cfg := config.LoadConfig()

	quranRepo := repository.NewQuranRepository(cfg)
	searchRepo, err := repository.NewQuranSearchRepository(cfg.SearchIndexPath)
	if err != nil {
		log.Fatalf("failed to create search repository: %v", err)
	}
	searchService := service.NewQuranSearchService(quranRepo, searchRepo)

	if *reindex {
		fmt.Println("Indexing Quran data...")
		if err := searchService.IndexQuran(); err != nil {
			log.Fatalf("failed to index quran data: %v", err)
		}
		fmt.Println("Indexing complete.")
		return
	}

	count, err := searchRepo.GetDocCount()
	if err == nil && count == 0 {
		if os.Getenv("AUTO_INDEX") == "true" {
			log.Println("Search index is empty. Starting automatic indexing in background...")
			go func() {
				if err := searchService.IndexQuran(); err != nil {
					log.Printf("Automatic indexing failed: %v", err)
				} else {
					log.Println("Automatic indexing complete.")
				}
			}()
		} else {
			log.Println("Warning: Search index is empty. Search functionality will not return results.")
			log.Println("To populate the index, run with -reindex flag or set AUTO_INDEX=true environment variable.")
		}
	}

	r := router.SetupRoute(cfg, quranRepo, searchRepo, searchService)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}