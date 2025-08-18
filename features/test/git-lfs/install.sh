#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_version "$(git lfs version)" "git-lfs/3.7.0 (GitHub; linux amd64; go 1.24.4; git 92dddf56)"
