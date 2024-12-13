package links

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tnfy-link/client-go/api"
)

const (
	keyIndex         = "links:index"
	keyTemplateMeta  = "links:%s:meta"
	keyTemplateStats = "links:%s:stats"

	fieldTargetUrl  = "targetUrl"
	fieldCreatedAt  = "createdAt"
	fieldValidUntil = "validUntil"

	fieldStatsTotal = "total"
)

type Labels map[string]string

type repository struct {
	redis *redis.Client
}

func (r *repository) Get(ctx context.Context, id string) (api.Link, error) {
	keyMeta := fmt.Sprintf(keyTemplateMeta, id)

	value, err := r.redis.HGetAll(ctx, keyMeta).Result()
	if err == redis.Nil || len(value) == 0 {
		return api.Link{}, ErrLinkNotFound
	} else if err != nil {
		return api.Link{}, fmt.Errorf("failed to get link: %w", err)
	}
	return newLink(id, value)
}

func (r *repository) GetTarget(ctx context.Context, id string) (string, error) {
	value, err := r.redis.HGet(ctx, fmt.Sprintf(keyTemplateMeta, id), fieldTargetUrl).Result()
	if err == redis.Nil {
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
	keyStats := fmt.Sprintf(keyTemplateStats, link.ID)
	pipe := r.redis.TxPipeline()
	pipe.HSet(ctx, keyMeta, map[string]string{
		fieldTargetUrl:  link.TargetURL,
		fieldCreatedAt:  link.CreatedAt.Format(time.RFC3339),
		fieldValidUntil: link.ValidUntil.Format(time.RFC3339),
	})
	pipe.HSet(ctx, keyStats, map[string]any{
		fieldStatsTotal: 0,
	})
	pipe.ExpireAt(ctx, keyMeta, link.ValidUntil)
	pipe.ExpireAt(ctx, keyStats, link.ValidUntil)
	pipe.HExpireAt(ctx, keyIndex, link.ValidUntil, link.ID)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set link: %w", err)
	}

	return nil
}

func (r *repository) RegisterStats(ctx context.Context, id string, labels Labels) error {
	keyStats := fmt.Sprintf(keyTemplateStats, id)
	exists, err := r.redis.Exists(ctx, keyStats).Result()
	if err != nil {
		return fmt.Errorf("failed to register stats: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("stats not found")
	}

	pipe := r.redis.TxPipeline()
	pipe.HIncrBy(ctx, keyStats, fieldStatsTotal, 1)

	for k, v := range labels {
		pipe.HIncrBy(ctx, keyStats, k+"|"+v, 1)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to register stats: %w", err)
	}

	return nil
}

func (r *repository) GetStats(ctx context.Context, id string) (api.Stats, error) {
	fields, err := r.redis.HGetAll(ctx, fmt.Sprintf(keyTemplateStats, id)).Result()
	if err == redis.Nil || len(fields) == 0 {
		return api.Stats{}, ErrLinkNotFound
	}
	if err != nil {
		return api.Stats{}, fmt.Errorf("failed to get stats: %w", err)
	}

	stats := api.Stats{
		Labels: make(map[string]map[string]int),
		Total:  0,
	}

	for k, v := range fields {
		switch k {
		case fieldStatsTotal:
			stats.Total, _ = strconv.Atoi(v)
		default:
			parts := strings.Split(k, "|")
			if len(parts) != 2 {
				continue
			}

			if stats.Labels[parts[0]] == nil {
				stats.Labels[parts[0]] = make(map[string]int)
			}

			stats.Labels[parts[0]][parts[1]], _ = strconv.Atoi(v)
		}
	}

	return stats, nil
}

func newRepository(redis *redis.Client) *repository {
	return &repository{
		redis: redis,
	}
}
