package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"gonelist/service/onedrive/auth"
)

// GetData performs a bounded request without exposing URLs, bodies or tokens in
// errors. Graph error responses must never be treated as successful operations.
func GetData(method, url string, headers map[string]string, data []byte) ([]byte, error) {
	client := auth.GetClient()
	if client == nil {
		return nil, errors.New("OneDrive client is not initialized")
	}
	return getData(context.Background(), client, method, url, headers, data)
}

func getData(ctx context.Context, client *http.Client, method, url string, headers map[string]string, data []byte) ([]byte, error) {
	attempts := 1
	if method == http.MethodGet || method == http.MethodHead {
		attempts = 3
	}
	for attempt := 0; attempt < attempts; attempt++ {
		body, retry, err := request(ctx, client, method, url, headers, data)
		if !retry || attempt == attempts-1 {
			return body, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second / 3):
		}
	}
	return nil, errors.New("request failed")
}

func request(ctx context.Context, client *http.Client, method, url string, headers map[string]string, data []byte) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(data))
	if err != nil {
		return nil, false, errors.New("invalid upstream request")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		var netErr net.Error
		retry := errors.As(err, &netErr) && netErr.Timeout()
		return nil, retry, errors.New("upstream request failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, fmt.Errorf("upstream returned HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, errors.New("could not read upstream response")
	}
	return body, false, nil
}
