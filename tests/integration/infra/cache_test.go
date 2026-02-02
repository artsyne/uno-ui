package infra

import (
	"fmt"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/infra/cache"
	framework "github.com/biuren/hebe/tests/framework"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisManager(t *testing.T) {
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer func() {
		if cleanupErr := ctx.Cleanup(); cleanupErr != nil {
			t.Logf("Cleanup failed: %v", cleanupErr)
		}
	}()

	config := ctx.GetConfig().GetAppConfig()

	// Direct client creation for testing
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		t.Logf("Redis ping failed: %v", err)
		// t.Skip("Redis not available, skipping cache tests")
	}

	manager := cache.NewManager(client)

	t.Run("Set and Get", func(t *testing.T) {
		key := "test_key_" + ctx.GetDataGenerator().GenerateUniqueString("k")
		val := "test_value"
		err := manager.Set(key, val, time.Minute)
		require.NoError(t, err)

		got, err := manager.Get(key)
		require.NoError(t, err)
		assert.Equal(t, val, got)
	})

	t.Run("SetJSON and GetJSON", func(t *testing.T) {
		type TestStruct struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		key := "test_json_" + ctx.GetDataGenerator().GenerateUniqueString("k")
		val := TestStruct{Name: "Alice", Age: 30}
		
		err := manager.SetJSON(key, val, time.Minute)
		require.NoError(t, err)

		var got TestStruct
		err = manager.GetJSON(key, &got)
		require.NoError(t, err)
		assert.Equal(t, val, got)
	})

	t.Run("Session", func(t *testing.T) {
		sid := "sess_" + ctx.GetDataGenerator().GenerateUniqueString("s")
		data := map[string]string{"user": "bob"}
		
		err := manager.SetSession(sid, data, 0)
		require.NoError(t, err)
		
		var got map[string]string
		err = manager.GetSession(sid, &got)
		require.NoError(t, err)
		assert.Equal(t, data["user"], got["user"])
		
		err = manager.DelSession(sid)
		require.NoError(t, err)
		
		err = manager.GetSession(sid, &got)
		assert.Error(t, err)
	})
}
