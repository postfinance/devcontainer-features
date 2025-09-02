#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(node -v)" "v24.4.1"
check_version "$(corepack -v)" "0.34.0"

corepack prepare yarn@4.9.2
corepack prepare pnpm@10.14.0
