#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

# Detect distro and codename
. /etc/os-release

if { [[ "$ID" = "debian" && "$VERSION_CODENAME" = "trixie" ]] || [[ "$ID" = "ubuntu" && "$VERSION_CODENAME" = "noble" ]]; }; then
    check_package_installed "libgtk2.0-0t64"
    check_package_installed "libgtk-3-0t64"
    check_package_installed "libasound2t64"
else
    check_package_installed "libgtk2.0-0"
    check_package_installed "libgtk-3-0"
    check_package_installed "libasound2"
fi

check_package_installed "libgbm-dev"
check_package_installed "libnotify-dev"
check_package_installed "libnss3"
check_package_installed "libxss1"
check_package_installed "libxtst6"
check_package_installed "xauth"
check_package_installed "xvfb"
