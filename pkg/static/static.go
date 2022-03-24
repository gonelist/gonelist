package static

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func DownloadStatic(version string) error {
	url := fmt.Sprintf("https://github.do/https://github.com/gonelist/gonelist-web/releases/download/%v/dist.tar.gz", version)
	log.Infoln("开始下载文件")
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile("dist.tar.gz", data, 0755)
	if err != nil {
		return err
	}
	err = DeCompress("dist.tar.gz", "./")
	if err != nil {
		return err
	}
	return nil
}

//解压 tar.gz
func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		if hdr.FileInfo().IsDir() {
			_ = os.Mkdir(filename, 0666)
			continue
		}

		file, err := createFile(filename)
		if err != nil {
			return err
		}
		_, _ = io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
