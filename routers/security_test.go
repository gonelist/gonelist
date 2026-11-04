package routers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/onedrive/model"
)

func testRouter(t *testing.T) (http.Handler, string, *bytes.Buffer) {
	t.Helper()
	old, writer := conf.UserSet, gin.DefaultWriter
	var logs bytes.Buffer
	gin.DefaultWriter = &logs
	base := t.TempDir()
	mount := filepath.Join(base, "mount")
	if err := os.Mkdir(mount, 0755); err != nil {
		t.Fatal(err)
	}
	conf.UserSet = &conf.AllSet{
		Server:   &conf.Server{DistPATH: filepath.Join(base, "dist")},
		Onedrive: &conf.Onedrive{Remote: "local", FolderSub: "/"},
		Local:    &conf.Local{Enable: true, Name: "local", Path: mount},
		Admin:    &conf.Admin{EnableWrite: true, Secret: "private-admin-secret"},
		WebDav:   &conf.WebDav{Enable: true, Account: "operator", Password: "private-webdav-password"},
	}
	t.Cleanup(func() { conf.UserSet = old; gin.DefaultWriter = writer })
	return InitRouter(), mount, &logs
}

func requestTo(r http.Handler, method, target string, body io.Reader, user, pass string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, body)
	if user != "" {
		req.SetBasicAuth(user, pass)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestWebDAVCredentials(t *testing.T) {
	r, _, _ := testRouter(t)
	if err := model.BatchInsertFile([]*model.FileNode{{ID: "security-test-root", Path: "/", IsFolder: true}}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { model.DeleteFile("security-test-root") })
	for _, tc := range []struct {
		user, pass string
		status     int
	}{
		{"", "", 401}, {"operator", "operator", 401}, {"operator", "wrong", 401}, {"wrong", "private-webdav-password", 401}, {"operator", "private-webdav-password", 204},
	} {
		w := requestTo(r, "OPTIONS", "/webdav/", nil, tc.user, tc.pass)
		if w.Code != tc.status {
			t.Fatalf("credentials %q: status %d want %d", tc.user, w.Code, tc.status)
		}
	}
	conf.UserSet.Admin.EnableWrite = false
	if w := requestTo(r, "MKCOL", "/webdav/new", nil, "operator", "private-webdav-password"); w.Code != 403 {
		t.Fatalf("read-only WebDAV: %d", w.Code)
	}
	if w := requestTo(r, "OPTIONS", "/webdav-other", nil, "operator", "private-webdav-password"); w.Code != 404 {
		t.Fatalf("mount boundary: %d", w.Code)
	}
}

func TestLocalRoutesConfinement(t *testing.T) {
	r, mount, _ := testRouter(t)
	parent := filepath.Dir(mount)
	for _, name := range []string{"outside", "README.md"} {
		if err := os.WriteFile(filepath.Join(parent, name), []byte("private-outside-content"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink("..", filepath.Join(mount, "escape")); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{
		"/d/local/../outside", "/d/local/%2e%2e/outside", "/d/local/escape/outside",
		"/onedrive/getpath?path=%2Flocal%2F..", "/README?path=%2Flocal%2F..",
		"/README?path=%2Flocal%2Fescape",
	} {
		w := requestTo(r, "GET", target, nil, "", "")
		if strings.Contains(w.Body.String(), "private-outside-content") || w.Code == http.StatusInternalServerError {
			t.Fatalf("unsafe response %q: %d %s", target, w.Code, w.Body.String())
		}
		if strings.HasPrefix(target, "/d/") && w.Code != 403 && w.Code != 400 {
			t.Fatalf("download not denied: %q %d", target, w.Code)
		}
	}
	target := "/onedrive/mkdir?path=/local&folder_name=" + url.QueryEscape("../created") + "&secret=private-admin-secret"
	if w := requestTo(r, "GET", target, nil, "", ""); w.Code != 403 {
		t.Fatalf("mkdir escape: %d", w.Code)
	}
	target = "/onedrive/delete_file?path=" + url.QueryEscape("/local/../outside") + "&secret=private-admin-secret"
	if w := requestTo(r, "GET", target, nil, "", ""); w.Code != 403 {
		t.Fatalf("delete escape: %d", w.Code)
	}
	var body bytes.Buffer
	mp := multipart.NewWriter(&body)
	part, err := mp.CreateFormFile("file", "hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("hello"))
	mp.Close()
	req := httptest.NewRequest("POST", "/onedrive/upload?path=/local&secret=private-admin-secret", &body)
	req.Header.Set("Content-Type", mp.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("upload: %d %s", w.Code, w.Body.String())
	}
	w = requestTo(r, "GET", "/d/local/hello.txt", nil, "", "")
	if w.Code != 200 || w.Body.String() != "hello" {
		t.Fatalf("download: %d %s", w.Code, w.Body.String())
	}
	conf.UserSet.Admin.EnableWrite = false
	w = requestTo(r, "GET", "/onedrive/delete_file?path=/local/hello.txt&secret=private-admin-secret", nil, "", "")
	var res app.Res
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Code != e.SECRET_ERROR {
		t.Fatalf("disabled write accepted: %s", w.Body.String())
	}
	if _, err := os.Stat(filepath.Join(mount, "hello.txt")); err != nil {
		t.Fatal("file removed while writes disabled")
	}
}

func TestPasswordDownloadsAndLogs(t *testing.T) {
	r, mount, logs := testRouter(t)
	if err := os.WriteFile(filepath.Join(mount, ".password"), []byte("private-password"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"/d/local/.password", "/d/private/.password", "/d/private/.PASSWORD", "/webdav/private/.password"} {
		w := requestTo(r, "GET", target, nil, "operator", "private-webdav-password")
		if w.Header().Get("Location") != "" || strings.Contains(w.Body.String(), "private-password") {
			t.Fatalf("password exposed: %s", target)
		}
		if strings.HasPrefix(target, "/d/private/") {
			var res app.Res
			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				t.Fatalf("multiple responses: %s", w.Body.String())
			}
			if res.Code != e.PASSWORD_FORBIT_DOWNLOAD {
				t.Fatalf("not blocked: %s", w.Body.String())
			}
		}
	}
	requestTo(r, "GET", "/update_permission?secret=private-admin-secret&code=private-oauth-code", nil, "", "")
	if strings.Contains(logs.String(), "private-admin-secret") || strings.Contains(logs.String(), "private-oauth-code") {
		t.Fatal("query credentials leaked into request logs")
	}
}
