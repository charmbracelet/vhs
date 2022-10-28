#!/bin/sh

curl https://raw.githubusercontent.com/atomcorp/themes/master/app/src/$1.json |
	sed  \
		-e 's/"purple":/"magenta":/g' \
		-e 's/"brightPurple":/"brightMagenta":/g' \
		-e 's/selectionBackground/selection/g' \
		-e 's/cursorColor/cursor/g' |
	jq -c >$2.json

