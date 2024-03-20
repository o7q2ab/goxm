package xmpath

import (
	"os"
	"path/filepath"
	"runtime"
)

func ListPathEnv() []string {
	var raw string
	if runtime.GOOS == "windows" {
		raw = os.Getenv("path")
	} else {
		raw = os.Getenv("PATH")
	}

	names := []string{}
	all := filepath.SplitList(raw)
	for _, one := range all {
		names = append(names, listdir(one)...)
	}

	return names
}

func List(p string) []string {
	stat, err := os.Stat(p)
	if err != nil {
		return []string{}
	}
	if stat.IsDir() {
		return listdir(p)
	}
	return []string{p}
}

func listdir(p string) []string {
	files, err := os.ReadDir(p)
	if err != nil {
		return nil
	}
	names := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		names = append(names, filepath.Join(p, file.Name()))
	}
	return names
}
