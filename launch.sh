#!/bin/bash

# Usage: 
# startx ./launch.sh -- -nocursor

# Disable screensaver.
xset s off
xset -dpms
xset s noblank

# Launch Chrome in kiosk mode. Note: Chromium's CLI switches are poorly documented; there's a good reference at https://peter.sh/experiments/chromium-command-line-switches/ 
chromium-browser "$(dirname $0)/report.html" \
--disable-component-update \
--disable-gpu \
--disable-infobars \
--disable-translate \
--disk-cache-dir=/dev/null \
--incognito \
--kiosk \
--no-first-run \
--noerrdialogs \
--remote-debugging-port=9222 \
--simulate-outdated-no-au='Tue, 31 Dec 2099 23:59:59 GMT' \
--window-size="$1,$2"
