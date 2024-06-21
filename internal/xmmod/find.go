package xmmod

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

var (
	errCannotFindGoMod = errors.New("cannot find go.mod file")
	errNotDir          = errors.New("not a directory")

	ignoreDirs = []string{".git", ".github", "vendor", ".vscode", ".idea"}
)

func Find(start string) (string, error) {
	if filepath.Base(start) == "go.mod" {
		return start, nil
	}

	stat, err := os.Stat(start)
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		curr, err := filepath.Abs(start)
		if err != nil {
			return "", err
		}
		for {
			p := filepath.Join(curr, "go.mod")
			_, err = os.Stat(p)
			if err == nil {
				return p, nil
			}
			next := filepath.Dir(curr)
			if next == curr {
				break
			}
			curr = next
		}
	}

	return "", errCannotFindGoMod
}

func FindAll(start string) ([]string, error) {
	stat, err := os.Stat(start)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%w: %s", errNotDir, start)
	}
	curr, err := filepath.Abs(start)
	if err != nil {
		return nil, err
	}
	mods := []string{}
	err = filepath.WalkDir(curr, func(path string, d fs.DirEntry, err error) error {
		name := d.Name()
		if slices.Contains(ignoreDirs, name) {
			return filepath.SkipDir
		}
		if name == "go.mod" {
			mods = append(mods, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mods, nil
}
