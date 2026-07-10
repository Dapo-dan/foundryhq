// Package version holds build metadata for the /version endpoint. Version
// and Commit are set via -ldflags at build time (see apps/api/Dockerfile);
// the defaults below are what `go run`/`go build` without ldflags produce.
package version

var (
	Version = "dev"
	Commit  = "unknown"
)
