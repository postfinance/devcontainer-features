#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(docker version -f '{{.Client.Version}}')" "28.3.3"
check_version "$(docker compose version)" "Docker Compose version v2.39.1"
check_version "$(docker buildx version)" "github.com/docker/buildx v0.21.2 1360a9e8d25a2c3d03c2776d53ae62e6ff0a843d"
