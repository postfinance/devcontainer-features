#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_file_exists "/usr/local/bin/apm"
check_version "$(apm --version)" "0.26.0"
