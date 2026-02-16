package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amine7536/preview-generator/config"
)

func TestLoad(t *testing.T) {
	yaml := `namespace: preview-test
services:
  - name: backend-1
    image_tag: "abc123"
  - name: front
    image_tag: "def456"
infra:
  - name: postgres
    chart: postgresql
    repoURL: https://charts.bitnami.com/bitnami
    targetRevision: "*"
    values:
      auth:
        postgresPassword: postgres
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Namespace != "preview-test" {
		t.Errorf("Namespace = %q, want %q", cfg.Namespace, "preview-test")
	}
	if len(cfg.Services) != 2 {
		t.Fatalf("Services count = %d, want 2", len(cfg.Services))
	}
	if cfg.Services[0].Name != "backend-1" || cfg.Services[0].ImageTag != "abc123" {
		t.Errorf("Services[0] = %+v", cfg.Services[0])
	}
	if len(cfg.Infra) != 1 {
		t.Fatalf("Infra count = %d, want 1", len(cfg.Infra))
	}
	if cfg.Infra[0].Chart != "postgresql" {
		t.Errorf("Infra[0].Chart = %q, want %q", cfg.Infra[0].Chart, "postgresql")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("Load() expected error for missing file")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	path := writeTemp(t, "{{invalid yaml")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("Load() expected error for invalid YAML")
	}
}

func TestLoadEmptyServices(t *testing.T) {
	yaml := `namespace: preview-empty
services: []
infra: []
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(cfg.Services) != 0 {
		t.Errorf("Services count = %d, want 0", len(cfg.Services))
	}
	if len(cfg.Infra) != 0 {
		t.Errorf("Infra count = %d, want 0", len(cfg.Infra))
	}
}

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "apps.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
