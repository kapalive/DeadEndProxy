// Package assets provides embedded static resources.
// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Static files such as logos are embedded at build time so the
// proxy can serve them without external dependencies.
package assets

import (
	"embed"
	"io/fs"
)
//go:embed static/*
var staticFiles embed.FS

// Static returns the embedded static file system used by the
// HTTP handlers to serve assets.
func Static() fs.FS {
	f, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return staticFiles
	}
	return f
}
