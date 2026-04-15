package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRenderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/render", RenderHandler)

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
		wantBody   string
	}{
		{
			name: "Successful render",
			body: map[string]any{
				"template": "{{ .CommonLabels.alertname }}",
				"data":     `{"commonLabels": {"alertname": "TestAlert"}}`,
				"mode":     "alertmanager",
			},
			wantStatus: http.StatusOK,
			wantBody:   "TestAlert",
		},
		{
			name: "Prometheus render",
			body: map[string]any{
				"template": "{{ .Value | humanize }}",
				"data":     `{"value": 1000}`,
				"mode":     "prometheus",
			},
			wantStatus: http.StatusOK,
			wantBody:   "1k",
		},
		{
			name: "Invalid YAML data",
			body: map[string]any{
				"template": "{{ .CommonLabels.alertname }}",
				"data":     `*invalid_alias`,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/render", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("RenderHandler() status = %v, want %v", w.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				var resp map[string]string
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp["result"] != tt.wantBody {
					t.Errorf("RenderHandler() got body = %v, want %v", resp["result"], tt.wantBody)
				}
			}
		})
	}
}
