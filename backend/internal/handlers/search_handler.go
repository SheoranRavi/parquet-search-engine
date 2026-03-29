package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/model"
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

func (sh *SearchHandler) Search(rw http.ResponseWriter, req *http.Request) {
	var query model.SearchRequest
	if err := json.NewDecoder(req.Body).Decode(&query); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	msgs, duration := sh.queryEngine.Query(query.Query)
	resp := &model.SearchResponse{
		Messages:   msgs,
		Duration:   duration,
		TotalCount: len(msgs),
	}
	_ = json.NewEncoder(rw).Encode(resp)
}

// func (sh *SearchHandler) Upload(rw http.ResponseWriter, req *http.Request) {

// }
