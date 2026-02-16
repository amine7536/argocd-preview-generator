package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amine7536/preview-generator/config"
	"github.com/amine7536/preview-generator/generator"
)

func main() {
	sourcePath := os.Getenv("ARGOCD_APP_SOURCE_PATH")
	if sourcePath == "" {
		fmt.Fprintf(os.Stderr, "ARGOCD_APP_SOURCE_PATH not set\n")
		os.Exit(1)
	}

	slug := filepath.Base(sourcePath)

	cfg, err := config.Load("apps.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	out, err := generator.Generate(cfg, slug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(out)
}
