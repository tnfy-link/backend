package stats

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/tnfy-link/client-go/api"
)

const (
	keyTemplateStats = "links:%s:stats"

	fieldStatsTotal = "total"
)

type repository struct {
	redis *redis.Client
}

func (r *repository) Get(ctx context.Context, id string) (api.Stats, error) {
	fields, err := r.redis.HGetAll(ctx, fmt.Sprintf(keyTemplateStats, id)).Result()
	if errors.Is(err, redis.Nil) || len(fields) == 0 {
		return api.Stats{}, ErrNotFound
	}
	if err != nil {
		return api.Stats{}, fmt.Errorf("failed to get stats: %w", err)
	}

	stats := api.Stats{
		Labels: make(map[string]map[string]int),
		Total:  0,
	}

	const partsCount = 2
	for k, v := range fields {
		switch k {
		case fieldStatsTotal:
			stats.Total, _ = strconv.Atoi(v)
		default:
			parts := strings.Split(k, "|")
			if len(parts) != partsCount {
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

func (r *repository) Incr(ctx context.Context, link api.Link, labels Labels) error {
	keyStats := fmt.Sprintf(keyTemplateStats, link.ID)

	pipe := r.redis.TxPipeline()
	pipe.HIncrBy(ctx, keyStats, fieldStatsTotal, 1)
	for k, v := range labels {
		pipe.HIncrBy(ctx, keyStats, k+"|"+v, 1)
	}
	pipe.ExpireAt(ctx, keyStats, link.ValidUntil)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to register stats: %w", err)
	}

	return nil
}

func newRepository(redis *redis.Client) *repository {
	if redis == nil {
		panic("redis client is nil")
	}

	return &repository{
		redis: redis,
	}
}
