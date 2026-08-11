#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(docker version -f '{{.Client.Version}}')" "27.2.1"

# Ensure the daemon is reachable and returns server metadata.
if ! docker info -f '{{.ServerVersion}}' >/dev/null 2>&1; then
	echo "dockerd is not running or not responding"
	exit 1
fi
