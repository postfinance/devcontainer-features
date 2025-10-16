#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_command_exists "sonar-scanner"
check_dir_not_exists "/usr/local/sonar-scanner/jre"
check_version "$(sonar-scanner --version | sed -n 3p | cut -b 20-)" "SonarScanner CLI 7.3.0.5189"