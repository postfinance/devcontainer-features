#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_file_exists "/usr/local/goreleaser/goreleaser"
check_version "$(goreleaser -v)" "2.4.8"
