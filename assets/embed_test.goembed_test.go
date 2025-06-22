// Package assets provides embedded static resources.
// © 2023 Devinsidercode CORP. Licensed under the MIT License.

package assets
// Unit test verifying that the embedded file system contains the
// expected static assets.

import "testing"

func TestStaticFS(t *testing.T) {
	fs := Static()
	f, err := fs.Open("logo-full.png")
	if err != nil {
		t.Fatalf("logo-full.png missing: %v", err)
	}
	defer f.Close()
}
