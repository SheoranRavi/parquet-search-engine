package store

import (
	"strings"
	"sync"

	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/model"
	"github.com/rs/zerolog"
)

type InMemoryStore struct {
	logger    zerolog.Logger
	messages  map[string]*model.Message // msgId -> message
	termIndex map[string][]string       // token -> []msgId
	//categoryIndex map[string]map[string][]string // category -> categoryInstance -> []msgId
	muMsg sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		messages:  make(map[string]*model.Message),
		termIndex: make(map[string][]string),
		logger:    logger.NewRepoLogger("InMemoryStore"),
	}
}

func (store *InMemoryStore) AddChunk(msgs []model.Message, termIndex map[string][]string) {
	store.muMsg.Lock()
	defer store.muMsg.Unlock()
	// add
	for _, msg := range msgs {
		if len(strings.TrimSpace(msg.MsgId)) == 0 {
			panic("Can't have empty msg id")
		}
		store.messages[msg.MsgId] = &msg
	}

	for word, indices := range termIndex {
		if _, ok := store.termIndex[word]; ok {
			store.termIndex[word] = append(store.termIndex[word], indices...)
		} else {
			store.termIndex[word] = indices
		}
	}
}

func (store *InMemoryStore) GetUnion(tokens []string) (map[string]model.Message, error) {
	store.muMsg.Lock()
	defer store.muMsg.Unlock()
	result := make(map[string]model.Message)
	for _, t := range tokens {
		mIds, ok := store.termIndex[t]
		if ok {
			for _, id := range mIds {
				result[id] = *store.messages[id]
			}
		}
	}
	return result, nil
}
