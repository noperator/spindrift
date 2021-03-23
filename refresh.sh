#!/bin/bash

DIR=$(dirname $0)
# PATH="/home/pi/go/bin/spindrift:$PATH"

# Prevent showing stale data if the next refresh fails for some reason.
echo 'Clearing old screenshots.'
find "$DIR/img" -iname '*.png' -delete

# Pull latest screenshots.
/home/pi/go/bin/spindrift

# Reload page.
source "$DIR/devtools.sh"
echo 'Reloading page.'
devtools '{"id":1,"method":"Page.reload"}'
