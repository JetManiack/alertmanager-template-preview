package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	tmpltext "text/template"
	"time"

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

func RenderPrometheus(tmplStr string, dataStr string, prometheusURL string) (string, error) {
	var data PrometheusData
	if err := yaml.Unmarshal([]byte(dataStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal alert data: %w", err)
	}

	funcs := tmpltext.FuncMap{
		"toUpper":   strings.ToUpper,
		"toLower":   strings.ToLower,
		"title":     strings.Title,
		"trimSpace": strings.TrimSpace,
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"match": regexp.MatchString,
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"humanize":          humanize,
		"humanize1024":      humanize1024,
		"humanizeDuration":  templates.HumanizeDuration,
		"humanizeTimestamp": templates.HumanizeTimestamp,
		"humanizePercentage": func(v any) (string, error) {
			f, err := templates.ConvertToFloat(v)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%.4g%%", f*100), nil
		},
		"query": func(q string) (any, error) {
			if prometheusURL == "" {
				return nil, fmt.Errorf("function 'query' requires a live Prometheus server (use --prometheus-url flag)")
			}
			return queryPrometheus(prometheusURL, q)
		},
		"first": func(v any) (any, error) {
			samples, ok := v.([]QueryResultSample)
			if !ok || len(samples) == 0 {
				return nil, nil
			}
			return samples[0], nil
		},
		"last": func(v any) (any, error) {
			samples, ok := v.([]QueryResultSample)
			if !ok || len(samples) == 0 {
				return nil, nil
			}
			return samples[len(samples)-1], nil
		},
		"value": func(v any) (any, error) {
			switch s := v.(type) {
			case QueryResultSample:
				return s.Value, nil
			case *QueryResultSample:
				if s == nil {
					return nil, nil
				}
				return s.Value, nil
			default:
				return nil, fmt.Errorf("value: expected sample, got %T", v)
			}
		},
		"label": func(label string, v any) (any, error) {
			switch s := v.(type) {
			case QueryResultSample:
				return s.Labels[label], nil
			case *QueryResultSample:
				if s == nil {
					return "", nil
				}
				return s.Labels[label], nil
			default:
				return nil, fmt.Errorf("label: expected sample, got %T", v)
			}
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

type QueryResultSample struct {
	Labels map[string]string
	Value  float64
}

type prometheusAPIResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []any             `json:"value"`
		} `json:"result"`
	} `json:"data"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}

func queryPrometheus(baseURL, q string) ([]QueryResultSample, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/api/v1/query"
	u.RawQuery = url.Values{"query": {q}}.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp prometheusAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "success" {
		return nil, fmt.Errorf("prometheus error: %s (%s)", apiResp.Error, apiResp.ErrorType)
	}

	samples := make([]QueryResultSample, len(apiResp.Data.Result))
	for i, r := range apiResp.Data.Result {
		samples[i].Labels = r.Metric
		if len(r.Value) >= 2 {
			valStr, ok := r.Value[1].(string)
			if ok {
				val, err := strconv.ParseFloat(valStr, 64)
				if err == nil {
					samples[i].Value = val
				}
			}
		}
	}

	return samples, nil
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
