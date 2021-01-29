package util_test

import (
	"github.com/stretchr/testify/assert"
	"send-fiule-to-emails/util"
	"testing"
)

func TestFindFileName(t *testing.T) {
	ret, err := util.FindFileName("张三", "testdata")
	assert.Nil(t, err)
	assert.Equal(t, "张三的文件", ret)

	ret, err = util.FindFileName("李四", "testdata")
	assert.Equal(t, "", ret)
	assert.NotNil(t, err)
}

func TestFindFileName2(t *testing.T) {
	ret, err := util.FindFileName1("张三", "testdata")
	assert.Nil(t, err)
	assert.Equal(t, "testdata/张三的文件", ret)

	ret, err = util.FindFileName1("李四", "testdata")
	assert.Equal(t, "", ret)
	assert.NotNil(t, err)
}
