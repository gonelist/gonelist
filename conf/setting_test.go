package conf

import (
	"testing"
)

func TestLoadUserConfig(t *testing.T) {
	filePath := "config.json"
	LoadUserConfig(filePath)
}