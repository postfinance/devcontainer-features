#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(locale -k LC_TIME)" "abday=\"Sun;Mon;Tue;Wed;Thu;Fri;Sat\""

check_version "$(env)" "LANG=en_US.UTF-8"
check_version "$(env)" "LANGUAGE=en_US.UTF-8"
check_version "$(env)" "LC_ALL=en_US.UTF-8"
