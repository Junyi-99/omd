#!/bin/bash
# Shared window/session dialog; resolve symlinks through the theme entry point.
config_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec python3 "$config_dir/new-item-prompt.py" "${1:-window}"
