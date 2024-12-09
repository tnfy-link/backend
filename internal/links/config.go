package links

import "time"

type Config struct {
	Hostname string
	TTL      time.Duration
}
