package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/vhs"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
)

const (
	host = "localhost"
	port = 1976
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			func(h ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					// Request for vhs must be passed in through stdin, which
					// implies that there is no PTY.
					//
					// In the future, we should support PTY by providing a
					// Bubble Tea interface for VHS.
					//
					// Ideally, users can SSH into the server and get a
					// walk through of how to write a .tape file.
					_, _, isPty := s.Pty()
					if isPty {
						wish.Println(s, "PTY is not supported")
						_ = s.Exit(1)
						return
					}

					// Read stdin passed from the client.
					// This is the .tape file which contains the VHS commands.
					//
					// ssh vhs.charm.sh < demo.tape
					var b bytes.Buffer
					_, err := io.Copy(&b, s)
					if err != nil {
						wish.Errorln(s, err)
						_ = s.Exit(1)
						return
					}

					err = vhs.Evaluate(b.String(), s)
					if err != nil {
						_ = s.Exit(1)
					}

					h(s)
				}
			},
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
