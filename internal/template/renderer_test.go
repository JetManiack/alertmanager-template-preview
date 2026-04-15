package template

import (
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
			got, err := Render(tt.tmpl, tt.data, tt.mode)
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
