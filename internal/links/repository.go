package links

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/itchyny/base58-go"
	"github.com/redis/go-redis/v9"
)

const (
	keyTemplateMeta   = "links:%s:meta"
	keyTemplateTarget = "links:%s:target"
	keyTemplateStats  = "links:%s:stats"

	fieldTargetUrl = "targetUrl"
	fieldCreatedAt = "createdAt"
)

type repository struct {
	redis *redis.Client

	ttl time.Duration
}

func (r *repository) IsExists(ctx context.Context, id string) (bool, error) {
	count, err := r.redis.Exists(ctx, fmt.Sprintf(keyTemplateMeta, id)).Result()
	return count != 0, err
}

func (r *repository) Get(ctx context.Context, id string) (Link, error) {
	key := fmt.Sprintf(keyTemplateMeta, id)
	fields, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return Link{}, fmt.Errorf("failed to get link: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, fields["createdAt"])
	if err != nil {
		return Link{}, fmt.Errorf("failed to parse createdAt: %w", err)
	}

	link := Link{
		ID:        id,
		TargetURL: fields[fieldTargetUrl],
		CreatedAt: createdAt,
	}

	return link, nil
}

func (r *repository) GetTarget(ctx context.Context, id string) (string, error) {
	return r.redis.Get(ctx, fmt.Sprintf(keyTemplateTarget, id)).Result()
}

func (r *repository) Create(ctx context.Context, target CreateLink) (Link, error) {
	id, err := r.nextID(ctx)
	if err != nil {
		return Link{}, fmt.Errorf("failed to generate id: %w", err)
	}

	link := Link{
		ID:        id,
		TargetURL: target.TargetURL,
		CreatedAt: time.Now(),
	}

	keyTarget := fmt.Sprintf(keyTemplateTarget, link.ID)

	done, err := r.redis.SetNX(ctx, keyTarget, link.TargetURL, r.ttl).Result()
	if !done || err != nil {
		return link, fmt.Errorf("failed to set link: %w", err)
	}

	keyMeta := fmt.Sprintf(keyTemplateMeta, link.ID)

	pipe := r.redis.Pipeline()
	pipe.HSet(ctx, keyMeta, map[string]string{
		fieldTargetUrl: link.TargetURL,
		fieldCreatedAt: link.CreatedAt.Format(time.RFC3339),
	})
	pipe.Expire(ctx, keyMeta, r.ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return link, fmt.Errorf("failed to set link: %w", err)
	}

	return link, nil
}

func (r *repository) nextID(_ context.Context) (string, error) {
	var randomValue uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &randomValue)
	if err != nil {
		return "", fmt.Errorf("failed to read random value: %w", err)
	}

	id := base58.FlickrEncoding.EncodeUint64(uint64(randomValue))

	return string(id), nil
}

func newRepository(redis *redis.Client, config Config) *repository {
	return &repository{
		redis: redis,

		ttl: config.TTL,
	}
}
