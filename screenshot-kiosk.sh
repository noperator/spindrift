#!/bin/bash

# Send a command through Chrome's DevTools Protocol.
devtools() {
    echo "$@" |
    /usr/local/bin/websocat -B 2000000 -n1 --jsonrpc $(
        curl -sk http://127.0.0.1:9222/json |
        jq '.[] | select(.url | index("report.html")) | .webSocketDebuggerUrl' -r
    )
}

# Extract and decode PNG screenshot.
SCREENSHOT="screenshot-$(date +%s).png"
echo "Saving screenshot to $SCREENSHOT"
devtools 'Page.captureScreenshot' | jq -r '.result | .data' | base64 -d > "$SCREENSHOT"
