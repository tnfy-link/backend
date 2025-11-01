package links

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tnfy-link/client-go/api"
)

const (
	keyIndex        = "links:index"
	keyTemplateMeta = "links:%s:meta"

	fieldTargetURL  = "targetUrl"
	fieldCreatedAt  = "createdAt"
	fieldValidUntil = "validUntil"
)

type repository struct {
	redis *redis.Client
}

func (r *repository) Get(ctx context.Context, id string) (api.Link, error) {
	keyMeta := fmt.Sprintf(keyTemplateMeta, id)

	value, err := r.redis.HGetAll(ctx, keyMeta).Result()
	if errors.Is(err, redis.Nil) || len(value) == 0 {
		return api.Link{}, ErrLinkNotFound
	} else if err != nil {
		return api.Link{}, fmt.Errorf("failed to get link: %w", err)
	}
	return newLink(id, value)
}

func (r *repository) GetTarget(ctx context.Context, id string) (string, error) {
	value, err := r.redis.HGet(ctx, fmt.Sprintf(keyTemplateMeta, id), fieldTargetURL).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrLinkNotFound
	} else if err != nil {
		return "", fmt.Errorf("failed to get target: %w", err)
	}
	return value, nil
}

func (r *repository) Create(ctx context.Context, link api.Link) error {
	canCreate, err := r.redis.HSetNX(ctx, keyIndex, link.ID, link.TargetURL).Result()
	if err != nil {
		return fmt.Errorf("failed to set link: %w", err)
	}
	if !canCreate {
		return fmt.Errorf("failed to set link: %w", ErrLinkAlreadyExists)
	}

	keyMeta := fmt.Sprintf(keyTemplateMeta, link.ID)
	pipe := r.redis.TxPipeline()
	pipe.HSet(ctx, keyMeta, map[string]string{
		fieldTargetURL:  link.TargetURL,
		fieldCreatedAt:  link.CreatedAt.Format(time.RFC3339),
		fieldValidUntil: link.ValidUntil.Format(time.RFC3339),
	})
	pipe.ExpireAt(ctx, keyMeta, link.ValidUntil)
	pipe.HExpireAt(ctx, keyIndex, link.ValidUntil, link.ID)

	if _, execErr := pipe.Exec(ctx); execErr != nil {
		return fmt.Errorf("failed to set link: %w", execErr)
	}

	return nil
}

func newRepository(redis *redis.Client) *repository {
	return &repository{
		redis: redis,
	}
}
