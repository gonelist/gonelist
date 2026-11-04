package local

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gonelist/conf"
)

func TestMountConfinement(t *testing.T) {
	base := t.TempDir()
	mount := filepath.Join(base, "mount")
	if err := os.Mkdir(mount, 0755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(base, "outside")
	if err := os.WriteFile(outside, []byte("private"), 0600); err != nil {
		t.Fatal(err)
	}
	old := conf.UserSet
	conf.UserSet = &conf.AllSet{Local: &conf.Local{Name: "local", Path: mount, Enable: true}}
	t.Cleanup(func() { conf.UserSet = old })

	for _, p := range []string{"/local/../outside", "/locality/outside", "/local/..\\outside", "/local/.password", "/local/.PASSWORD"} {
		if f, err := Open(p); err == nil {
			f.Close()
			t.Errorf("opened forbidden path %q", p)
		}
		if err := Delete(p); err == nil {
			t.Errorf("removed forbidden path %q", p)
		}
	}
	for _, name := range []string{"../outside", "..\\outside", "/outside", ".password", ".", "..", "a/b"} {
		if err := Upload("/local", name, strings.NewReader("changed")); err == nil {
			t.Errorf("uploaded forbidden name %q", name)
		}
		if err := Mkdir("/local", name); err == nil {
			t.Errorf("created forbidden directory %q", name)
		}
	}
	if err := os.Symlink("../outside", filepath.Join(mount, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("..", filepath.Join(mount, "parent")); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"/local/link", "/local/parent/outside"} {
		if f, err := Open(p); err == nil {
			f.Close()
			t.Errorf("opened escaping symlink %q", p)
		}
	}
	if err := Upload("/local", "link", strings.NewReader("changed")); err == nil {
		t.Error("wrote through escaping symlink")
	}
	if err := Delete("/local/parent/outside"); err == nil {
		t.Error("deleted through escaping symlink")
	}
	if err := Mkdir("/local/parent", "new"); err == nil {
		t.Error("created through escaping symlink")
	}
	b, err := os.ReadFile(outside)
	if err != nil || string(b) != "private" {
		t.Fatalf("outside changed: %q %v", b, err)
	}
	if err := Mkdir("/local", "ok"); err != nil {
		t.Fatal(err)
	}
	if err := Upload("/local/ok", "hello.txt", strings.NewReader("hello")); err != nil {
		t.Fatal(err)
	}
	f, err := Open("/local/ok/hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, err = io.ReadAll(f)
	f.Close()
	if err != nil || string(b) != "hello" {
		t.Fatalf("read: %q %v", b, err)
	}
	if err := Delete("/local/ok/hello.txt"); err != nil {
		t.Fatal(err)
	}
}
