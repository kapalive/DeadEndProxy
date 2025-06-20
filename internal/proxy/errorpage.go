// © 2023 Devinsidercode CORP. Licensed under the MIT License.
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
      background: linear-gradient(135deg, #0c1445 0%%, #1a237e 50%%, #000051 100%%);
      background-attachment: fixed;
      color: #ffffff;
      text-align: center;
      padding-top: 80px;
      min-height: 100vh;
      margin: 0;
    }
    img {
      width: 740px;
      height: 76px;
      margin-bottom: 40px;
    }
    h1 {
      font-size: 38px;
      margin-bottom: 12px;
      color: #ffffff;
    }
    p {
      color: #e0e0e0;
    }
    .footer {
      position: fixed;
      bottom: 20px;
      width: 100%%;
      font-size: 14px;
      color: #b0b0b0;
    }
    .footer a {
      color: #64b5f6;
      text-decoration: none;
    }
    .footer a:hover {
      color: #90caf9;
      text-decoration: underline;
    }
  </style>
</head>
<body>
  <img src="/static/logo-full.png" alt="Devinsider Proxy">
  <h1>%d %s</h1>
  <p>Something went wrong while proxying your request.</p>
  <div class="footer">
<a href="https://devinsidercode.com" target="_blank" rel="noopener noreferrer">Devinsidercode CORP</a> &copy; %d. All Rights Reserved.
  </div>
</body>
</html>`, status, text, status, text, time.Now().Year())

	_, _ = w.Write([]byte(body))
}
