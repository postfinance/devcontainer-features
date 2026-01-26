#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(rustup --version 2>/dev/null)" "rustup 1.27.1 (54dd3d00f 2024-04-24)"
check_version "$(rustc --version)" "rustc 1.76.0 (07dca489a 2024-02-04)"
check_version "$(rustup target list --installed)" $'x86_64-pc-windows-gnu\nx86_64-unknown-linux-gnu'
