#!/bin/bash

DIR=$(dirname $0)

DIMENSIONS=$(
    for DIM in width height; do
        awk -v "dim=$DIM" '$1 == dim {printf $NF ","}' "$DIR/config/params.toml" |
        sed 's/,$//'
    done
)

startx "$DIR/launch.sh" "$DIMENSIONS" -- -nocursor &
