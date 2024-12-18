package links

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/tnfy-link/backend/internal/id"
	"github.com/tnfy-link/client-go/api"
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

type Service struct {
	idgen *id.Generator

	links  *repository
	logger *zap.Logger

	hostname string
	ttl      time.Duration
}

func (s *Service) Create(ctx context.Context, target api.CreateLink) (api.Link, error) {
	if target.TargetURL == "" {
		return api.Link{}, newValidationError("targetUrl", fmt.Errorf("value is empty"))
	}
	if len(target.TargetURL) > maxTargetURLLength {
		return api.Link{}, newValidationError("targetUrl", fmt.Errorf("value too long"))
	}

	parsedUrl, err := url.Parse(target.TargetURL)
	if err != nil {
		return api.Link{}, newValidationError("targetUrl", fmt.Errorf("invalid url: %w", err))
	}
	if parsedUrl.Scheme != "https" {
		return api.Link{}, newValidationError("targetUrl", fmt.Errorf("scheme must be https"))
	}

	id, err := s.idgen.New()
	if err != nil {
		return api.Link{}, fmt.Errorf("failed to generate id: %w", err)
	}

	link := api.Link{
		ID:         id,
		URL:        s.hostname + "/" + id,
		TargetURL:  target.TargetURL,
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().Add(s.ttl),
	}

	return link, s.links.Create(ctx, link)
}

func (s *Service) Get(ctx context.Context, id string) (api.Link, error) {
	if err := s.idgen.Validate(id); err != nil {
		return api.Link{}, newValidationError("id", err)
	}

	link, err := s.links.Get(ctx, id)
	if err != nil {
		return api.Link{}, err
	}

	link.URL = s.hostname + "/" + link.ID

	return link, nil
}

func (s *Service) GetTarget(ctx context.Context, id string) (string, error) {
	if err := s.idgen.Validate(id); err != nil {
		return "", newValidationError("id", err)
	}

	return s.links.GetTarget(ctx, id)
}

func (s *Service) RegisterStats(ctx context.Context, id string, query string) error {
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

func (s *Service) GetStats(ctx context.Context, id string) (api.Stats, error) {
	if err := s.idgen.Validate(id); err != nil {
		return api.Stats{}, newValidationError("id", err)
	}

	return s.links.GetStats(ctx, id)
}

func (s *Service) ValidateID(id string) error {
	return s.idgen.Validate(id)
}

func NewService(
	idgen *id.Generator,
	links *repository,
	logger *zap.Logger,
	config Config,
) *Service {
	if idgen == nil {
		panic("id generator is required")
	}
	if links == nil {
		panic("links repository is required")
	}
	if logger == nil {
		panic("logger is required")
	}

	return &Service{
		idgen: idgen,

		links:  links,
		logger: logger,

		hostname: config.Hostname,
		ttl:      config.TTL,
	}
}
