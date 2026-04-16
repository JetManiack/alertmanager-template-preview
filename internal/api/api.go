package api

import (
	"context"
	"net/http"
	"time"

	"github.com/JetManiack/alertmanager-template-preview/internal/template"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RenderRequest struct {
	Template string `json:"template" binding:"required"`
	Data     string `json:"data" binding:"required"`
	Mode     string `json:"mode"` // "alertmanager" (default) or "prometheus"
}

// RenderHandler handles the template rendering request.
func RenderHandler(c *gin.Context, prometheusURL string) {
	var req RenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mode := req.Mode
	if mode == "" {
		mode = "alertmanager"
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	result, err := template.Render(ctx, req.Template, req.Data, mode, prometheusURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// SetupRouter initializes the Gin engine with all routes.
func SetupRouter(prometheusURL string) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/api/render", func(c *gin.Context) {
		// Limit request body size to 1MB to prevent DoS
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1024*1024)
		RenderHandler(c, prometheusURL)
	})

	SetupUI(r)

	return r
}
