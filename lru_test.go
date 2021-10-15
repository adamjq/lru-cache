package lru

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLRUCache_New(t *testing.T) {
	assert := require.New(t)

	_, err := New(0)
	assert.Error(err)

	_, err = New(-5)
	assert.Error(err)

	_, err = New(1)
	assert.NoError(err)
}

func TestLRUCache_simple(t *testing.T) {
	assert := require.New(t)

	cache, err := New(2)
	assert.NoError(err)

	v := cache.Get("key")
	assert.Nil(v)
	assert.Equal(0, len(cache.cache))

	cache.Put("key", "value")
	v = cache.Get("key")
	assert.Equal(*v, "value")
	assert.Equal(1, len(cache.cache))
}

func TestLRUCache_updateSameKey(t *testing.T) {
	assert := require.New(t)

	cache, err := New(2)
	assert.NoError(err)

	v := cache.Get("key")
	assert.Nil(v)
	assert.Equal(0, len(cache.cache))

	cache.Put("key", "value")
	v = cache.Get("key")
	assert.Equal(*v, "value")
	assert.Equal(1, len(cache.cache))

	cache.Put("key", "new-value")
	v = cache.Get("key")
	assert.Equal(*v, "new-value")
	assert.Equal(1, len(cache.cache))
}

func TestLRUCache_eviction(t *testing.T) {
	assert := require.New(t)

	cache, err := New(2)
	assert.NoError(err)

	cache.Put("key-1", "value")
	assert.Equal(1, len(cache.cache))
	cache.Put("key-2", "value-2")
	assert.Equal(2, len(cache.cache))
	cache.Put("key-3", "value-3")
	assert.Equal(2, len(cache.cache))

	v := cache.Get("key-1")
	assert.Nil(v)

	v = cache.Get("key-2")
	assert.Equal(*v, "value-2")

	v = cache.Get("key-3")
	assert.Equal(*v, "value-3")
}

func TestLRUCache_sequentialPut(t *testing.T) {
	assert := require.New(t)

	cacheSize := 100
	cache, err := New(cacheSize)
	assert.NoError(err)

	for i := 0; i < cacheSize; i++ {
		cache.Put(fmt.Sprintf("key-%v", i), fmt.Sprintf("value-%v", i))
	}
	assert.Equal(cacheSize, len(cache.cache))
}

func TestLRUCache_sequentialPutExceedsCacheSize(t *testing.T) {
	assert := require.New(t)

	cacheSize := 100
	cache, err := New(cacheSize)
	assert.NoError(err)

	for i := 0; i < cacheSize*2; i++ {
		cache.Put(fmt.Sprintf("key-%v", i), fmt.Sprintf("value-%v", i))
	}
	assert.Equal(cacheSize, len(cache.cache))
}

func TestLRUCache_parallelPut(t *testing.T) {
	assert := require.New(t)

	cacheSize := 500
	cache, err := New(cacheSize)
	assert.NoError(err)

	for i := 0; i < cacheSize; i++ {
		go cache.Put(fmt.Sprintf("key-%v", i), fmt.Sprintf("value-%v", i))
	}
	time.Sleep(time.Second)
	assert.Equal(cacheSize, len(cache.cache))
}