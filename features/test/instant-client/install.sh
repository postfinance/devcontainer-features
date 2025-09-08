#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_file_exists "/opt/oracle/instantclient_23_9/BASIC_README"
check_version "$(ldconfig -v)" "/opt/oracle/instantclient_23_9"
