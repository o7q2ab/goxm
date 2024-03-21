package build

import (
	"runtime/debug"
)

var version = "unknown"

func Version() string {
	if version != "unknown" {
		return version
	}
	if c := readVersion(); c != "" {
		return c
	}
	return version
}

func readVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" && s.Value != "" {
			return s.Value
		}
	}
	return ""
}
