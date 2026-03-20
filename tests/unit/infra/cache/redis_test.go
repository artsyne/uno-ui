package cache_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/biuren/hebe/internal/infra/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*cache.Manager, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	manager := cache.NewManager(client)
	return manager, s
}

func TestManager_Set_Get(t *testing.T) {
	manager, s := setupTestRedis(t)
	defer s.Close()

	key := "test_key"
	value := "test_value"

	err := manager.Set(key, value, time.Minute)
	assert.NoError(t, err)

	got, err := manager.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, got)

	// Test Expiration
	s.FastForward(2 * time.Minute)
	_, err = manager.Get(key)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestManager_Del(t *testing.T) {
	manager, s := setupTestRedis(t)
	defer s.Close()

	key := "test_key"
	err := manager.Set(key, "value", 0)
	assert.NoError(t, err)

	err = manager.Del(key)
	assert.NoError(t, err)

	exists, err := manager.Exists(key)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestManager_JSON(t *testing.T) {
	manager, s := setupTestRedis(t)
	defer s.Close()

	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	key := "test_json"
	val := TestStruct{Name: "Alice", Age: 30}

	err := manager.SetJSON(key, val, 0)
	assert.NoError(t, err)

	var got TestStruct
	err = manager.GetJSON(key, &got)
	assert.NoError(t, err)
	assert.Equal(t, val, got)
}
