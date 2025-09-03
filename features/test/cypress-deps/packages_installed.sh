#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_package_installed "libgtk2.0-0"
check_package_installed "libgtk-3-0"
check_package_installed "libgbm-dev"
check_package_installed "libnotify-dev"
check_package_installed "libnss3"
check_package_installed "libxss1"
check_package_installed "libasound2"
check_package_installed "libxtst6"
check_package_installed "xauth"
check_package_installed "xvfb"
