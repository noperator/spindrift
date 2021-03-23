#!/bin/bash

# Send a command through Chrome's DevTools Protocol.
devtools() {
    <<< "$@" /usr/local/bin/websocat -n -1 $(
        curl -sk http://127.0.0.1:9222/json |
        jq -r '.[] | select(.url | index("report.html")) | .webSocketDebuggerUrl'
    )
}
