REV = $(shell git rev-parse --short=7 HEAD)
make:
	go install -ldflags "-s -w -X=main.Version=$(REV)" cmd/vhs/vhs.go

man:
	pandoc --standalone --to man docs/vhs.1.md -o docs/vhs.1
