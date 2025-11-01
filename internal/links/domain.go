package links

import (
	"fmt"
	"net/url"
)

type NewLink struct {
	targetURL string
}

func NewNewLink(targetURL string) NewLink {
	return NewLink{
		targetURL: targetURL,
	}
}

func (n NewLink) TargetURL() string {
	return n.targetURL
}

const (
	maxTargetURLLength = 2048
)

func (n NewLink) Validate() error {
	if n.targetURL == "" {
		return fmt.Errorf("%w: targetURL is empty", ErrValidationFailed)
	}
	if len(n.targetURL) > maxTargetURLLength {
		return fmt.Errorf("%w: targetURL too long", ErrValidationFailed)
	}

	parsedURL, err := url.Parse(n.targetURL)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	if parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: scheme must be https", ErrValidationFailed)
	}

	return nil
}
