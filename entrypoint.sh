#!/bin/sh -l

TAPE="$1"

GIF="$(grep "^Output .*.gif$" $TAPE | sed 's/Output //')"
ASCII="$(grep "^Output .*.ascii$" $TAPE | sed 's/Output //')"

echo $GIF
echo $ASCII

vhs $TAPE
cat $ASCII