package handlers

import (
	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/services"
	"github.com/rs/zerolog"
)

type SearchHandler struct {
	indexer     *services.Indexer
	queryEngine *services.QueryEngine
	logger      zerolog.Logger
}

func NewSearchHandler(indexer *services.Indexer, queryEngine *services.QueryEngine) *SearchHandler {
	return &SearchHandler{
		indexer:     indexer,
		queryEngine: queryEngine,
		logger:      logger.NewHandlerLogger("SearchHandler"),
	}
}
