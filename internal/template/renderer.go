package template

// Render parses the template and renders it using the provided YAML/JSON data, based on the mode.
func Render(tmplStr string, dataStr string, mode string, prometheusURL string) (string, error) {
	switch mode {
	case "prometheus":
		return RenderPrometheus(tmplStr, dataStr, prometheusURL)
	case "alertmanager":
		return RenderAlertmanager(tmplStr, dataStr)
	default:
		return RenderAlertmanager(tmplStr, dataStr)
	}
}
