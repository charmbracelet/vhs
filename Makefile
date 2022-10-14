REV = $(shell git rev-parse --short=7 HEAD)
make:
	go install -ldflags "-s -w -X=main.Version=$(REV)" cmd/vhs/vhs.go
