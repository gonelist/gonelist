package file

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestFatherPath(t *testing.T) {

	assert.Equal(t, FatherPath("/"), "/")
	assert.Equal(t, FatherPath("/public"), "/")
	assert.Equal(t, FatherPath("/public/test"), "/public/")

}
