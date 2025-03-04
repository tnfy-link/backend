package id

import (
	"context"
	"fmt"

	"github.com/itchyny/base58-go"
)

type idProvider interface {
	New(context.Context) (uint32, error)
}

type Generator struct {
	encoder    *base58.Encoding
	idProvider idProvider
}

func (g *Generator) New(ctx context.Context) (string, error) {
	randomValue, err := g.idProvider.New(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate new id: %w", err)
	}

	return string(g.encoder.EncodeUint64(uint64(randomValue))), nil
}

func (g *Generator) Validate(id string) error {
	if id == "" {
		return ErrInvalidID
	}
	if _, err := g.encoder.DecodeUint64([]byte(id)); err != nil {
		return ErrInvalidID
	}
	return nil
}

func NewGenerator(source idProvider) *Generator {
	if source == nil {
		panic("source is required")
	}

	return &Generator{
		encoder:    base58.FlickrEncoding,
		idProvider: source,
	}
}
