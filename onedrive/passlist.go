package onedrive

import (
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"strings"
)

// 用户设置的 目录密码
var passList = map[string]string{}

func InitPass(user *conf.UserSetting) {
	if user.PassList == nil {
		return
	}

	for _, pass := range user.PassList {
		if pass.Path != "" {
			passList[pass.Path] = pass.Pass
		}
	}
	log.WithField("passlist", passList).Info("成功导入目录密码")
}

// 判断输入目录和密码是否正确
func CheckPassCorrect(path, pass string) bool {
	isCorrect := true
	if pa, ok := passList[path]; ok && pa == pass {
		return true
	}
	// 对输入路径进行拆分
	list := GetPathArray(path)
	// 判断每一个路径
	for k := range list {
		if !CheckSinglePath(list[k], pass) {
			isCorrect = false
		}
	}
	return isCorrect
}

func CheckSinglePath(p, pass string) bool {
	pa, ok := passList[p]
	if !ok {
		return true
	}
	if pa == pass {
		return true
	}
	return false
}

func GetPathArray(path string) []string {
	list := strings.Split(path, "/")

	current := ""
	reList := []string{"/"}
	for i := 0; i < len(list); i++ {
		if list[i] == "" {
			continue
		}
		current += "/" + list[i]
		reList = append(reList, current)
	}

	return reList
}
