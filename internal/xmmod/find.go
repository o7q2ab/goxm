package xmmod

import (
	"errors"
	"os"
	"path/filepath"
)

var errCannotFindGoMod = errors.New("cannot find go.mod file")

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
