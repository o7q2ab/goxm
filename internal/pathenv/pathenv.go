package pathenv

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func List() []string {
	var raw string
	if runtime.GOOS == "windows" {
		raw = os.Getenv("path")
	} else {
		raw = os.Getenv("PATH")
	}

	names := []string{}
	all := filepath.SplitList(raw)
	for _, one := range all {
		files, err := os.ReadDir(one)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			names = append(names, filepath.Join(one, file.Name()))
		}
	}

	return names
}
