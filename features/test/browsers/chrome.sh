#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(google-chrome --version)" "Google Chrome 130.0.6723.69"

check_command_not_exists "chrome"
check_command_not_exists "firefox"
