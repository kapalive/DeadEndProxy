// Package assets provides embedded static resources.
// © 2023 Devinsidercode CORP. Licensed under the MIT License.
package assets

import (
	"embed"
	"io/fs"
)

//go:embed /assets/static/*
var staticFiles embed.FS

// Static returns the embedded static file system.
func Static() fs.FS {
	f, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return staticFiles
	}
	return f
}
