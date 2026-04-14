package api

import (
	"io/fs"
	"net/http"

	"github.com/JetManiack/alertmanager-template-preview/assets"
	"github.com/gin-gonic/gin"
)

// SetupUI adds the UI routes to the router.
func SetupUI(r *gin.Engine) {
	// Root of the UI assets
	subFS, err := fs.Sub(assets.UIAssets, "ui/dist")
	if err != nil {
		panic(err)
	}

	// Serve static files at /ui
	// Gin's StaticFS handles serving index.html for directory requests if it exists.
	r.StaticFS("/ui", http.FS(subFS))

	// Redirect root to /ui/
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})
}
