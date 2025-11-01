package version

import "strconv"

const notSet string = "not set"

// These variables are populated at build time using ldflags.
// Example: `go build -ldflags "-X github.com/tnfy-link/backend/internal/version.AppVersion=0.1 -X github.com/tnfy-link/backend/internal/version.AppRelease=123"`.
//
//nolint:gochecknoglobals // build metadata
var (
	AppVersion = notSet
	AppRelease = notSet
)

func AppReleaseID() int {
	id, _ := strconv.Atoi(AppRelease)

	return id
}
