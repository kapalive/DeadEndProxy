#!/bin/sh
# © 2023 Devinsidercode CORP. Licensed under the MIT License.
CONFIG_PATH="/etc/deadendproxy/config.yaml"
BINARY="/usr/local/bin/deadendproxy-bin"

if [ "$1" = "config" ]; then
    editor="${EDITOR:-vi}"
    "$editor" "$CONFIG_PATH"
else
    exec "$BINARY" "$@"
fi
