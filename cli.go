//go:build !wasm

package main

import (
	"awasm-portfolio/cmd"
	"awasm-portfolio/internal/preload"
	"awasm-portfolio/internal/repository"
	"os"
)

func main() {
	repo := repository.NewInMemoryRepository()
	preload.PreloadData(repo)
	rootCmd := cmd.NewRootCommand(repo)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
