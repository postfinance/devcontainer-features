#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(docker version -f '{{.Client.Version}}')" "28.3.3"
check_file_exists "/home/vscode/.docker/config.json"
cat /home/vscode/.docker/config.json | grep "######" >/dev/null 2>&1 || (echo "Custom Docker config.json has a wrong content!" && exit 1)
