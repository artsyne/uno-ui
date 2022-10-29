package helpers

import (
	"time"
	"strconv"

	"chinadiary.com/website/system"
	"github.com/go-redis/redis"
)

//创建redis client
func getRedisClient() (*redis.Client) {
	//获取配置信息
	config := system.GetConfig()

	// 链接redis
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.RedisServer + ":" + strconv.Itoa(config.Redis.RedisPort),
		Password: config.Redis.RedisPwd,
		DB:       0,
	})

	return client
}

// 设置Session,默认保存4个小时
func SetSession(name string, val string) error {

	// 获取client
	client := getRedisClient()

	err := client.Set(name, val, 4*time.Hour).Err()

	return err
}

//获取session的值
func GetSession(name string) (string, error) {

	//获取client
	client := getRedisClient()

	val, err := client.Get(name).Result()

	return val, err
}

//删除一个session
func DelSession(name string) (int64,error){
	//获取client
	client := getRedisClient()

	num, err := client.Del(name).Result()

	return num, err
}