package template

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		tmpl     string
		data     string
		mode     string
		expected string
		wantErr  bool
	}{
		{
			name:     "Alertmanager: Simple label render",
			tmpl:     "{{ .CommonLabels.alertname }}",
			data:     `{"commonLabels": {"alertname": "TestAlert"}}`,
			mode:     "alertmanager",
			expected: "TestAlert",
			wantErr:  false,
		},
		{
			name:     "Alertmanager: Template with function",
			tmpl:     "{{ .CommonLabels.alertname | toUpper }}",
			data:     `{"commonLabels": {"alertname": "testalert"}}`,
			mode:     "alertmanager",
			expected: "TESTALERT",
			wantErr:  false,
		},
		{
			name:     "Prometheus: Simple value render",
			tmpl:     "Value is {{ .Value | humanize }}",
			data:     `{"value": 1234.56}`,
			mode:     "prometheus",
			expected: "Value is 1.235k",
			wantErr:  false,
		},
		{
			name:     "Prometheus: Labels render",
			tmpl:     "Alert {{ .Labels.alertname }} is {{ .Value }}",
			data:     "labels:\n  alertname: HighLoad\nvalue: 99",
			mode:     "prometheus",
			expected: "Alert HighLoad is 99",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.tmpl, tt.data, tt.mode, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("Render() got = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRenderWithQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"instance":"localhost:9090","job":"prometheus"},"value":[1617181717.333,"1"]}]}}`)
	}))
	defer ts.Close()

	tmpl := `{{ with query "up" | first }}{{ .Labels.job }}: {{ .Value }}{{ end }}`
	got, err := Render(tmpl, `{}`, "prometheus", ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "prometheus: 1"
	if got != expected {
		t.Errorf("Render() got = %q, want %q", got, expected)
	}
}

func TestRenderWithQueryScalar(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status":"success","data":{"resultType":"scalar","result":[1617181717.333,"123.45"]}}`)
	}))
	defer ts.Close()

	tmpl := `{{ query "time()" | first | value }}`
	got, err := Render(tmpl, `{}`, "prometheus", ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "123.45"
	if got != expected {
		t.Errorf("Render() got = %q, want %q", got, expected)
	}
}
