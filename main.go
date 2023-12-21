package main

import (
	"fmt"
	"os"

	"github.com/o7q2ab/goxm/internal/commands"
)

func main() {
	if err := commands.NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: [%T] %v\n", err, err)
	}
}
