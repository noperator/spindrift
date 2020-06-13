#!/bin/bash

# Send a command through Chrome's DevTools Protocol.
devtools() {
    <<< "$@" /usr/local/bin/websocat -B 2000000 -n1 --jsonrpc $(
        curl -sk http://127.0.0.1:9222/json |
        jq '.[] | select(.url | index("report.html")) | .webSocketDebuggerUrl' -r
    )
}
