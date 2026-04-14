package template

import (
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		tmpl     string
		data     string // Alert data in JSON format
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple label render",
			tmpl:     "{{ .CommonLabels.alertname }}",
			data:     `{"commonLabels": {"alertname": "TestAlert"}}`,
			expected: "TestAlert",
			wantErr:  false,
		},
		{
			name:     "Template with function",
			tmpl:     "{{ .CommonLabels.alertname | toUpper }}",
			data:     `{"commonLabels": {"alertname": "testalert"}}`,
			expected: "TESTALERT",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("Render() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
