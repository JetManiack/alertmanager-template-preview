package template

import "context"

// Render parses the template and renders it using the provided YAML/JSON data, based on the mode.
func Render(ctx context.Context, tmplStr string, dataStr string, mode string, prometheusURL string) (string, error) {
	switch mode {
	case "prometheus":
		return RenderPrometheus(ctx, tmplStr, dataStr, prometheusURL)
	case "alertmanager":
		return RenderAlertmanager(ctx, tmplStr, dataStr)
	default:
		return RenderAlertmanager(ctx, tmplStr, dataStr)
	}
}
