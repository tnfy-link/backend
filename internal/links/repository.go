package links

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	keyIndex         = "links:index"
	keyTemplateMeta  = "links:%s:meta"
	keyTemplateStats = "links:%s:stats"

	fieldTargetUrl  = "targetUrl"
	fieldCreatedAt  = "createdAt"
	fieldValidUntil = "validUntil"
)

type repository struct {
	redis *redis.Client
}

func (r *repository) GetTarget(ctx context.Context, id string) (string, error) {
	return r.redis.HGet(ctx, fmt.Sprintf(keyTemplateMeta, id), fieldTargetUrl).Result()
}

func (r *repository) Create(ctx context.Context, link Link) error {
	canCreate, err := r.redis.HSetNX(ctx, keyIndex, link.ID, link.TargetURL).Result()
	if err != nil {
		return fmt.Errorf("failed to set link: %w", err)
	}
	if !canCreate {
		return fmt.Errorf("failed to set link: %w", ErrLinkAlreadyExists)
	}

	keyMeta := fmt.Sprintf(keyTemplateMeta, link.ID)
	pipe := r.redis.Pipeline()
	pipe.HSet(ctx, keyMeta, map[string]string{
		fieldTargetUrl:  link.TargetURL,
		fieldCreatedAt:  link.CreatedAt.Format(time.RFC3339),
		fieldValidUntil: link.ValidUntil.Format(time.RFC3339),
	})
	pipe.ExpireAt(ctx, keyMeta, link.ValidUntil)
	pipe.HExpireAt(ctx, keyIndex, link.ValidUntil, link.ID)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set link: %w", err)
	}

	return nil
}

func newRepository(redis *redis.Client) *repository {
	return &repository{
		redis: redis,
	}
}
