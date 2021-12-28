package main

import (
	"os"

	"go.reefassistant.com/apex-proxy/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
