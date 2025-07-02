#!/bin/bash

# ========================================
# Wildcard SSL setup via Cloudflare (dns-01)
# Author: Devinsidercode CORP
# Description: Fully automated script to:
# - install Certbot with Cloudflare support
# - request wildcard SSL certificate
# - configure auto-renewal with systemd proxy restart
# ========================================

DOMAIN="eyesync.app"
WILDCARD="*.${DOMAIN}"
EMAIL="admin@${DOMAIN}"  # You can change this to your actual email

# === Step 1: Install certbot and Cloudflare DNS plugin ===
echo "✅ Installing certbot and cloudflare DNS plugin..."
apt update && apt install -y certbot python3-certbot-dns-cloudflare

# === Step 2: Prepare Cloudflare API credentials ===
echo "✅ Preparing Cloudflare credentials..."
mkdir -p ~/.secrets
chmod 700 ~/.secrets

CLOUDFLARE_TOKEN_FILE=~/.secrets/cloudflare.ini

# If token file does not exist, create a placeholder
if [ ! -f "$CLOUDFLARE_TOKEN_FILE" ]; then
    echo "dns_cloudflare_api_token = <INSERT_YOUR_TOKEN_HERE>" > "$CLOUDFLARE_TOKEN_FILE"
    echo "⚠️ Please edit ~/.secrets/cloudflare.ini and insert your actual Cloudflare API token"
    exit 1
fi

# Set strict permissions
chmod 600 "$CLOUDFLARE_TOKEN_FILE"

# === Step 3: Request wildcard SSL certificate ===
echo "✅ Requesting wildcard SSL certificate from Let's Encrypt..."
certbot certonly \
  --dns-cloudflare \
  --dns-cloudflare-credentials "$CLOUDFLARE_TOKEN_FILE" \
  -d "$WILDCARD" -d "$DOMAIN" \
  --agree-tos --non-interactive --email "$EMAIL"

# === Step 4: Configure auto-renewal with proxy restart ===
echo "✅ Configuring cron job for automatic renewal..."
CRON_CMD='0 3 * * * certbot renew --quiet --deploy-hook "systemctl restart deadendproxy"'

# Add or replace existing certbot renew task
( crontab -l 2>/dev/null | grep -v 'certbot renew' ; echo "$CRON_CMD" ) | crontab -

# === Final Output ===
echo ""
echo "🎉 All done! Wildcard SSL certificate for *.${DOMAIN} and ${DOMAIN} has been installed."
echo "📍 Certificate path: /etc/letsencrypt/live/${DOMAIN}/fullchain.pem"
echo "📍 Key path:         /etc/letsencrypt/live/${DOMAIN}/privkey.pem"
echo "🔁 Auto-renewal is scheduled daily at 03:00 server time."
echo ""
