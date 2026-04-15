package template

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	tmpltext "text/template"

	"github.com/goccy/go-yaml"
	"github.com/prometheus/common/helpers/templates"
)

// PrometheusData is the data structure passed to Prometheus templates (e.g. for alerting rules).
type PrometheusData struct {
	Labels         map[string]string `json:"labels"`
	ExternalLabels map[string]string `json:"externalLabels"`
	ExternalURL    string            `json:"externalURL"`
	Value          float64           `json:"value"`
}

func RenderPrometheus(tmplStr string, dataStr string) (string, error) {
	var data PrometheusData
	if err := yaml.Unmarshal([]byte(dataStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal alert data: %w", err)
	}

	funcs := tmpltext.FuncMap{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"title":   strings.Title,
		"trimSpace": strings.TrimSpace,
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"match": regexp.MatchString,
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"humanize":         humanize,
		"humanize1024":     humanize1024,
		"humanizeDuration": templates.HumanizeDuration,
		"humanizeTimestamp": templates.HumanizeTimestamp,
		"humanizePercentage": func(v any) (string, error) {
			f, err := templates.ConvertToFloat(v)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%.4g%%", f*100), nil
		},
		"query": func(q string) (any, error) {
			return nil, fmt.Errorf("function 'query' is not supported in the previewer yet (requires a live Prometheus server)")
		},
		"first": func(v any) (any, error) {
			return nil, fmt.Errorf("function 'first' is not supported in the previewer yet")
		},
		"last": func(v any) (any, error) {
			return nil, fmt.Errorf("function 'last' is not supported in the previewer yet")
		},
		"value": func(v any) (any, error) {
			return nil, fmt.Errorf("function 'value' is not supported in the previewer yet")
		},
	}

	tmpl, err := tmpltext.New("prometheus").Funcs(funcs).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func humanize(v any) (string, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return "", err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Sprintf("%.4g", f), nil
	}
	if math.Abs(f) < 1 {
		return fmt.Sprintf("%.4g", f), nil
	}

	units := []string{"", "k", "M", "G", "T", "P", "E"}
	i := 0
	for math.Abs(f) >= 1000 && i < len(units)-1 {
		f /= 1000
		i++
	}
	return strconv.FormatFloat(f, 'g', 4, 64) + units[i], nil
}

func humanize1024(v any) (string, error) {
	f, err := templates.ConvertToFloat(v)
	if err != nil {
		return "", err
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Sprintf("%.4g", f), nil
	}
	if math.Abs(f) < 1 {
		return fmt.Sprintf("%.4g", f), nil
	}

	units := []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
	i := 0
	for math.Abs(f) >= 1024 && i < len(units)-1 {
		f /= 1024
		i++
	}
	return strconv.FormatFloat(f, 'g', 4, 64) + units[i], nil
}
