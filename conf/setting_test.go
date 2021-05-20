package conf

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestLoadUserConfig(t *testing.T) {
	filePath := "../config.yml"
	LoadUserConfig(filePath)
	if UserSet == nil || UserSet.Server == nil {
		log.Error("not find Userset")
	}
	assert.Equal(t, UserSet.Server.Port, 8000)
	log.Infof("default config.yml: %v", UserSet)
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
