package conf

import (
	"fmt"
	"testing"
)

func TestLoadUserConfig(t *testing.T) {
	filePath := "config.json"
	LoadUserConfig(filePath)
	fmt.Println(UserSetting)
}
