#!/bin/bash

DIR=$(dirname $0)

# Prevent showing stale data if the next refresh fails for some reason.
echo 'Clearing old screenshots.'
find "$DIR/img" -iname '*.png' -delete

# Pull latest screenshots.
"$HOME/go/bin/check-forecast"

# Reload page.
source "$DIR/devtools.sh"
echo 'Reloading page.'
devtools '{"id":1,"method":"Page.reload"}'
