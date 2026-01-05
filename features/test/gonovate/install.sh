#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_file_exists "/usr/local/bin/gonovate"
check_version "$(gonovate 2>&1 | head -n 1)" "gonovate v0.10.5"
