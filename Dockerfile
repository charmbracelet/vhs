FROM tsl0922/ttyd:alpine as ttyd
FROM golang:1.19.2-alpine

COPY --from=ttyd /usr/bin/ttyd /usr/bin/ttyd
WORKDIR /src/vhs

# Install Fonts
RUN apk add \
    font-adobe-source-code-pro font-source-code-pro-nerd \
    font-bitstream-vera-sans-mono-nerd \
    font-dejavu font-dejavu-sans-mono-nerd \
    font-fira-code font-fira-code-nerd \
    font-hack font-hack-nerd \
    font-ibm-plex-mono-nerd \
    font-inconsolata font-inconsolata-nerd \
    font-inconsolata font-inconsolata-nerd \
    font-jetbrains-mono font-jetbrains-mono-nerd \
    font-liberation font-liberation-mono-nerd \
    font-noto \
    font-roboto-mono \
    font-ubuntu font-ubuntu-mono-nerd \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing
        
# Install Dependencies
RUN apk add ffmpeg chromium bash

# Verify Go Dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Source + Install
COPY . .
RUN go install cmd/vhs/vhs.go && \
    go install cmd/serve/serve.go

ENTRYPOINT ["./entrypoint.sh"]
