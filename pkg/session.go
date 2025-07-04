package pkg

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// RedisConfig Redis配置结构
type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

// 创建redis client
func getRedisClient(config *RedisConfig) *redis.Client {
	// 链接redis
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       0,
	})

	return client
}

// 设置Session,默认保存4个小时
func SetSession(config *RedisConfig, name string, val string) error {
	// 获取client
	client := getRedisClient(config)

	err := client.Set(name, val, 4*time.Hour).Err()

	return err
}

// 获取session的值
func GetSession(config *RedisConfig, name string) (string, error) {
	//获取client
	client := getRedisClient(config)

	val, err := client.Get(name).Result()

	return val, err
}

// 删除一个session
func DelSession(config *RedisConfig, name string) (int64, error) {
	//获取client
	client := getRedisClient(config)

	num, err := client.Del(name).Result()

	return num, err
}
