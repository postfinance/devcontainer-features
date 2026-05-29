#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

check_command_exists java
check_command_exists gradle
check_env_var_exists JAVA_HOME
check_env_var_exists GRADLE_HOME
check_version "$(java -version 2>&1)" "openjdk"
check_version "$(gradle --version)" "Gradle"
