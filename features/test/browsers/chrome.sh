#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(google-chrome --version)" "Google Chrome 155.0.8059.39"

check_command_not_exists "chrome"
check_command_not_exists "firefox"
