#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_command_exists "dotnet"
check_version "$(dotnet --version | head -1)" "8.0."
check_file_exists /usr/bin/dotnet
check_command_exists /usr/bin/dotnet
