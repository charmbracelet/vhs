package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"
)

const (
	maxNumber = 1000000000
	timeout   = 30 * time.Second
)

type config struct {
	Port int    `env:"VHS_PORT" envDefault:"1976"`
	Host string `env:"VHS_HOST" envDefault:""`
	Key  string `env:"VHS_KEY" envDefault:""`
	GID  int    `env:"VHS_GID" envDefault:"0"`
	UID  int    `env:"VHS_UID" envDefault:"0"`

	AuthorizedKeysPath string `env:"VHS_AUTHORIZED_KEYS_PATH" envDefault:"${HOME}/.ssh/authorized_keys" envExpand:"true"`
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the VHS SSH server",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config
		if err := env.Parse(&cfg); err != nil {
			return err
		}
		key := cfg.Key
		if key == "" {
			key = filepath.Join(".ssh", "vhs_ed25519")
		}
		if err := ensureAuthorizedKeysFile(cfg.AuthorizedKeysPath); err != nil {
			return err
		}
		s, err := wish.NewServer(
			wish.WithAddress(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
			wish.WithHostKeyPath(key),
			wish.WithAuthorizedKeys(cfg.AuthorizedKeysPath),
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

						//nolint:gosec
						rand := rand.Int63n(maxNumber)
						tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("vhs-%d.gif", rand))

						err = Evaluate(b.String(), s.Stderr(), func(v *VHS) {
							v.Options.Video.Output.GIF = tempFile
							// Disable generating MP4 & WebM.
							v.Options.Video.Output.MP4 = ""
							v.Options.Video.Output.WebM = ""
						})
						if err != nil {
							_ = s.Exit(1)
						}

						gif, _ := os.ReadFile(tempFile)
						wish.Print(s, string(gif))
						_ = os.Remove(tempFile)

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
		log.Printf("Starting SSH server on %s:%d", cfg.Host, cfg.Port)
		go func() {
			gid, uid := cfg.GID, cfg.UID
			if gid != 0 && uid != 0 {
				log.Printf("Starting server with GID: %d, UID: %d", gid, uid)
				if err := dropUserPrivileges(gid, uid); err != nil {
					log.Fatalln(err)
				}
			}
			if err = s.ListenAndServe(); err != nil {
				log.Fatalln(err)
			}
		}()

		<-done
		log.Println("Stopping SSH server")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer func() { cancel() }()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
		return nil
	},
}

func ensureAuthorizedKeysFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), 0o644); err != nil {
			return fmt.Errorf("could not create %s: %w", path, err)
		}
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			return fmt.Errorf("could not create %s: %w", path, err)
		}
	}
	return nil
}
