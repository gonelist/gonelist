package onedrive

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"golang.org/x/oauth2"
)

//参考 https://www.kancloud.cn/mutouzhang/gocookbook/702034

// Storage 是我们的通用存储接口
type Storage interface {
	GetToken() (*oauth2.Token, error)
	SetToken(*oauth2.Token) error
}

// GetToken 检索github oauth2令牌
func GetToken(ctx context.Context, conf Config, code string) (*oauth2.Token, error) {
	token, err := conf.Storage.GetToken()
	if err == nil && token.Valid() {
		return token, err
	}

	return conf.Exchange(ctx, code)
}

// FileStorage 满足storage 接口
type FileStorage struct {
	Path string
	mu   sync.RWMutex
}

// GetToken 从文件中检索令牌
func (f *FileStorage) GetToken() (*oauth2.Token, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	in, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	var t *oauth2.Token
	data := json.NewDecoder(in)
	return t, data.Decode(&t)
}

// SetToken 将令牌存储在文件中
func (f *FileStorage) SetToken(t *oauth2.Token) error {
	if t == nil || !t.Valid() {
		return errors.New("bad token")
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	out, err := os.OpenFile(f.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	data, err := json.Marshal(&t)
	if err != nil {
		return err
	}

	_, err = out.Write(data)
	return err
}
