package onedrive

import (
	"gonelist/conf"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestCheckPassCorrect(t *testing.T) {
	var valid bool

	conf.LoadUserConfig("../example/config.json")
	InitPass(conf.UserSet)

	valid = CheckPassCorrect("/test-pass", "123456")
	assert.Equal(t, valid, true)
	valid = CheckPassCorrect("/test-pass", "12345")
	assert.Equal(t, valid, false)

	valid = CheckPassCorrect("/", "")
	assert.Equal(t, valid, true)
}

func TestGetPathArray(t *testing.T) {
	list := GetPathArray("/public/test-pass/folder")
	correctList := []string{"/", "/public", "/public/test-pass", "/public/test-pass/folder"}

	t.Log(list)
	for k := range list {
		assert.Equal(t, list[k], correctList[k])
	}

	list = GetPathArray("/")
	correctList = []string{"/"}
	for k := range list {
		assert.Equal(t, list[k], correctList[k])
	}

	list = GetPathArray("/test/")
	correctList = []string{"/", "/test"}
	for k := range list {
		assert.Equal(t, list[k], correctList[k])
	}
}
