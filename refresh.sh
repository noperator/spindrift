#!/bin/bash

# Usage (cron):
# 0 * * * * /bin/bash /home/pi/spindrift/refresh.sh

# Load config.
source "$(dirname $0)/.env"

# Update screenshots.
node /home/pi/spindrift/screenshot.js "$SPOT"

# Send a command through Chrome's DevTools Protocol.
devtools() {
    echo "$@" |
    websocat -n1 --jsonrpc $(
        curl -sk http://127.0.0.1:9222/json |
        jq '.[] | select(.url | index("report.html")) | .webSocketDebuggerUrl' -r
    )
}

# Reload page.
echo 'Reloading page.'
devtools 'Page.reload'
