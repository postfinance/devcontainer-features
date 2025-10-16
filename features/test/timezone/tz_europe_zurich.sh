#!/bin/bash
set -e

[[ -f "$(dirname "$0")/../functions.sh" ]] && source "$(dirname "$0")/../functions.sh"
[[ -f "$(dirname "$0")/functions.sh" ]] && source "$(dirname "$0")/functions.sh"

timezone="$(date +%Z)"
if [ "$timezone" != "CET" ] && [ "$timezone" != "CEST" ]; then
    echo "Timezone '$timezone' is incorrect"
    exit 1
fi

check_version "$(cat /etc/timezone)" "Europe/Zurich"
