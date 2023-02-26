package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/keygen"
	"github.com/mattn/go-isatty"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	ghostHost = "ghost.charm.sh"
	ghostPort = 22
)

var publishCmd = &cobra.Command{
	Use:           "publish <gif>",
	Short:         "Publish your GIF to vhs.charm.sh and get a shareable URL",
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true, // we print our own errors
	RunE: func(cmd *cobra.Command, args []string) error {
		initLogger(logLevel)

		file := args[0]

		if strings.HasSuffix(file, ".tape") {
			cmd.Printf("Use vhs %s --publish flag to publish tapes\n", file)
			return errors.New("must pass a GIF file")
		}

		if !strings.HasSuffix(file, ".gif") {
			return errors.New("must pass a GIF file")
		}

		url, err := Publish(cmd.Context(), file)
		if err != nil {
			return err
		}

		if isatty.IsTerminal(os.Stdout.Fd()) {
			publishShareInstructions(url)
		}
		logger.Println(URLStyle.Render(url))
		if isatty.IsTerminal(os.Stdout.Fd()) {
			logger.Println()
		}
		return nil
	},
}

func dataPath() (string, error) {
	scope := gap.NewScope(gap.User, "vhs")
	dataPath, err := scope.DataPath("")
	if err != nil {
		return "", err
	}
	return dataPath, nil
}

// hostKeyCallback returns a callback that will be used to verify the host key.
//
// it creates a file in the given path, and uses that to verify hosts and keys.
// if the host does not exist there, it adds it so its available next time, as plain old `ssh` does.
func hostKeyCallback(path string) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		kh, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600) //nolint:gomnd
		if err != nil {
			return fmt.Errorf("failed to open known_hosts: %w", err)
		}
		defer func() { _ = kh.Close() }()

		callback, err := knownhosts.New(kh.Name())
		if err != nil {
			return fmt.Errorf("failed to check known_hosts: %w", err)
		}

		if err := callback(hostname, remote, key); err != nil {
			var kerr *knownhosts.KeyError
			if errors.As(err, &kerr) {
				if len(kerr.Want) > 0 {
					return fmt.Errorf("possible man-in-the-middle attack: %w", err)
				}
				// if want is empty, it means the host was not in the known_hosts file, so lets add it there.
				fmt.Fprintln(kh, knownhosts.Line([]string{hostname}, key))
				return nil
			}
			return fmt.Errorf("failed to check known_hosts: %w", err)
		}
		return nil
	}
}

func sshSession() (*ssh.Session, error) {
	dp, err := dataPath()
	if err != nil {
		return nil, err
	}
	kp, err := keygen.NewWithWrite(filepath.Join(dp, "vhs"), nil, keygen.Ed25519)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.NewSignerFromKey(kp.PrivateKey())
	if err != nil {
		return nil, err
	}

	pkam := ssh.PublicKeys(signer)
	sshConfig := &ssh.ClientConfig{
		User:            "vhs",
		Auth:            []ssh.AuthMethod{pkam},
		HostKeyCallback: hostKeyCallback(filepath.Join(dp, "known_hosts")),
	}

	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ghostHost, ghostPort), sshConfig)
	if err != nil {
		return nil, err
	}
	s, err := c.NewSession()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// publishShareInstructions log shareable URL
// If log level is set to `logLevelQuiet` the log message will be forced
func publishShareInstructions(url string) {
	// Set log level to verbose if current mode is quiet
	forcedLog := false
	if logLevel == logLevelQuiet {
		forcedLog = true
		setLogLevel(logLevelVerbose)
	}

	logger.Println("\n" + GrayStyle.Render("  Share your GIF with Markdown:"))
	logger.Println(CommandStyle.Render("  ![Made with VHS]") + URLStyle.Render("("+url+")"))
	logger.Println(GrayStyle.Render("\n  Or HTML (with badge):"))
	logger.Println(CommandStyle.Render("  <img ") + CommandStyle.Render("src=") + URLStyle.Render(`"`+url+`"`) + CommandStyle.Render(" alt=") + URLStyle.Render(`"Made with VHS"`) + CommandStyle.Render(">"))
	logger.Println(CommandStyle.Render("  <a ") + CommandStyle.Render("href=") + URLStyle.Render(`"https://vhs.charm.sh"`) + CommandStyle.Render(">"))
	logger.Println(CommandStyle.Render("    <img ") + CommandStyle.Render("src=") + URLStyle.Render(`"https://stuff.charm.sh/vhs/badge.svg"`) + CommandStyle.Render(">"))
	logger.Println(CommandStyle.Render("  </a>"))
	logger.Println(GrayStyle.Render("\n  Or link to it:"))
	logger.Printf("  ")

	// If log was forced restore original log level quiet
	if forcedLog {
		setLogLevel(logLevelQuiet)
	}
}

// Publish publishes the given GIF file to the web.
func Publish(ctx context.Context, path string) (string, error) {
	s, err := sshSession()
	if err != nil {
		return "", err
	}
	defer s.Close() //nolint:errcheck

	// Close connection when context is done
	go func() {
		<-ctx.Done()
		_ = s.Close()
	}()

	in, err := s.StdinPipe()
	if err != nil {
		return "", err
	}
	out, err := s.StdoutPipe()
	if err != nil {
		return "", err
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close() //nolint:errcheck

	if err := s.Start(""); err != nil {
		return "", err
	}

	_, err = io.Copy(in, f)
	if err != nil {
		return "", err
	}
	_ = in.Close()

	b, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
