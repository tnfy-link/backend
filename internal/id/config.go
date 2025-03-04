package id

type Provider string

const (
	ProviderRandom   = "random"
	ProviderCombined = "combined"
)

type Config struct {
	Provider Provider
}
