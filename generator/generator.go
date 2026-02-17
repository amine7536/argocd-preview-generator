package generator

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/amine7536/preview-generator/config"
)

const (
	SyncWaveProject   = "-3"
	SyncWaveNamespace = "-2"
	SyncWaveService   = "0"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type appProjectData struct {
	Slug     string
	SyncWave string
}

type namespaceData struct {
	Slug     string
	SyncWave string
}

type serviceAppData struct {
	Name        string
	Slug        string
	ServiceName string
	ImageTag    string
	SyncWave    string
	HelmParams  []helmParam
}

type helmParam struct {
	Name  string
	Value string
}

func Generate(cfg *config.AppsConfig, slug string) (string, error) {
	var out strings.Builder

	templates, err := buildTemplates()
	if err != nil {
		return "", err
	}

	err = templates.ExecuteTemplate(&out, "appproject.yaml.tmpl", appProjectData{
		Slug:     slug,
		SyncWave: SyncWaveProject,
	})
	if err != nil {
		return "", fmt.Errorf("appproject template: %w", err)
	}

	err = templates.ExecuteTemplate(&out, "namespace.yaml.tmpl", namespaceData{
		Slug:     slug,
		SyncWave: SyncWaveNamespace,
	})
	if err != nil {
		return "", fmt.Errorf("namespace template: %w", err)
	}

	for _, svc := range cfg.Services {
		var params []helmParam
		for _, p := range svc.HelmParams {
			params = append(params, helmParam{Name: p.Name, Value: p.Value})
		}

		data := serviceAppData{
			Name:        slug + "-" + svc.Name,
			Slug:        slug,
			ServiceName: svc.Name,
			ImageTag:    svc.ImageTag,
			SyncWave:    SyncWaveService,
			HelmParams:  params,
		}

		err = templates.ExecuteTemplate(&out, "service-app.yaml.tmpl", data)
		if err != nil {
			return "", fmt.Errorf("service-app template: %w", err)
		}
	}

	return out.String(), nil
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
