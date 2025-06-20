# DeadEndProxy

DeadEndProxy is a lightweight reverse proxy written in Go.  
It reads a YAML configuration file and supports hot reload, TLS with SNI, dynamic domain routing via Redis, DNS TXT, and fallback API.  
Perfect for SaaS platforms where each customer brings their own domain.

---

## ✨ Features

- 🧩 Dynamic domain resolution via:
    - Redis cache (`routing:<domain>`)
    - DNS TXT records (`username_website_<user>`)
    - Core API fallback
- ⚡ Fast Go-based reverse proxy with HTTP & HTTPS
- 🔁 Hot-reloadable YAML config (`config.yaml`)
- 🔐 TLS with SNI support
- 🧾 `systemd` integration with ready-to-use service file
- 🖼 Static assets support (`/static/`)
- 🛠 No need for Nginx or Apache

---

## 🛠 Project Structure
```
DeadEndProxy/
├── LICENSE
├── README.md
├── assets/
│ └── embed.go
├── cmd/
│ └── main.go
├── config/
│ └── config.go
├── config.yaml # Example config
├── deadendproxy-bin # Compiled binary
├── devinsider-proxy-v1.0 # Optional release artifact
├── go.mod
├── go.sum
├── install.sh # Full auto-install script
├── internal/
│ ├── proxy/
│ │ ├── cors.go
│ │ ├── dynamic.go
│ │ ├── errorpage.go
│ │ ├── override.go
│ │ ├── proxy.go
│ │ ├── redirect.go
│ │ ├── reverseproxy.go
│ │ └── router.go
│ ├── router/
│ │ └── resolver.go
│ └── security/
│ └── security.go
├── scripts/
│ └── deadendproxy # CLI wrapper script
├── systemd/
│ └── deadendproxy.service # systemd unit file
├── test.png # Debug downloaded image
└── webroot/
└── static/
└── logo-full.png # Static asset
```

---

## ⚙️ Building

```bash
go build -o deadendproxy-bin ./cmd
````
🧾 Configuration

Put your config.yaml in /etc/deadendproxy/config.yaml.
Minimal example:
```yaml
listen:
  http: ":80"
  https: ":443"

domains:
  - domain: picture-proof.com
    ssl:
      cert_file: /etc/letsencrypt/live/picture-proof.com/fullchain.pem
      key_file: /etc/letsencrypt/live/picture-proof.com/privkey.pem
    redirect_to_https: true
    routes:
      - path: "/core/"
        proxy_pass: "http://127.0.0.1:8085"
      - path: "/storage/"
        proxy_pass: "http://127.0.0.1:9090"
      - path: "/"
        proxy_pass: "http://127.0.0.1:3000"
```
Don't forget to place your static files like logo-full.png into:

/etc/deadendproxy/webroot/static/

🚀 Quick Install

chmod +x install.sh
./install.sh

The script will:

    Create /etc/deadendproxy/ and copy config.yaml and static files

    Compile the binary and copy it to /usr/local/bin/

    Install the systemd unit

    Reload and restart the service

🧰 Manual systemd Setup

# Compile and set permissions
```bash
go build -o deadendproxy-bin ./cmd
sudo cp deadendproxy-bin /usr/local/bin/
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/deadendproxy-bin
```

# Copy config and static files
```bash
sudo mkdir -p /etc/deadendproxy/webroot/static/
sudo cp config.yaml /etc/deadendproxy/
sudo cp webroot/static/logo-full.png /etc/deadendproxy/webroot/static/
```

# Add user
```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin deadendproxy
sudo chown -R deadendproxy:deadendproxy /etc/deadendproxy
```

# Copy and enable systemd service
```bash
sudo cp systemd/deadendproxy.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now deadendproxy
```

🔁 Editing config live

```bash
sudo deadendproxy config
```

This opens /etc/deadendproxy/config.yaml in $EDITOR (default: nano).
You can symlink a helper script for convenience:

sudo cp scripts/deadendproxy /usr/local/bin/deadendproxy

🌐 Dynamic Routing via DNS TXT

On incoming request:

    Check Redis: routing:<domain>

    If not found → TXT record like username_website_john

    If not found → API call to
    http://127.0.0.1:8080/core/domains/resolve?domain=<domain>

This allows customer domains to be dynamically routed based on their DNS or backend configuration.
🧼 Uninstall

```bash
sudo systemctl stop deadendproxy
sudo systemctl disable deadendproxy
sudo rm /usr/local/bin/deadendproxy-bin
sudo rm /etc/systemd/system/deadendproxy.service
sudo rm -rf /etc/deadendproxy
sudo userdel deadendproxy
```

© 2023 Devinsidercode CORP. Licensed under the MIT License.
