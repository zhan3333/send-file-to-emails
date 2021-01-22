package cache_test

import (
	"github.com/stretchr/testify/assert"
	"send-fiule-to-emails/cache"
	"testing"
)

func TestInit(t *testing.T) {
	cache.Init("testdata/cache.json")
}

func TestSet(t *testing.T) {
	cache.Init("testdata/cache.json")
	cache.Set("foo", "bar")
}

func TestGet(t *testing.T) {
	cache.Init("testdata/cache.json")
	cache.Set("foo", "bar")
	val := cache.Get("foo")
	assert.Equal(t, val, "bar")
	noExistVal := cache.Get("test")
	assert.Nil(t, noExistVal)
}
