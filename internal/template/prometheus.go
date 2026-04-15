package template

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		"date":               date,
		"tz":                 tz,
		"since":              time.Since,
		"list":               list,
		"append":             appendFunc,
		"dict":               dict,
		"urlUnescape":        url.QueryUnescape,
		"humanize":           humanize,
		"humanize1024":       humanize1024,
		"humanizeDuration":   templates.HumanizeDuration,
		"humanizeTimestamp":  humanizeTimestamp,
		"humanizePercentage": humanizePercentage,
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
		"round":      round,
		"toJS":       toJson,
		"toJson":     toJson,
		"toTime":     toTime,
		"toDuration": toDuration,
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
		ResultType string          `json:"resultType"`
		Result     json.RawMessage `json:"result"`
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

	var samples []QueryResultSample
	switch apiResp.Data.ResultType {
	case "vector":
		var vector []struct {
			Metric map[string]string `json:"metric"`
			Value  []any             `json:"value"`
		}
		if err := json.Unmarshal(apiResp.Data.Result, &vector); err != nil {
			return nil, err
		}
		samples = make([]QueryResultSample, len(vector))
		for i, v := range vector {
			samples[i].Labels = v.Metric
			if len(v.Value) >= 2 {
				valStr, ok := v.Value[1].(string)
				if ok {
					val, err := strconv.ParseFloat(valStr, 64)
					if err == nil {
						samples[i].Value = val
					}
				}
			}
		}
	case "scalar":
		var scalar []any
		if err := json.Unmarshal(apiResp.Data.Result, &scalar); err != nil {
			return nil, err
		}
		if len(scalar) >= 2 {
			valStr, ok := scalar[1].(string)
			if ok {
				val, err := strconv.ParseFloat(valStr, 64)
				if err == nil {
					samples = []QueryResultSample{{Value: val}}
				}
			}
		}
	case "matrix":
		return nil, fmt.Errorf("result type matrix is not supported in the 'query' function yet")
	case "string":
		var strResult []any
		if err := json.Unmarshal(apiResp.Data.Result, &strResult); err != nil {
			return nil, err
		}
		if len(strResult) >= 2 {
			// Prometheus string results aren't directly usable as samples with values,
			// but we can return them as a sample with Value 0 and labels if needed,
			// or just error out as it's not common in templates.
			return nil, fmt.Errorf("result type string is not supported")
		}
	default:
		return nil, fmt.Errorf("unknown result type: %s", apiResp.Data.ResultType)
	}

	return samples, nil
}
