package services

import (
	"os"
	"path/filepath"
	"strings"
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
}

func NewIndexer(store *store.InMemoryStore) *Indexer {
	return &Indexer{
		logger: logger.NewServiceLogger("Indexer"),
		store:  store,
	}
}

// index a bunch of files
func (indexer *Indexer) Index(parentDir string, files []os.DirEntry) {
	tStart := time.Now()
	indexer.logger.Info().Msg("Starting the indexing")
	for _, file := range files {
		filePath := filepath.Join(parentDir, file.Name())
		_, err := indexer.IndexFile(filePath)
		if err != nil {
			indexer.logger.Error().Msgf("Error indexing file: %s, error: %s", file, err)
		}
	}
	tElapsed := time.Since(tStart)
	indexer.logger.Info().Msgf("Time taken for indexing all files: %d ms", tElapsed.Milliseconds())
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
	indexer.store.AddChunk(rows, termIndex)
	return time.Since(tStart), nil
}
