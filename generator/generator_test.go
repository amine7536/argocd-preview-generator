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

func TestGenerateEmptyServicesAndInfra(t *testing.T) {
	cfg := &config.AppsConfig{
		Namespace: "preview-empty",
		Services:  []config.Service{},
		Infra:     []config.Infra{},
	}

	got, err := generator.Generate(cfg, "empty")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Should contain AppProject and Namespace only
	if !contains(got, "kind: AppProject") {
		t.Error("missing AppProject")
	}
	if !contains(got, "kind: Namespace") {
		t.Error("missing Namespace")
	}
	// Should not contain any Application
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
		Infra: []config.Infra{},
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

func TestGenerateSpecialCharsInValues(t *testing.T) {
	cfg := &config.AppsConfig{
		Namespace: "preview-special",
		Services:  []config.Service{},
		Infra: []config.Infra{
			{
				Name:           "redis",
				Chart:          "redis",
				RepoURL:        "https://charts.bitnami.com/bitnami",
				TargetRevision: "18.x",
				Values: map[string]interface{}{
					"auth": map[string]interface{}{
						"password": "p@ss:w0rd!#$",
					},
				},
			},
		},
	}

	got, err := generator.Generate(cfg, "special")
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	if !contains(got, "p@ss:w0rd!#$") {
		t.Error("special characters not preserved in values")
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
