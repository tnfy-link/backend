package provider

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const keyLastId = "links:lastId"

type CombinedProvider struct {
	storage *redis.Client
}

func (g *CombinedProvider) New(ctx context.Context) (uint32, error) {
	// Generate 16-bit random part for the lower bits of the ID
	var randomPart uint16
	err := binary.Read(rand.Reader, binary.BigEndian, &randomPart)
	if err != nil {
		return 0, fmt.Errorf("failed to read random value: %w", err)
	}

	// Get sequential part for the upper bits of the ID
	sequentialPart, err := g.storage.Incr(ctx, keyLastId).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment last id: %w", err)
	}

	// Combine sequential part (upper 16 bits) and random part (lower 16 bits)
	return uint32(sequentialPart)<<16 | uint32(randomPart), nil
}

func NewCombinedGenerator(storage *redis.Client) *CombinedProvider {
	if storage == nil {
		panic("storage cannot be nil")
	}

	return &CombinedProvider{
		storage: storage,
	}
}
