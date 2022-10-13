FROM golang:1.19.2-alpine
WORKDIR /src/vhs
RUN apk add ffmpeg ttyd chromium bash
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go install cmd/vhs/vhs.go && \
    go install cmd/serve/serve.go
ENTRYPOINT ["entrypoint.sh"]
