package conf

import (
	"github.com/go-playground/assert/v2"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUserConfig(t *testing.T) {
	old := UserSet
	t.Cleanup(func() { UserSet = old })
	t.Setenv("CONF_PATH", "")
	filePath := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(filePath, []byte("server:\n  port: 8000\nlist_type: onedrive\nonedrive:\n  remote: local\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := LoadUserConfig(filePath); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, UserSet.Server.Port, 8000)

}

func TestGetBindAddr(t *testing.T) {
	assert.Equal(t, GetBindAddr(false, 8000), "127.0.0.1:8000")
	assert.Equal(t, GetBindAddr(true, 8000), ":8000")
}

func TestGetTokenPath(t *testing.T) {
	var (
		configPath string
		tokenPath  string
	)
	configPath = "/etc/gonelist/config.yml"
	tokenPath = GetTokenPath(configPath)
	assert.Equal(t, tokenPath, "/etc/gonelist/.token")

	configPath = "config.yml"
	tokenPath = GetTokenPath(configPath)
	assert.Equal(t, tokenPath, ".token")
}
