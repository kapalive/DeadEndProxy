#!/bin/sh
# // © 2023 Devinsidercode CORP. Licensed under the MIT License.
set -e

echo "[*] Checking the environment..."

if ! command -v go >/dev/null 2>&1; then
    echo "[!] Go is not installed. Install Go and try again."
    exit 1
fi

if ! pidof systemd >/dev/null; then
    echo "[!] Systemd is not running. This script only works on systemd systems."
    exit 1
fi

echo "[*] Creating directories..."
sudo mkdir -p /etc/deadendproxy/webroot/static/

echo "[*] Copying the config..."
sudo cp config.yaml /etc/deadendproxy/

echo "[*] I copy statics..."
sudo cp webroot/static/logo-full.png /etc/deadendproxy/webroot/static/

echo "[*] I'm compiling a binary..."
go build -o deadendproxy-bin ./cmd

echo "[*] I copy the binary to /usr/local/bin/..."
sudo cp deadendproxy-bin /usr/local/bin/

echo "[*] Installing cap_net_bind_service..."
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/deadendproxy-bin

echo "[*] Installing systemd unit..."
sudo mkdir -p /etc/systemd/system/
sudo cp systemd/deadendproxy.service /etc/systemd/system/deadendproxy.service

echo "[*] Restarting systemd..."
sudo systemctl daemon-reload
sudo systemctl enable deadendproxy
sudo systemctl restart deadendproxy

echo "[✔] DeadEndProxy is installed and running as a systemd service."
