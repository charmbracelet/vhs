#!/bin/sh -l

TAPE="$1"


vhs $TAPE

ASCII="$(cat "$(grep "^Output .*.ascii$" $TAPE | sed 's/Output //')")"
GIF="$(base64 -w 0 "$(grep "^Output .*.gif$" $TAPE | sed 's/Output //')")"

echo "::set-output name=ascii::$ASCII"
echo "::set-output name=gif::$GIF"
