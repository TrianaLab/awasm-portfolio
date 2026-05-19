//go:build !wasm

package main

import (
	"github.com/TrianaLab/awasm-portfolio/cmd"
	"github.com/TrianaLab/awasm-portfolio/internal/preload"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"os"
)

func main() {
	repo := repository.NewInMemoryRepository()
	preload.PreloadData(repo)
	rootCmd := cmd.NewRootCommand(repo)

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
