package proxy

import (
	"fmt"
	"net/http"
)

// writeErrorPage outputs a simple HTML error message similar to nginx style.
func writeErrorPage(w http.ResponseWriter, status int) {
	text := http.StatusText(status)
	if text == "" {
		text = "Error"
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	body := fmt.Sprintf(`<html>
<head><title>%d %s</title></head>
<body>
<center><h1>DevinsiderProxy %d %s</h1></center>
</body>
</html>`, status, text, status, text)
	_, _ = w.Write([]byte(body))
}
