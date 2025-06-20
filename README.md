# DeadEndProxy

DeadEndProxy is a lightweight reverse proxy written in Go. The proxy reads a YAML
configuration file and supports automatic reloads when the file changes. Dynamic
routing is resolved through Redis caching, DNS TXT records and a fallback core
API.

---

## ✨ Features

- HTTP & HTTPS proxying with optional redirect from HTTP to HTTPS
- Dynamic domain resolution via:
   - Redis cache
   - DNS TXT records
   - Core API fallback
- YAML-based configuration with hot reload
- Ready to run as a `systemd` service on Linux
- No Nginx or Apache required

---

## Building

```bash
# build the binary
GOOS=linux GOARCH=amd64 go build -o deadendproxy-bin ./cmd
```

## Configuration

routes:
- domain: picture-proof.com
  target: http://127.0.0.1:5000

- domain: manage.eyesync.app
  target: http://127.0.0.1:7000



The default configuration file is `config.yaml`. Copy it to `/etc/deadendproxy/config.yaml`
and edit it to define your domains and routes.

## Running as a service

An example `systemd` service file is provided in `systemd/deadendproxy.service`.
Install it to `/etc/systemd/system/deadendproxy.service` and reload systemd:

```bash
sudo mkdir /etc/deadendproxy/
sudo setcap 'cap_net_bind_service=+ep' ./deadendproxy-bin

./deadendproxy-bin \
  -config /etc/deadendproxy/config.yaml \
  -port-http 80 \
  -port-proxy 443
  

```
/etc/systemd/system/deadendproxy.service

```vim
[Unit]
Description=DeadEnd Reverse Proxy
After=network.target

[Service]
User=deadendproxy
Group=deadendproxy
WorkingDirectory=/etc/deadendproxy
ExecStart=/usr/local/bin/deadendproxy-bin -config /etc/deadendproxy/config.yaml -port-http 80 -port-proxy 443
Restart=always
RestartSec=2

[Install]
WantedBy=multi-user.target
```

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin deadendproxy
sudo mkdir -p /etc/deadendproxy
sudo cp config.yaml /etc/deadendproxy/
sudo chown -R deadendproxy:deadendproxy /etc/deadendproxy
sudo chmod 640 /etc/deadendproxy/config.yaml

sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/deadendproxy-bin
```

After installation you can manage the proxy like any other service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now deadendproxy

sudo systemctl stop deadendproxy
```

## Editing the config quickly

A helper script `scripts/deadendproxy` opens the configuration file when invoked
with the `config` argument. Copy the script and binary to `/usr/local/bin`:

```bash
sudo cp deadendproxy-bin /usr/local/bin/
sudo cp scripts/deadendproxy /usr/local/bin/deadendproxy
```

Now running `sudo deadendproxy config` opens `/etc/deadendproxy/config.yaml` in
`$EDITOR` (defaults to `nano`). Other arguments are forwarded to the binary.

## Dynamic routing via DNS TXT records

When a request arrives, the resolver performs the following steps in order:

1. Check Redis for a cached entry `routing:<domain>`.
2. If not found, lookup TXT records for the domain. A record in the format
   `username_website_<USER>` indicates that requests should be proxied to the
   service for `<USER>` (by default `http://127.0.0.1:9999`). The result is cached
   in Redis for one hour.
3. If no TXT record is present, the resolver queries the core API at
   `http://127.0.0.1:8080/core/domains/resolve?domain=<domain>` to determine the
   target, then caches the result.

This mechanism allows customers to point their own domain at the proxy and
control routing using a simple DNS TXT entry.
