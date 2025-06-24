#!/bin/sh
# Auto-update script for DeadEndProxy
# Fetches the latest code from git, rebuilds the binary,
# and restarts the systemd service.
set -e

REPO_DIR="$(dirname "$(realpath "$0")")/.."
cd "$REPO_DIR"

echo "[*] Fetching latest code..."
if git pull --rebase --stat; then
    echo "[*] Build latest binary..."
    go build -o deadendproxy-bin ./cmd

    echo "[*] Installing binary..."
    sudo cp deadendproxy-bin /usr/local/bin/
    sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/deadendproxy-bin

    echo "[*] Restarting service..."
    sudo systemctl restart deadendproxy
    echo "[✔] DeadEndProxy updated and restarted."
else
    echo "[!] Git pull failed."
    exit 1
fi