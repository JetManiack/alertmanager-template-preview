package template

import (
	"fmt"
	tmplhtml "html/template"
	tmpltext "text/template"

	"github.com/goccy/go-yaml"
	"github.com/prometheus/alertmanager/template"
)

// Render parses the template and renders it using the provided YAML/JSON data, based on the mode.
func Render(tmplStr string, dataStr string, mode string, prometheusURL string) (string, error) {
	if mode == "prometheus" {
		return RenderPrometheus(tmplStr, dataStr, prometheusURL)
	}

	var data template.Data
	if err := yaml.Unmarshal([]byte(dataStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal alert data: %w", err)
	}

	// FromGlobs(nil) initializes a template with default functions and built-in templates.
	tmpl, err := template.FromGlobs(nil, func(text *tmpltext.Template, html *tmplhtml.Template) {
		text.Funcs(tmpltext.FuncMap{
			"round":              round,
			"toTime":             toTime,
			"toDuration":         toDuration,
			"toJson":             toJson,
			"toJS":               toJson,
			"humanize":           humanize,
			"humanize1024":       humanize1024,
			"humanizeTimestamp":  humanizeTimestamp,
			"humanizePercentage": humanizePercentage,
		})
	})
	if err != nil {
		return "", fmt.Errorf("failed to initialize templates: %w", err)
	}

	return tmpl.ExecuteTextString(tmplStr, &data)
}
