package build

import (
	"runtime/debug"
)

var version = "unknown"

func Version() string {
	if version != "unknown" {
		return version
	}
	if c := readCommit(); c != "" {
		return c
	}

	return version
}

func readCommit() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, s := range bi.Settings {
		if s.Key == "vsc.revision" && s.Value != "" {
			return s.Value
		}
	}
	return ""
}
