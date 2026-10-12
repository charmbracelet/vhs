package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProbeDevTools(t *testing.T) {
	t.Run("ready endpoint", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		_, port, err := net.SplitHostPort(ts.Listener.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		if err := probeDevTools(mustAtoi(t, port)); err != nil {
			t.Errorf("expected the endpoint to be ready, got: %v", err)
		}
	})

	t.Run("error status code", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		_, port, err := net.SplitHostPort(ts.Listener.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		if err := probeDevTools(mustAtoi(t, port)); err == nil {
			t.Error("expected an error for a not ready endpoint")
		}
	})

	t.Run("nothing listening", func(t *testing.T) {
		addr, err := net.Listen("tcp", "127.0.0.1:0") //nolint:gosec
		if err != nil {
			t.Fatal(err)
		}
		port := addr.Addr().(*net.TCPAddr).Port
		if err := addr.Close(); err != nil {
			t.Fatal(err)
		}

		if err := probeDevTools(port); err == nil {
			t.Error("expected an error when nothing is listening")
		}
	})
}

func TestWaitForDevTools(t *testing.T) {
	t.Run("becomes ready", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		_, port, err := net.SplitHostPort(ts.Listener.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		if !waitForDevTools(context.Background(), mustAtoi(t, port), nil, time.Second) {
			t.Error("expected the endpoint to become ready")
		}
	})

	t.Run("times out when nothing is listening", func(t *testing.T) {
		addr, err := net.Listen("tcp", "127.0.0.1:0") //nolint:gosec
		if err != nil {
			t.Fatal(err)
		}
		port := addr.Addr().(*net.TCPAddr).Port
		if err := addr.Close(); err != nil {
			t.Fatal(err)
		}

		if waitForDevTools(context.Background(), port, nil, 250*time.Millisecond) {
			t.Error("expected waitForDevTools to give up")
		}
	})

	t.Run("gives up when the browser exits", func(t *testing.T) {
		exited := make(chan struct{})
		close(exited)

		if waitForDevTools(context.Background(), 1, exited, time.Second) {
			t.Error("expected waitForDevTools to give up on exit")
		}
	})

	t.Run("respects the context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if waitForDevTools(ctx, 1, nil, time.Second) {
			t.Error("expected waitForDevTools to give up on a canceled context")
		}
	})
}

func mustAtoi(tb testing.TB, s string) int {
	tb.Helper()
	port, err := net.LookupPort("tcp", s)
	if err != nil {
		tb.Fatal(err)
	}
	return port
}
