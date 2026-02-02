#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(npm list -g @devcontainers/cli 2>&1)" "@devcontainers/cli"
check_version "$(npm list -g which@3.0.1 2>&1)" "which@3.0.1"
