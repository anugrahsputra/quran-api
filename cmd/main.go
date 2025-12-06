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
	// Try to load .env file, but don't fail if it doesn't exist
	// In Docker, environment variables are set directly, so .env is optional
	err := godotenv.Load()
	if err != nil {
		// Only log in development mode, in production this is expected
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

	if *reindex {
		fmt.Println("Indexing Quran data...")
		quranRepo := repository.NewQuranRepository(cfg)
		searchRepo, err := repository.NewQuranSearchRepository(cfg.SearchIndexPath)
		if err != nil {
			log.Fatalf("failed to create search repository: %v", err)
		}
		searchService := service.NewQuranSearchService(quranRepo, searchRepo)
		if err := searchService.IndexQuran(); err != nil {
			log.Fatalf("failed to index quran data: %v", err)
		}
		fmt.Println("Indexing complete.")
		return
	}

	r := router.SetupRoute()

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
		// Production-ready timeouts
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
