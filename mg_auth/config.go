package mg_auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

var IsLogin bool

// Config 包含了 oauth2.Config和 Storage接口
type Config struct {
	*oauth2.Config
	Storage
}

// Exchange 在接收到令牌后将其存储
func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	if err := c.Storage.SetToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

// TokenSource 可以传递已被存储的令牌
// 或当新令牌被接收时将其转换为oauth2.TokenSource
func (c *Config) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	if t == nil || !t.Valid() {
		if tok, err := c.Storage.GetToken(); err == nil {
			t = tok
		}
	}
	ts := c.Config.TokenSource(ctx, t)
	return &storageTokenSource{c, ts}
}

// Client 附加到TokenSource
func (c *Config) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx, t))
}

type storageTokenSource struct {
	*Config
	oauth2.TokenSource
}

// Token满足TokenSource接口
func (s *storageTokenSource) Token() (*oauth2.Token, error) {
	if token, err := s.Config.Storage.GetToken(); err == nil && token.Valid() {
		return token, err
	}
	token, err := s.TokenSource.Token()
	if err != nil {
		return token, err
	}
	if err := s.Config.Storage.SetToken(token); err != nil {
		return nil, err
	}
	return token, nil
}
