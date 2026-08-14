#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

# Detect distro and codename
. /etc/os-release

if { [[ "$ID" = "debian" && "$VERSION_CODENAME" = "bookworm" ]] }; then
    check_package_installed "libwebkit2gtk-4.0-37"
else
    check_package_installed "libwebkit2gtk-4.1-0"
fi

check_package_installed "libswt-gtk-4-jni"
