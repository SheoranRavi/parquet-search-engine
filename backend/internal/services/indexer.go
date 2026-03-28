package services

import (
	"strings"
	"sync"
	"time"

	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/model"
	"github.com/SheoranRavi/parquet-search-engine/internal/store"
	"github.com/SheoranRavi/parquet-search-engine/internal/util"
	"github.com/google/uuid"
	"github.com/parquet-go/parquet-go"
	"github.com/rs/zerolog"
)

type Indexer struct {
	logger zerolog.Logger
	store  *store.InMemoryStore
	muIdx  sync.Mutex
}

func NewIndexer(store *store.InMemoryStore) *Indexer {
	return &Indexer{
		logger: logger.NewServiceLogger("Indexer"),
		store:  store,
	}
}

// index a bunch of files
func (indexer *Indexer) Index(files []string) {
	tStart := time.Now()
	indexer.logger.Info().Msg("Starting the indexing")
	for _, filePath := range files {
		_, err := indexer.IndexFile(filePath)
		if err != nil {
			indexer.logger.Error().Msgf("Error indexing file: %s, error: %s", filePath, err)
		}
	}
	tElapsed := time.Since(tStart)
	indexer.logger.Info().Msgf("Time taken for indexing all files: %d", tElapsed)
}

// index one file
func (indexer *Indexer) IndexFile(filePath string) (time.Duration, error) {
	tStart := time.Now()
	rows, err := parquet.ReadFile[model.Message](filePath)
	if err != nil {
		return time.Since(tStart), err
	}

	termIndex := make(map[string][]string)
	for _, row := range rows {
		tokens := util.Tokenize(row.MessageRaw)
		tokens = util.FilterStopWords(tokens)
		if len(strings.TrimSpace(row.MsgId)) == 0 {
			row.MsgId = uuid.NewString()
		}
		for _, t := range tokens {
			termIndex[t] = append(termIndex[t], row.MsgId)
		}
	}
	indexer.muIdx.Lock()
	defer indexer.muIdx.Unlock()
	indexer.store.AddChunk(rows, termIndex)
	return time.Since(tStart), nil
}
