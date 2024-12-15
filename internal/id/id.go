package id

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"github.com/itchyny/base58-go"
)

type Generator struct {
	encoder *base58.Encoding
}

func (g *Generator) New() (string, error) {
	var randomValue uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &randomValue)
	if err != nil {
		return "", fmt.Errorf("failed to read random value: %w", err)
	}

	id := g.encoder.EncodeUint64(uint64(randomValue))

	return string(id), nil
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

func NewGenerator() *Generator {
	return &Generator{
		encoder: base58.FlickrEncoding,
	}
}
