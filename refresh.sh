#!/bin/bash

# Usage (cron):
# 0 * * * * /bin/bash /home/pi/refresh.sh

# Send a command through Chrome's DevTools Protocol.
devtools() {
    /bin/echo "$@" |\
    /usr/local/bin/websocat -n1 --jsonrpc $(\
    /usr/bin/curl -sk http://127.0.0.1:9222/json |\
    /usr/bin/jq '.[] | select(.url | index("https://magicseaweed.com")) | .webSocketDebuggerUrl' -r \
    )
}

# Reload page.
/bin/echo 'Reloading page.'
devtools 'Page.reload'
/bin/sleep 10

# Wait for sidebar to load.
while [[ $(devtools "Runtime.evaluate {\"expression\": \"\$('div.msw-col-fixed.sidebar.msw-js-sidebar')\"}" | jq '.result | .result | .subtype' -r) == 'error' ]]; do
    /bin/echo 'Waiting for sidebar to load.'
    /bin/sleep 10
done

# Remove ads.
/bin/echo 'Removing sidebar.'
devtools "Runtime.evaluate {\"expression\": \"\$('div.msw-col-fixed.sidebar.msw-js-sidebar').remove()\"}"
