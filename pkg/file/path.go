package file

import (
	"os"
	"strings"
)

// 替换去掉 Sub 路径
// 如果设置了 folderSub 为 /public
// 那么 /public 替换为 /, /public/test 替换为 /test
func RemoveSubPath(pSrc, foldSub string) string {
	if foldSub != "/" {
		pSrc = strings.Replace(pSrc, foldSub, "", 1)
	}
	return pSrc
}

// 获取上一级目录
// 如 / 返回 /, /public 返回 /, /public/test 返回 /public/
func FatherPath(pSrc string) string {
	if pSrc == "/" || len(pSrc) == 0 {
		return "/"
	}

	index := strings.LastIndex(pSrc, string(os.PathSeparator))
	return pSrc[:index+1]
}
