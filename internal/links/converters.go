package links

import (
	"fmt"
	"time"
)

func newLink(id string, values map[string]string) (Link, error) {
	createdAt, err := time.Parse(time.RFC3339, values[fieldCreatedAt])
	if err != nil {
		return Link{}, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	validUntil, err := time.Parse(time.RFC3339, values[fieldValidUntil])
	if err != nil {
		return Link{}, fmt.Errorf("failed to parse validUntil: %w", err)
	}

	return Link{
		ID:         id,
		TargetURL:  values[fieldTargetUrl],
		CreatedAt:  createdAt,
		ValidUntil: validUntil,
	}, nil
}
