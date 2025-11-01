package provider

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

type RandomProvider struct {
}

func NewRandomGenerator() *RandomProvider {
	return &RandomProvider{}
}

func (g *RandomProvider) New(context.Context) (uint32, error) {
	var val uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &val)
	if err != nil {
		return 0, fmt.Errorf("failed to read random value: %w", err)
	}

	return val, nil
}
