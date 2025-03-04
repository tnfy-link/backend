package links

import (
	"context"
	"fmt"
	"time"

	"github.com/tnfy-link/backend/internal/id"
	"github.com/tnfy-link/client-go/api"
	"go.uber.org/zap"
)

type Service struct {
	idgen *id.Generator

	links  *repository
	logger *zap.Logger

	hostname string
	ttl      time.Duration
}

func (s *Service) Create(ctx context.Context, target NewLink) (api.Link, error) {
	if err := target.Validate(); err != nil {
		return api.Link{}, newValidationError("link", err)
	}

	id, err := s.idgen.New(ctx)
	if err != nil {
		return api.Link{}, fmt.Errorf("failed to generate id: %w", err)
	}

	link := api.Link{
		ID:         id,
		URL:        s.hostname + "/" + id,
		TargetURL:  target.TargetURL(),
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
