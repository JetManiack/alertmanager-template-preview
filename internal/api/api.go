package api

import (
	"net/http"

	"github.com/JetManiack/alertmanager-template-preview/internal/template"
	"github.com/gin-gonic/gin"
)

type RenderRequest struct {
	Template string `json:"template" binding:"required"`
	Data     string `json:"data" binding:"required"`
	Mode     string `json:"mode"` // "alertmanager" (default) or "prometheus"
}

// RenderHandler handles the template rendering request.
func RenderHandler(c *gin.Context) {
	var req RenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mode := req.Mode
	if mode == "" {
		mode = "alertmanager"
	}

	result, err := template.Render(req.Template, req.Data, mode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// SetupRouter initializes the Gin engine with all routes.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/api/render", RenderHandler)

	SetupUI(r)

	return r
}
