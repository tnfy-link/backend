package links

import (
	"errors"
	"strings"
	"time"
)

type CreateLink struct {
	TargetURL string `json:"targetUrl" validate:"required,http_url"`
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

func (r *PostLinksRequest) Validate() error {
	if !strings.HasPrefix(r.Link.TargetURL, "https://") {
		return errors.New("targetUrl must start with https://")
	}

	return nil
}

type PostLinksResponse struct {
	Link Link `json:"link"`
}
