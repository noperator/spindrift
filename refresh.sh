#!/bin/bash

# Usage (cron):
# 0 * * * * /bin/bash /home/pi/spindrift/refresh.sh

# Load config.
source "$(dirname $0)/.env"

# Update screenshots.
node /home/pi/spindrift/screenshot.js "$SPOT"

# Send a command through Chrome's DevTools Protocol.
devtools() {
    /bin/echo "$@" |\
    /usr/local/bin/websocat -n1 --jsonrpc $(\
    /usr/bin/curl -sk http://127.0.0.1:9222/json |\
    /usr/bin/jq '.[] | select(.url | index("report.html")) | .webSocketDebuggerUrl' -r \
    )
}

# Reload page.
/bin/echo 'Reloading page.'
devtools 'Page.reload'
