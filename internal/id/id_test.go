package id_test

import (
	"testing"

	"github.com/tnfy-link/backend/internal/id"
)

func TestGenerator_New(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	g := id.NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.New()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Generator.New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratorValidate(t *testing.T) {
	g := id.NewGenerator()

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid ID",
			id:      func() string { id, _ := g.New(); return id }(),
			wantErr: false,
		},
		{
			name:    "invalid ID",
			id:      " invalid", // an invalid base58 encoded string
			wantErr: true,
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
		},
		{
			name:    "ID with special characters",
			id:      "3ah4!@#",
			wantErr: true,
		},
		{
			name:    "very long ID",
			id:      "3ah4Vb3ah4Vb3ah4Vb3ah4Vb3ah4Vb3ah4Vb",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := g.Validate(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkGeneratorNew(b *testing.B) {
	g := id.NewGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.New()
	}
}

func BenchmarkGeneratorNewParallel(b *testing.B) {
	g := id.NewGenerator()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = g.New()
		}
	})
}

func BenchmarkGeneratorValidate(b *testing.B) {
	g := id.NewGenerator()
	validID, _ := g.New() // Generate a valid ID

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Validate(validID)
	}
}

func BenchmarkGeneratorValidateInvalid(b *testing.B) {
	g := id.NewGenerator()
	invalidID := "invalid-id"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Validate(invalidID)
	}
}

func BenchmarkGeneratorValidateEmpty(b *testing.B) {
	g := id.NewGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Validate("")
	}
}
