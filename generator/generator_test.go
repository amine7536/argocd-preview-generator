package generator_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/amine7536/preview-generator/config"
	"github.com/amine7536/preview-generator/generator"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

func TestGenerateBasic(t *testing.T) {
	cfg, err := config.Load(filepath.Join(testdataDir(), "basic.yaml"))
	if err != nil {
		t.Fatalf("config.Load() error: %v", err)
	}

	got, err := generator.Generate(cfg, "feature-add-pricing")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	expected, err := os.ReadFile(filepath.Join(testdataDir(), "basic_expected.yaml"))
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}

	if got != string(expected) {
		t.Errorf("Generate() output differs from expected.\nGot:\n%s\nExpected:\n%s", got, string(expected))
	}
}

func TestGenerateEmptyServices(t *testing.T) {
	cfg := &config.AppsConfig{
		Namespace: "preview-empty",
		Services:  []config.Service{},
	}

	got, err := generator.Generate(cfg, "empty")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if !contains(got, "kind: AppProject") {
		t.Error("missing AppProject")
	}
	if !contains(got, "kind: Namespace") {
		t.Error("missing Namespace")
	}
	if contains(got, "kind: Application") {
		t.Error("unexpected Application in empty config")
	}
}

func TestGenerateServicesOnly(t *testing.T) {
	cfg := &config.AppsConfig{
		Namespace: "preview-svc-only",
		Services: []config.Service{
			{Name: "api", ImageTag: "sha123"},
		},
	}

	got, err := generator.Generate(cfg, "svc-only")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if !contains(got, `"preview-svc-only-api"`) {
		t.Error("missing service application name")
	}
	if !contains(got, `value: "sha123"`) {
		t.Error("missing image tag")
	}
}

func TestGenerateBackend1DatabaseName(t *testing.T) {
	cfg := &config.AppsConfig{
		Namespace: "preview-test",
		Services: []config.Service{
			{Name: "backend-1", ImageTag: "abc123"},
		},
	}

	got, err := generator.Generate(cfg, "my-feature")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if !contains(got, "database.name") {
		t.Error("missing database.name helm parameter for backend-1")
	}
	if !contains(got, `"backend-1-my-feature"`) {
		t.Error("database.name should be backend-1-<slug>")
	}
}

func TestGenerateNonBackend1NoDatabaseParam(t *testing.T) {
	cfg := &config.AppsConfig{
		Namespace: "preview-test",
		Services: []config.Service{
			{Name: "front", ImageTag: "abc123"},
		},
	}

	got, err := generator.Generate(cfg, "my-feature")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if contains(got, "database.name") {
		t.Error("front should not have database.name helm parameter")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
