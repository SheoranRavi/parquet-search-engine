package services

import (
	"slices"
	"time"

	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/model"
	"github.com/SheoranRavi/parquet-search-engine/internal/store"
	"github.com/SheoranRavi/parquet-search-engine/internal/util"
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

func (q *QueryEngine) Query(input string) ([]model.Message, time.Duration) {
	t := time.Now()
	// tokenize the query
	tokens := util.Tokenize(input)
	tokens = util.FilterStopWords(tokens)
	messages, _ := q.store.Get(tokens)
	// order messages by timestamp
	slices.SortFunc(messages, func(a, b model.Message) int {
		if a.NanoTimeStamp > b.NanoTimeStamp {
			return 1
		} else if a.NanoTimeStamp < b.NanoTimeStamp {
			return -1
		}
		return 0
	})
	elapsed := time.Since(t)
	q.logger.Info().Msgf("Fetched %d messages in %d", len(messages), elapsed)
	return messages, elapsed
}
