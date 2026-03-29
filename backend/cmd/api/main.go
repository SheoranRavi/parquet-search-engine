package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SheoranRavi/parquet-search-engine/internal/handlers"
	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/middleware"
	"github.com/SheoranRavi/parquet-search-engine/internal/server"
	"github.com/SheoranRavi/parquet-search-engine/internal/services"
	"github.com/SheoranRavi/parquet-search-engine/internal/store"
)

func main() {
	if err := logger.Initialize(); err != nil {
		log.Fatal("Failed to initialize logger")
	}
	defer logger.Close()

	store := store.NewInMemoryStore()
	indexer := services.NewIndexer(store)
	queryEngine := services.NewQueryEngine(store)
	searchHandler := handlers.NewSearchHandler(indexer, queryEngine)

	loggingMiddleware := middleware.Logging()
	router := server.NewRouter(searchHandler, loggingMiddleware)

	// start indexing
	pathStr, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(pathStr)
	pFilesPath := filepath.Join(path, "..", "parquet_files")
	files, err := os.ReadDir(pFilesPath)
	if err != nil {
		panic(err)
	}
	go indexer.Index(pFilesPath, files)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
