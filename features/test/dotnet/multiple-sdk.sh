#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_command_exists "dotnet"
check_file_exists /usr/bin/dotnet
check_command_exists /usr/bin/dotnet
check_version "$(dotnet --list-sdks)" "8.0."
check_version "$(dotnet --list-sdks)" "9.0."
check_version "$(dotnet --list-runtimes)" "NETCore.App 6.0."
check_version "$(dotnet --list-runtimes)" "AspNetCore.App 8.0.15"
check_version "$(dotnet workload list)" "wasm-tools-net8"
check_version "$(dotnet workload list)" "maui-android"
