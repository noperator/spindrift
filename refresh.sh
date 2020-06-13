#!/bin/bash

DIR=$(dirname $0)

# Prevent showing stale data if the next refresh fails for some reason.
echo 'Clearing old screenshots.'
find "$DIR/img" -iname '*.png' -delete

# Load config and pull latest screenshots.
source "$DIR/config/.env"
node /home/pi/spindrift/screenshot-surf.js    "$SPOT"
node /home/pi/spindrift/screenshot-weather.js "$LOCATION"

# Reload page.
source "$DIR/devtools.sh"
echo 'Reloading page.'
devtools 'Page.reload'
pkill 'node.*screenshot.js'
