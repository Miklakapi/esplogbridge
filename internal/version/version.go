package version

import (
	"runtime/debug"
)

func VersionString() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	if info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "devel"
	}

	return info.Main.Version
}
