package links

import (
	"time"
)

type CreateLink struct {
	TargetURL string `json:"targetUrl" validate:"required"`
}

type Link struct {
	ID        string `json:"id"`
	TargetURL string `json:"targetUrl"`
	URL       string `json:"url"`

	CreatedAt  time.Time `json:"createdAt"`
	ValidUntil time.Time `json:"validUntil"`
}

type PostLinksRequest struct {
	Link CreateLink `json:"link"`
}

type PostLinksResponse struct {
	Link Link `json:"link"`
}

// Stats
type Stats struct {
	Labels map[string]map[string]int `json:"labels"`
	Total  int                       `json:"total"`
}

type GetStatsResponse struct {
	Stats Stats `json:"stats"`
}
