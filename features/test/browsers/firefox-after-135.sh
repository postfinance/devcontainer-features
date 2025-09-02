#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(firefox --version)" "Mozilla Firefox 135.0"

check_command_not_exists "google-chrome"
check_command_not_exists "chrome"
