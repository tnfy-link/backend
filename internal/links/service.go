package links

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/tnfy-link/backend/internal/id"
	"go.uber.org/zap"
)

const (
	maxUTMLabelValueLength = 64
	maxTargetURLLength     = 2048
)

var utmLabels = map[string]string{
	"utm_source":   "source",
	"utm_medium":   "medium",
	"utm_campaign": "campaign",
}

type service struct {
	idgen *id.Generator

	links  *repository
	logger *zap.Logger

	ttl time.Duration
}

func (s *service) Create(ctx context.Context, target CreateLink) (Link, error) {
	if target.TargetURL == "" {
		return Link{}, newValidationError("targetUrl", fmt.Errorf("value is empty"))
	}
	if len(target.TargetURL) > maxTargetURLLength {
		return Link{}, newValidationError("targetUrl", fmt.Errorf("value too long"))
	}

	parsedUrl, err := url.Parse(target.TargetURL)
	if err != nil {
		return Link{}, newValidationError("targetUrl", fmt.Errorf("invalid url: %w", err))
	}
	if parsedUrl.Scheme != "https" {
		return Link{}, newValidationError("targetUrl", fmt.Errorf("scheme must be https"))
	}

	id, err := s.idgen.New()
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

func (s *service) Get(ctx context.Context, id string) (Link, error) {
	if err := s.idgen.Validate(id); err != nil {
		return Link{}, newValidationError("id", err)
	}

	return s.links.Get(ctx, id)
}

func (s *service) GetTarget(ctx context.Context, id string) (string, error) {
	if err := s.idgen.Validate(id); err != nil {
		return "", newValidationError("id", err)
	}

	return s.links.GetTarget(ctx, id)
}

func (s *service) RegisterStats(ctx context.Context, id string, query string) error {
	values, err := url.ParseQuery(query)
	if err != nil {
		// not a fatal error, just log
		s.logger.Warn("failed to parse query", zap.Error(err))
	}

	labels := Labels{}

	for k, v := range utmLabels {
		if val := values.Get(k); val != "" {
			if err := validateUTMValue(val); err != nil {
				s.logger.Warn("invalid utm value", zap.String("id", id), zap.String("label", v), zap.String("value", val), zap.Error(err))
				continue
			}

			if len(val) > maxUTMLabelValueLength {
				s.logger.Warn("label value too long", zap.String("id", id), zap.String("label", v), zap.String("value", val))
				val = val[:maxUTMLabelValueLength]
			}
			labels[v] = val
		}
	}

	return s.links.RegisterStats(ctx, id, labels)
}

func (s *service) GetStats(ctx context.Context, id string) (Stats, error) {
	if err := s.idgen.Validate(id); err != nil {
		return Stats{}, newValidationError("id", err)
	}

	return s.links.GetStats(ctx, id)
}

func (s *service) ValidateID(id string) error {
	return s.idgen.Validate(id)
}

func newService(
	idgen *id.Generator,
	links *repository,
	logger *zap.Logger,
	config Config,
) *service {
	if idgen == nil {
		panic("id generator is required")
	}
	if links == nil {
		panic("links repository is required")
	}
	if logger == nil {
		panic("logger is required")
	}

	return &service{
		idgen: idgen,

		links:  links,
		logger: logger,

		ttl: config.TTL,
	}
}
