package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

type publishResponse struct {
	output     string
	status     uint32
	omitStatus bool
	rejectExec bool
	stall      bool
}

func newPublishTestSession(t *testing.T, response publishResponse) (*ssh.Session, <-chan []byte) {
	t.Helper()
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		t.Fatal(err)
	}
	config := &ssh.ServerConfig{NoClientAuth: true}
	config.AddHostKey(signer)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	received := make(chan []byte, 1)
	release := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		defer connection.Close()
		server, channels, requests, err := ssh.NewServerConn(connection, config)
		if err != nil {
			return
		}
		defer server.Close()
		go ssh.DiscardRequests(requests)
		for incoming := range channels {
			channel, requests, err := incoming.Accept()
			if err != nil {
				return
			}
			for request := range requests {
				if request.Type != "exec" || response.rejectExec {
					_ = request.Reply(false, nil)
					_ = channel.Close()
					continue
				}
				_ = request.Reply(true, nil)
				body, err := io.ReadAll(channel)
				if err != nil {
					_ = channel.Close()
					return
				}
				received <- body
				if response.stall {
					<-release
				}
				_, _ = io.WriteString(channel, response.output)
				if response.status != 0 {
					_, _ = io.WriteString(channel.Stderr(), "file is too large\n")
				}
				if !response.omitStatus {
					_, _ = channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{response.status}))
				}
				_ = channel.Close()
			}
		}
	}()
	client, err := ssh.Dial("tcp", listener.Addr().String(), &ssh.ClientConfig{
		User:            "fixture",
		HostKeyCallback: ssh.FixedHostKey(signer.PublicKey()),
		Timeout:         5 * time.Second,
	})
	if err != nil {
		_ = listener.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		close(release)
		_ = client.Close()
		_ = listener.Close()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("SSH fixture did not exit")
		}
	})
	session, err := client.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	return session, received
}

func publishTestGIF(t *testing.T) (string, []byte) {
	t.Helper()
	body, err := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "fixture.gif")
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	return path, body
}

func TestPublishSession(t *testing.T) {
	const url = "https://example.invalid/recording.gif"
	for _, test := range []struct {
		name      string
		response  publishResponse
		wantURL   string
		wantError bool
	}{
		{"success", publishResponse{output: url}, url, false},
		{"trailing newline", publishResponse{output: url + "\r\n"}, url, false},
		{"upload rejected", publishResponse{status: 1}, "", true},
		{"failure with output", publishResponse{output: url, status: 42}, "", true},
		{"empty successful response", publishResponse{}, "", true},
		{"whitespace response", publishResponse{output: " \r\n\t"}, "", true},
		{"missing exit status", publishResponse{output: url, omitStatus: true}, "", true},
		{"exec rejected", publishResponse{rejectExec: true}, "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			path, body := publishTestGIF(t)
			session, received := newPublishTestSession(t, test.response)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			got, err := publishSession(ctx, path, session)
			if (err != nil) != test.wantError {
				t.Errorf("error=%v, wantError=%v", err, test.wantError)
			}
			if got != test.wantURL {
				t.Errorf("URL=%q, want=%q", got, test.wantURL)
			}
			if test.response.status != 0 {
				var exitError *ssh.ExitError
				if !errors.As(err, &exitError) || exitError.ExitStatus() != int(test.response.status) {
					t.Errorf("remote exit status was lost: %v", err)
				}
			}
			if !test.response.rejectExec {
				select {
				case data := <-received:
					if !bytes.Equal(data, body) {
						t.Error("uploaded GIF bytes changed")
					}
				case <-ctx.Done():
					t.Error("server did not receive the upload")
				}
			}
		})
	}
}

func TestPublishSessionMissingFile(t *testing.T) {
	session, _ := newPublishTestSession(t, publishResponse{})
	_, err := publishSession(context.Background(), filepath.Join(t.TempDir(), "missing.gif"), session)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected missing-file error, got %v", err)
	}
}

func TestPublishSessionCancellation(t *testing.T) {
	path, _ := publishTestGIF(t)
	session, received := newPublishTestSession(t, publishResponse{stall: true})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := publishSession(ctx, path, session)
		done <- err
	}()
	select {
	case <-received:
	case <-time.After(5 * time.Second):
		t.Fatal("server did not receive the upload")
	}
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Error("canceled publish reported success")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("canceled publish did not return")
	}
}
