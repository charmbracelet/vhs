#!/bin/sh

curl $1 |
	sed  \
		-e 's/"purple":/"magenta":/g' \
		-e 's/"brightPurple":/"brightMagenta":/g' \
		-e 's/selectionBackground/selection/g' \
		-e 's/cursorColor/cursor/g' |
	jq >$2.json
