package proxy

import (
	"fmt"
	"net/http"
	"time"
)

// writeErrorPage outputs a simple HTML error message similar to nginx style.
func writeErrorPage(w http.ResponseWriter, status int) {
	text := http.StatusText(status)
	if text == "" {
		text = "Error"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>%d %s</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f4;
      color: #222;
      text-align: center;
      padding-top: 80px;
    }
    img {
      width: 740px;
      height: 76px;
      margin-bottom: 40px;
    }
    h1 {
      font-size: 38px;
      margin-bottom: 12px;
    }
    .footer {
      position: fixed;
      bottom: 20px;
      width: 100%%;
      font-size: 14px;
      color: #999;
    }
  </style>
</head>
<body>
  <img src="/static/logo-full.png" alt="Devinsider Proxy">
  <h1>%d %s</h1>
  <p>Something went wrong while proxying your request.</p>
  <div class="footer">
    Devinsidercode CORP &copy; %d. All Rights Reserved.
  </div>
</body>
</html>`, status, text, status, text, time.Now().Year())

	_, _ = w.Write([]byte(body))
}
