package template

import (
	"encoding/json"
	"fmt"

	"github.com/prometheus/alertmanager/template"
)

// Render parses the Alertmanager template and renders it using the provided JSON alert data.
func Render(tmplStr string, jsonData string) (string, error) {
	var data template.Data
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal alert data: %w", err)
	}

	// FromGlobs(nil) initializes a template with default functions and built-in templates.
	tmpl, err := template.FromGlobs(nil)
	if err != nil {
		return "", fmt.Errorf("failed to initialize templates: %w", err)
	}

	return tmpl.ExecuteTextString(tmplStr, &data)
}
