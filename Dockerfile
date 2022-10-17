FROM tsl0922/ttyd:alpine as ttyd
FROM golang:1.19.2-alpine
COPY --from=ttyd /usr/bin/ttyd /usr/bin/ttyd
WORKDIR /src/vhs
RUN apk add font-jetbrains-mono font-bitstream-vera-sans-mono-nerd \
    --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community
RUN apk add ffmpeg chromium bash
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go install cmd/vhs/vhs.go && \
    go install cmd/serve/serve.go
ENTRYPOINT ["./entrypoint.sh"]
