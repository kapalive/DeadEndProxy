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
  - domain: your-domain.com
    ssl:
      cert_file: /etc/letsencrypt/live/your-domain.com/fullchain.pem
      key_file: /etc/letsencrypt/live/your-domain.com/privkey.pem
    redirect_to_https: true
    routes:
      - path: "/api/"
        proxy_pass: "http://127.0.0.1:8080"
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


🐡 BSD Manual Install (FreeBSD / TrueNAS / HardenedBSD)
```
pkg install -y go git
cd /usr/local/src
git clone https://github.com/devinsidercode/DeadEndProxy.git
cd DeadEndProxy
go build -o deadendproxy-bin ./cmd
mkdir -p /usr/local/etc/deadendproxy/webroot/static
cp config.yaml /usr/local/etc/deadendproxy/
cp webroot/static/logo-full.png /usr/local/etc/deadendproxy/webroot/static/
cp deadendproxy-bin /usr/local/sbin/
```

Create rc.d script

Save to /usr/local/etc/rc.d/deadendproxy:
```
#!/bin/sh
# PROVIDE: deadendproxy
# REQUIRE: DAEMON
# KEYWORD: shutdown
. /etc/rc.subr

name="deadendproxy"
rcvar=deadendproxy_enable

load_rc_config $name
: ${deadendproxy_enable:="NO"}

pidfile="/var/run/${name}.pid"
deadendproxy_command="/usr/local/sbin/deadendproxy-bin"
deadendproxy_flags="-config /usr/local/etc/deadendproxy/config.yaml -port-http 80 -port-proxy 443"

command="/usr/sbin/daemon"
command_args="-f -p ${pidfile} ${deadendproxy_command} ${deadendproxy_flags}"

run_rc_command "$1"
```

```
chmod +x /usr/local/etc/rc.d/deadendproxy
sysrc deadendproxy_enable=YES
service deadendproxy start
```


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

Or on BSD:
```
service deadendproxy stop
rm /usr/local/etc/rc.d/deadendproxy
rm -rf /usr/local/etc/deadendproxy
rm /usr/local/sbin/deadendproxy-bin
```

© 2023 Devinsidercode CORP. Licensed under the MIT License.
