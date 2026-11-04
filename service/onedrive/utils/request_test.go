package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error { b.closed = true; return nil }

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) { return 0, errors.New("private-response") }

func TestResponseHandling(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		reader  io.Reader
		wantErr bool
	}{
		{"success", 200, strings.NewReader("ok"), false},
		{"unauthorized", 401, strings.NewReader("private-response"), true},
		{"not-found", 404, strings.NewReader("private-response"), true},
		{"rate-limit", 429, strings.NewReader("private-response"), true},
		{"server-error", 500, strings.NewReader("private-response"), true},
		{"read-error", 200, brokenReader{}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &trackedBody{Reader: tc.reader}
			client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if _, ok := r.Context().Deadline(); !ok {
					t.Error("missing request deadline")
				}
				return &http.Response{StatusCode: tc.status, Body: body, Header: make(http.Header)}, nil
			})}
			_, err := getData(context.Background(), client, "GET", "https://example.invalid/?token=private-url", nil, nil)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error: %v", err)
			}
			if !body.closed {
				t.Error("response body was not closed")
			}
			if err != nil && strings.Contains(err.Error(), "private") {
				t.Error("error leaked credentials")
			}
		})
	}
}

func TestInvalidRequest(t *testing.T) {
	called := false
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		called = true
		return nil, errors.New("unexpected request")
	})}
	if _, err := getData(context.Background(), client, "GET", "://invalid", map[string]string{"X-Test": "1"}, nil); err == nil {
		t.Fatal("invalid URL accepted")
	}
	if called {
		t.Error("sent invalid request")
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "private-url timed out" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestWriteTimeoutIsNotReplayed(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		body, err := io.ReadAll(r.Body)
		if err != nil || string(body) != "payload" {
			t.Fatalf("request body: %q %v", body, err)
		}
		r.Body.Close()
		return nil, timeoutError{}
	})}
	_, err := getData(context.Background(), client, "PUT", "https://example.invalid/", nil, []byte("payload"))
	if err == nil || calls != 1 {
		t.Fatalf("write was retried: %d %v", calls, err)
	}
	if strings.Contains(err.Error(), "private-url") {
		t.Fatal("transport error leaked URL")
	}
}
