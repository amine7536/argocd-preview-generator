package generator

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/amine7536/preview-generator/config"
)

const (
	SyncWaveProject   = "-3"
	SyncWaveNamespace = "-2"
	SyncWaveInfra     = "0"
	SyncWaveService   = "1"
	valuesIndent      = 8
)

//go:embed templates/*.tmpl
var templateFS embed.FS

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

	templates, err := buildTemplates()
	if err != nil {
		return "", err
	}

	err = templates.ExecuteTemplate(&out, "appproject.yaml.tmpl", appProjectData{
		Slug:      slug,
		Namespace: cfg.Namespace,
		SyncWave:  SyncWaveProject,
	})
	if err != nil {
		return "", fmt.Errorf("appproject template: %w", err)
	}

	err = templates.ExecuteTemplate(&out, "namespace.yaml.tmpl", namespaceData{
		Namespace: cfg.Namespace,
		SyncWave:  SyncWaveNamespace,
	})
	if err != nil {
		return "", fmt.Errorf("namespace template: %w", err)
	}

	for _, infra := range cfg.Infra {
		valuesYAML, valuesErr := marshalValuesYAML(infra.Values)
		if valuesErr != nil {
			return "", fmt.Errorf("infra %s values: %w", infra.Name, valuesErr)
		}

		if execErr := templates.ExecuteTemplate(&out, "infra-app.yaml.tmpl", infraAppData{
			Name:           cfg.Namespace + "-" + infra.Name,
			Slug:           slug,
			RepoURL:        infra.RepoURL,
			Chart:          infra.Chart,
			TargetRevision: infra.TargetRevision,
			ValuesYAML:     valuesYAML,
			Namespace:      cfg.Namespace,
			SyncWave:       SyncWaveInfra,
		}); execErr != nil {
			return "", fmt.Errorf("infra-app template: %w", execErr)
		}
	}

	for _, svc := range cfg.Services {
		err = templates.ExecuteTemplate(&out, "service-app.yaml.tmpl", serviceAppData{
			Name:        cfg.Namespace + "-" + svc.Name,
			Slug:        slug,
			ServiceName: svc.Name,
			ImageTag:    svc.ImageTag,
			Namespace:   cfg.Namespace,
			SyncWave:    SyncWaveService,
		})
		if err != nil {
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
	err = json.Unmarshal(valuesJSON, &valuesMap)
	if err != nil {
		return "", err
	}
	valuesYAML, err := yaml.Marshal(valuesMap)
	if err != nil {
		return "", err
	}

	return indentYAML(string(valuesYAML), valuesIndent), nil
}

func buildTemplates() (*template.Template, error) {
	funcMap := template.FuncMap{
		"quote": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return tmpl, nil
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
