#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_command_exists uv
grep -q 'index-url = "https://pypi.org/simple"' /etc/uv/uv.toml
grep -q 'extra-index-url = \["https://example.com/simple"\]' /etc/uv/uv.toml
