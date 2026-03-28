package store

import (
	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/rs/zerolog"
)

type InMemoryStore struct {
	logger zerolog.Logger
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		logger: logger.NewRepoLogger("InMemoryStore"),
	}
}
