package assets

import (
	"embed"
)

// UIAssets is the embedded UI assets.
//
//go:embed all:ui/dist
var UIAssets embed.FS
