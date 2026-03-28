package services

import (
	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/store"
	"github.com/rs/zerolog"
)

type Indexer struct {
	logger zerolog.Logger
	store  *store.InMemoryStore
}

func NewIndexer(store *store.InMemoryStore) *Indexer {
	return &Indexer{
		logger: logger.NewServiceLogger("Indexer"),
		store:  store,
	}
}

func (indexer *Indexer) Start() error {

}

func (indexer *Indexer) IndexFile(filePath string) error {

}
