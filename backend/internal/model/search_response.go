package model

import "time"

type SearchResponse struct {
	Messages   []Message
	Duration   time.Duration
	TotalCount int
}
