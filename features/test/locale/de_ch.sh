#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(locale -k LC_TIME)" "abday=\"So;Mo;Di;Mi;Do;Fr;Sa\""

check_version "$(env)" "LANG=de_CH.UTF-8"
check_version "$(env)" "LANGUAGE=de_CH.UTF-8"
check_version "$(env)" "LC_ALL=de_CH.UTF-8"
