package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JetManiack/alertmanager-template-preview/internal/template"
)

func TestRequestBodySizeLimit(t *testing.T) {
	router := SetupRouter("")

	// Create a large body (> 1MB)
	largeData := strings.Repeat("a", 1024*1024+100)
	body := RenderRequest{
		Template: "test",
		Data:     largeData,
	}

	importBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/render", bytes.NewBuffer(importBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// http.MaxBytesReader causes 413 or an error during reading
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusRequestEntityTooLarge {
		// Gin handles MaxBytesReader error as 500 by default unless we handle it
		// But let's check it failed
		if w.Code == http.StatusOK {
			t.Errorf("Expected failure for large request body, got %v", w.Code)
		}
	}
}

func TestRenderTimeout(t *testing.T) {
	// This is harder to test via API without a slow template
	// Let's test the template.Render directly with a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := template.Render(ctx, "test", "{}", "alertmanager", "")
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}
