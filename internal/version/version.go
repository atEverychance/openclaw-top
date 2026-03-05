package version

import (
	"runtime/debug"
)

// These variables are set via ldflags during build.
var (
	// Version is the semantic version of the binary.
	Version = "dev"
	// Commit is the git commit hash.
	Commit = "unknown"
	// Date is the build date.
	Date = "unknown"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if Version == "dev" {
			Version = info.Main.Version
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if Commit == "unknown" {
					Commit = setting.Value
				}
			case "vcs.time":
				if Date == "unknown" {
					Date = setting.Value
				}
			}
		}
	}
}

// Info returns formatted version information.
func Info() string {
	return Version
}

// FullInfo returns full version information including commit and date.
func FullInfo() string {
	return Version + " (commit: " + Commit + ", built: " + Date + ")"
}