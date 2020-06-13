#!/bin/bash

DIR=$(dirname $0)

# Extract and decode PNG screenshot.
source "$DIR/devtools.sh"
SCREENSHOT="$DIR/screenshot-$(date +%s).png"
echo "Saving screenshot to $SCREENSHOT"
devtools 'Page.captureScreenshot' | jq -r '.result | .data' | base64 -d > "$SCREENSHOT"
