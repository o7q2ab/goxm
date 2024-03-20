package xmmod

import (
	"encoding/json"
	"os/exec"
)

type goListOutput struct {
	Versions []string `json:"Versions"`
}

// GetLatest returns the latest available version for the given module.
func GetLatest(modpath string) string {
	cmd := exec.Command(
		// More about `go list`: https://pkg.go.dev/cmd/go#hdr-List_packages_or_modules
		"go", "list", "-m", "-e", "-json=Versions", "-mod=mod", "-versions", modpath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	result := goListOutput{}
	if err = json.Unmarshal(out, &result); err != nil {
		return ""
	}
	if len(result.Versions) == 0 {
		return ""
	}
	return result.Versions[len(result.Versions)-1]
}
