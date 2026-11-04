package conf

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestSensitiveConfiguration(t *testing.T) {
	old, out := UserSet, log.StandardLogger().Out
	t.Cleanup(func() { UserSet = old; log.SetOutput(out) })
	t.Setenv("CONF_PATH", "")
	var logs bytes.Buffer
	log.SetOutput(&logs)
	file := filepath.Join(t.TempDir(), "config.yml")
	base := "server:\n  port: 8000\nlist_type: onedrive\nonedrive:\n  remote: onedrive\n  client_id: custom-id\n  client_secret: private-client-secret\n  pass_list:\n    - path: /private\n      pass: private-folder-pass\n"
	for _, tc := range []struct {
		name, extra string
		bad         bool
	}{
		{"configured", "admin:\n  enable_write: true\n  secret: private-admin-secret\n", false},
		{"blank-secret", "admin:\n  enable_write: true\n", true},
		{"old-default", "admin:\n  enable_write: true\n  secret: 123465\n", true},
		{"empty-webdav", "web_dav:\n  enable: true\n  account: admin\n", true},
		{"valid-webdav", "web_dav:\n  enable: true\n  account: admin\n  password: private-webdav-password\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(file, []byte(base+tc.extra), 0600); err != nil {
				t.Fatal(err)
			}
			err := LoadUserConfig(file)
			if (err != nil) != tc.bad {
				t.Fatalf("load error: %v", err)
			}
		})
	}
	for _, secret := range []string{"private-client-secret", "private-folder-pass", "private-admin-secret", "private-webdav-password"} {
		if strings.Contains(logs.String(), secret) {
			t.Fatalf("configuration log leaked %s", secret)
		}
	}
}
