#!/bin/sh

curl $1 |
	sed  \
		-e 's/"purple":/"magenta":/g' \
		-e 's/"brightPurple":/"brightMagenta":/g' \
		-e 's/selectionBackground/selection/g' \
		-e 's/cursorColor/cursor/g' |
	jq 'map(del(.meta.credits,.cursor,.selection))' >$2.json

