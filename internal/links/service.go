package links

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/itchyny/base58-go"
)

type service struct {
	links *repository

	ttl time.Duration
}

func (s *service) Create(ctx context.Context, target CreateLink) (Link, error) {
	id, err := s.newID()
	if err != nil {
		return Link{}, fmt.Errorf("failed to generate id: %w", err)
	}

	link := Link{
		ID:         id,
		TargetURL:  target.TargetURL,
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().Add(s.ttl),
	}

	return link, s.links.Create(ctx, link)
}

func (s *service) GetTarget(ctx context.Context, id string) (string, error) {
	return s.links.GetTarget(ctx, id)
}

func (s *service) newID() (string, error) {
	var randomValue uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &randomValue)
	if err != nil {
		return "", fmt.Errorf("failed to read random value: %w", err)
	}

	id := base58.FlickrEncoding.EncodeUint64(uint64(randomValue))

	return string(id), nil
}

func newService(links *repository, config Config) *service {
	return &service{
		links: links,

		ttl: config.TTL,
	}
}
