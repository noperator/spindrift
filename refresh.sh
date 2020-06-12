#!/bin/bash

# Usage (cron):
# 0 * * * * /bin/bash /home/pi/spindrift/refresh.sh

# Load config.
source "$(dirname $0)/config/.env"

# Update screenshots.
node /home/pi/spindrift/screenshot-surf.js    "$SPOT"
node /home/pi/spindrift/screenshot-weather.js "$LOCATION"

# Send a command through Chrome's DevTools Protocol.
devtools() {
    echo "$@" |
    /usr/local/bin/websocat -n1 --jsonrpc $(
        curl -sk http://127.0.0.1:9222/json |
        jq '.[] | select(.url | index("report.html")) | .webSocketDebuggerUrl' -r
    )
}

# Reload page.
echo 'Reloading page.'
devtools 'Page.reload'
pkill 'node.*screenshot.js'

# Prevent showing stale data if the next refresh fails for some reason.
echo 'Clearing screenshots.'
sleep 5
find img -iname '*.png' -delete
