// © 2023 Devinsidercode CORP. Licensed under the MIT License.
//
// Package proxy contains helpers for rendering a simple HTML
// error page when upstream services are unavailable.
package proxy

import (
	"fmt"
	"net/http"
	"time"
)

// writeErrorPage renders a minimalistic HTML error page similar
// to the style used by nginx.
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
  <svg
   version="1.1"
   id="svg1"
   width="740"
   height="76.000008"
   viewBox="0 0 740 76.000008"
   sodipodi:docname="logo-full.png"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
   xmlns:xlink="http://www.w3.org/1999/xlink"
   xmlns="http://www.w3.org/2000/svg"
   xmlns:svg="http://www.w3.org/2000/svg">
  <defs
     id="defs1" />
  <sodipodi:namedview
     id="namedview1"
     pagecolor="#ffffff"
     bordercolor="#000000"
     borderopacity="0.25"
     inkscape:showpageshadow="2"
     inkscape:pageopacity="0.0"
     inkscape:pagecheckerboard="0"
     inkscape:deskcolor="#d1d1d1">
    <inkscape:page
       x="0"
       y="0"
       width="740"
       height="76.000008"
       id="page2"
       margin="0"
       bleed="0" />
  </sodipodi:namedview>
  <g
     inkscape:groupmode="layer"
     inkscape:label="Image"
     id="g1"
     transform="translate(8.2698679,-112.71523)">
    <image
       width="740"
       height="76"
       preserveAspectRatio="none"
       xlink:href="data:image/png;base64,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&#10;"
       id="image1"
       x="-8.2698679"
       y="112.71523" />
  </g>
</svg>
  <h1>%d %s</h1>
  <p>Something went wrong while proxying your request.</p>
  <div class="footer">
<a href="https://devinsidercode.com" target="_blank" rel="noopener noreferrer">Devinsidercode CORP</a> &copy; %d. All Rights Reserved.
  </div>
</body>
</html>`, status, text, status, text, time.Now().Year())

	_, _ = w.Write([]byte(body))
}
