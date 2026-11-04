// Package local
// 挂载本地目录
package local

import (
	"io"
	"os"
	"strings"
	"time"

	"gonelist/conf"
	"gonelist/service/onedrive/model"
)

// IsPath matches the mount itself or a child, never a similarly prefixed mount.
func IsPath(p string) bool {
	mount := "/" + conf.UserSet.Local.Name
	return p == mount || strings.HasPrefix(p, mount+"/")
}

// openRoot keeps resolution and file operations inside the mount, including
// concurrent symlink replacement. Never turn this back into an absolute path.
func openRoot(p string) (*os.Root, string, error) {
	if !IsPath(p) || strings.ContainsAny(p, "\\\x00") {
		return nil, "", os.ErrPermission
	}
	name := strings.TrimPrefix(strings.TrimPrefix(p, "/"+conf.UserSet.Local.Name), "/")
	for _, part := range strings.Split(name, "/") {
		if part == ".." || strings.EqualFold(part, ".password") {
			return nil, "", os.ErrPermission
		}
	}
	if name == "" {
		name = "."
	}
	root, err := os.OpenRoot(conf.UserSet.Local.Path)
	return root, name, err
}

func Open(p string) (*os.File, error) {
	root, name, err := openRoot(p)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return root.Open(name)
}

func childPath(p, name string) (string, error) {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\\x00:") {
		return "", os.ErrPermission
	}
	return strings.TrimRight(p, "/") + "/" + name, nil
}

func Mkdir(p, name string) error {
	child, err := childPath(p, name)
	if err != nil {
		return err
	}
	root, rel, err := openRoot(child)
	if err != nil {
		return err
	}
	defer root.Close()
	return root.Mkdir(rel, 0755)
}

func Upload(p, name string, src io.Reader) error {
	child, err := childPath(p, name)
	if err != nil {
		return err
	}
	root, rel, err := openRoot(child)
	if err != nil {
		return err
	}
	defer root.Close()
	dst, err := root.OpenFile(rel, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func Delete(p string) error {
	root, rel, err := openRoot(p)
	if err != nil {
		return err
	}
	defer root.Close()
	if rel == "." {
		return os.ErrPermission
	}
	return root.Remove(rel)
}

func GetPath(path string) ([]*model.FileNode, error) {
	dir, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	entries, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	var nodes []*model.FileNode
	for _, file := range entries {
		name := file.Name()
		if strings.EqualFold(name, ".password") {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &model.FileNode{
			ID:             "",
			Name:           name,
			Path:           path + "/" + name,
			READMEUrl:      "",
			IsFolder:       file.IsDir(),
			DownloadURL:    "",
			LastModifyTime: info.ModTime(),
			Size:           info.Size(),
			Children:       nil,
			RefreshTime:    time.Time{},
			Password:       "",
			PasswordURL:    "",
			ParentID:       "",
		})
	}
	return nodes, err
}
