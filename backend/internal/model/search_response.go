package model

type SearchResponse struct {
	Messages   []Message `json:"messages"`
	Duration   int64     `json:"duration"`
	TotalCount int       `json:"totalCount"`
}
