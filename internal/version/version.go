// Package version provides build-time version information for openclaw-top.
package version

// These variables are set via ldflags during build.
// For go install, they will have default values.
var (
	// Version is the semantic version of the binary.
	Version = "dev"
	// Commit is the git commit hash.
	Commit = "unknown"
	// Date is the build date.
	Date = "unknown"
)

// Info returns formatted version information.
func Info() string {
	return Version
}

// FullInfo returns full version information including commit and date.
func FullInfo() string {
	return Version + " (commit: " + Commit + ", built: " + Date + ")"
}