package links

import (
	"fmt"
	"time"

	"github.com/tnfy-link/client-go/api"
)

func newLink(id string, values map[string]string) (api.Link, error) {
	createdAt, err := time.Parse(time.RFC3339, values[fieldCreatedAt])
	if err != nil {
		return api.Link{}, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	validUntil, err := time.Parse(time.RFC3339, values[fieldValidUntil])
	if err != nil {
		return api.Link{}, fmt.Errorf("failed to parse validUntil: %w", err)
	}

	return api.Link{
		ID:         id,
		TargetURL:  values[fieldTargetUrl],
		CreatedAt:  createdAt,
		ValidUntil: validUntil,
	}, nil
}
