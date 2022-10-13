#!/bin/sh -l

TAPE="$1"

GIF="$(grep "^Output .*.gif$" $TAPE | sed 's/Output //')"


vhs $TAPE

ASCII_OUTPUT="$(cat "$(grep "^Output .*.ascii$" $TAPE | sed 's/Output //')")"

echo "::set-output name=ascii::$ASCII_OUTPUT"
echo "::set-output name=gif::https://stuff.charm.sh/bubbletea-examples/chat.gif"
