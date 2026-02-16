package generator

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/amine7536/preview-generator/config"
	"gopkg.in/yaml.v3"
)

const (
	SyncWaveProject  = "-3"
	SyncWaveNamespace = "-2"
	SyncWaveInfra    = "0"
	SyncWaveService  = "1"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"quote": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
	}
	templates = template.Must(
		template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl"),
	)
}

type appProjectData struct {
	Slug      string
	Namespace string
	SyncWave  string
}

type namespaceData struct {
	Namespace string
	SyncWave  string
}

type infraAppData struct {
	Name           string
	Slug           string
	RepoURL        string
	Chart          string
	TargetRevision string
	ValuesYAML     string
	Namespace      string
	SyncWave       string
}

type serviceAppData struct {
	Name        string
	Slug        string
	ServiceName string
	ImageTag    string
	Namespace   string
	SyncWave    string
}

func Generate(cfg *config.AppsConfig, slug string) (string, error) {
	var out strings.Builder

	if err := templates.ExecuteTemplate(&out, "appproject.yaml.tmpl", appProjectData{
		Slug:      slug,
		Namespace: cfg.Namespace,
		SyncWave:  SyncWaveProject,
	}); err != nil {
		return "", fmt.Errorf("appproject template: %w", err)
	}

	if err := templates.ExecuteTemplate(&out, "namespace.yaml.tmpl", namespaceData{
		Namespace: cfg.Namespace,
		SyncWave:  SyncWaveNamespace,
	}); err != nil {
		return "", fmt.Errorf("namespace template: %w", err)
	}

	for _, infra := range cfg.Infra {
		valuesYAML, err := marshalValuesYAML(infra.Values)
		if err != nil {
			return "", fmt.Errorf("infra %s values: %w", infra.Name, err)
		}

		if err := templates.ExecuteTemplate(&out, "infra-app.yaml.tmpl", infraAppData{
			Name:           cfg.Namespace + "-" + infra.Name,
			Slug:           slug,
			RepoURL:        infra.RepoURL,
			Chart:          infra.Chart,
			TargetRevision: infra.TargetRevision,
			ValuesYAML:     valuesYAML,
			Namespace:      cfg.Namespace,
			SyncWave:       SyncWaveInfra,
		}); err != nil {
			return "", fmt.Errorf("infra-app template: %w", err)
		}
	}

	for _, svc := range cfg.Services {
		if err := templates.ExecuteTemplate(&out, "service-app.yaml.tmpl", serviceAppData{
			Name:        cfg.Namespace + "-" + svc.Name,
			Slug:        slug,
			ServiceName: svc.Name,
			ImageTag:    svc.ImageTag,
			Namespace:   cfg.Namespace,
			SyncWave:    SyncWaveService,
		}); err != nil {
			return "", fmt.Errorf("service-app template: %w", err)
		}
	}

	return out.String(), nil
}

func marshalValuesYAML(values interface{}) (string, error) {
	valuesJSON, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	var valuesMap interface{}
	json.Unmarshal(valuesJSON, &valuesMap)
	valuesYAML, _ := yaml.Marshal(valuesMap)

	return indentYAML(string(valuesYAML), 8), nil
}

func indentYAML(s string, spaces int) string {
	prefix := strings.Repeat(" ", spaces)
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}
