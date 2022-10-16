REV = $(shell git rev-parse --short=7 HEAD)
make:
	go install -ldflags "-s -w -X=main.Version=$(REV)" cmd/vhs/vhs.go

man:
	pandoc --standalone --to man vhs.1.md -o vhs.1
