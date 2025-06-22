package assets

import "testing"

func TestStaticFS(t *testing.T) {
    fs := Static()
    f, err := fs.Open("logo-full.png")
    if err != nil {
        t.Fatalf("logo-full.png missing: %v", err)
    }
    defer f.Close()
}
