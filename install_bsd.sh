#!/bin/sh
# © 2023 Devinsidercode CORP. Licensed under the MIT License.
#
# Automated installation script for DeadEndProxy on FreeBSD.

set -e

echo "[*] Checking the environment..."

if ! command -v go >/dev/null 2>&1; then
    echo "[!] Go is not installed. Please install Go via 'pkg install go'."
    exit 1
fi

echo "[*] Creating directories..."
sudo mkdir -p /usr/local/etc/deadendproxy/webroot/static/

echo "[*] Copying config..."
sudo cp config.yaml /usr/local/etc/deadendproxy/

echo "[*] Copying static files..."
sudo cp webroot/static/logo-full.png /usr/local/etc/deadendproxy/webroot/static/

echo "[*] Building binary..."
go build -buildvcs=false -o deadendproxy-bsd ./cmd

echo "[*] Installing binary to /usr/local/bin/..."
sudo cp deadendproxy-bsd /usr/local/bin/

echo "[✔] DeadEndProxy installed on FreeBSD."

echo
echo "[ℹ] You can now run it manually:"
echo "    $ /usr/local/bin/deadendproxy-bsd"
echo
echo "[⚙] To autostart it, create an rc.d script or run via tmux/supervisord/etc."

# echo "[✔] DeadEndProxy fully removed." # can remove dir
