package xmmod

import (
	"os"

	"golang.org/x/mod/modfile"
)

func Read(path string) (*modfile.File, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return modfile.Parse(path, f, nil)
}
