#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(node -v)" "v20.11.1"
check_version "$(npm -v)" "10.5.0"
check_version "$(yarn -v)" "1.22.22"
check_version "$(pnpm -v)" "9.14.2"
