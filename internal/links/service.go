package links

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/url"
	"time"

	"github.com/itchyny/base58-go"
	"go.uber.org/zap"
)

const (
	maxLabelValueLength = 64
)

var utmLabels = map[string]string{
	"utm_source":   "source",
	"utm_medium":   "medium",
	"utm_campaign": "campaign",
}

type service struct {
	links  *repository
	logger *zap.Logger

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

func (s *service) RegisterStats(ctx context.Context, id string, query string) error {
	values, err := url.ParseQuery(query)
	if err != nil {
		s.logger.Error("failed to parse query", zap.Error(err))
	}

	labels := Labels{}

	for k, v := range utmLabels {
		if val := values.Get(k); v != "" {
			if len(val) > maxLabelValueLength {
				s.logger.Warn("label value too long", zap.String("id", id), zap.String("label", v), zap.String("value", val))
				val = val[:maxLabelValueLength]
			}
			labels[v] = val
		}
	}

	return s.links.RegisterStats(ctx, id, labels)
}

func (s *service) GetStats(ctx context.Context, id string) (Stats, error) {
	return s.links.GetStats(ctx, id)
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

func newService(links *repository, logger *zap.Logger, config Config) *service {
	return &service{
		links:  links,
		logger: logger,

		ttl: config.TTL,
	}
}
