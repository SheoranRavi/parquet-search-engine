package services

import (
	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/store"
	"github.com/rs/zerolog"
)

type QueryEngine struct {
	logger zerolog.Logger
	store  *store.InMemoryStore
}

func NewQueryEngine(store *store.InMemoryStore) *QueryEngine {
	return &QueryEngine{
		logger: logger.NewServiceLogger("QueryEngine"),
		store:  store,
	}
}
