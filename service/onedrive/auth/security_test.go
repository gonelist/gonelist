package auth

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gonelist/conf"
)

func TestTokenLogsAreRedacted(t *testing.T) {
	oldClient, oldConfig, oldState := client, oauthConfig, oauthStateString
	oldID, oldSecret := clientID, clientSecret
	out := log.StandardLogger().Out
	t.Cleanup(func() {
		client = oldClient
		oauthConfig = oldConfig
		oauthStateString = oldState
		clientID = oldID
		clientSecret = oldSecret
		log.SetOutput(out)
	})
	var logs bytes.Buffer
	log.SetOutput(&logs)
	tokenPath := filepath.Join(t.TempDir(), ".token")
	storage := &FileStorage{Path: tokenPath}
	if err := storage.SetToken(&oauth2.Token{AccessToken: "private-access-token", RefreshToken: "private-refresh-token", Expiry: time.Now().Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	SetOnedriveInfo(&conf.AllSet{
		Admin:    &conf.Admin{},
		Onedrive: &conf.Onedrive{ClientID: "test-id", ClientSecret: "private-client-secret", TokenPath: tokenPath, State: "test-state"},
	}, func() {})
	if err := GetAccessToken(ReceiveCode{Code: "private-auth-code", State: "test-state"}); err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"private-access-token", "private-refresh-token", "private-client-secret", "private-auth-code"} {
		if strings.Contains(logs.String(), secret) {
			t.Errorf("token log leaked %s", secret)
		}
	}
}
