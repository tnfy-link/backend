package links

import "time"

type CreateLink struct {
	TargetURL string `json:"targetUrl" format:"uri"`
}

type Link struct {
	ID        string `json:"id"`
	TargetURL string `json:"targetUrl"`
	URL       string `json:"url"`

	CreatedAt time.Time `json:"createdAt"`
}

type PostLinksRequest struct {
	Link CreateLink `json:"link"`
}

type PostLinksResponse struct {
	Link Link `json:"link"`
}
