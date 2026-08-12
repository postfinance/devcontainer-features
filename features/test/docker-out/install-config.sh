#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

NON_ROOT_USER=$(resolve_non_root_user)
check_version "$(docker version -f '{{.Client.Version}}')" "28.3.3"
check_file_exists "/home/${NON_ROOT_USER}/.docker/config.json"
cat "/home/${NON_ROOT_USER}/.docker/config.json" | grep "######" >/dev/null 2>&1 || (echo "Custom Docker config.json has a wrong content!" && exit 1)
