package links

import (
	"errors"
	"fmt"
	"net/url"
)

type NewLink struct {
	targetURL string
}

func (n NewLink) TargetURL() string {
	return n.targetURL
}

const (
	maxTargetURLLength = 2048
)

func (n NewLink) Validate() error {
	if n.targetURL == "" {
		return errors.New("targetUrl is empty")
	}
	if len(n.targetURL) > maxTargetURLLength {
		return errors.New("targetUrl too long")
	}

	parsedUrl, err := url.Parse(n.targetURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	if parsedUrl.Scheme != "https" {
		return errors.New("scheme must be https")
	}

	return nil
}

func NewNewLink(targetUrl string) NewLink {
	return NewLink{
		targetURL: targetUrl,
	}
}
