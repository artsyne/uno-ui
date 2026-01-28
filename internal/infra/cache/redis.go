package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/biuren/hebe/internal/utils/config"
	"github.com/go-redis/redis"
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Manager 封装Redis客户端，提供便捷的缓存管理方法
type Manager struct {
	client *redis.Client
}

// NewManager 创建一个新的缓存管理器
// client: 已初始化的 Redis 客户端
func NewManager(client *redis.Client) *Manager {
	return &Manager{
		client: client,
	}
}

// Set 存储基本类型值
// key: 缓存键
// value: 缓存值（支持 string, int, bool, []byte 等基本类型）
// expiration: 过期时间，0 表示永不过期。单位：time.Duration (如 time.Second, time.Hour)
func (m *Manager) Set(key string, value interface{}, expiration time.Duration) error {
	return m.client.Set(key, value, expiration).Err()
}

// Get 获取字符串值
// key: 缓存键
// 返回: 字符串值，如果键不存在返回 redis.Nil 错误
func (m *Manager) Get(key string) (string, error) {
	return m.client.Get(key).Result()
}

// Del 删除键
// key: 缓存键
// 返回: 错误信息
func (m *Manager) Del(key string) error {
	return m.client.Del(key).Err()
}

// Exists 检查键是否存在
// key: 缓存键
// 返回: true 存在, false 不存在
func (m *Manager) Exists(key string) (bool, error) {
	count, err := m.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Expire 设置过期时间
// key: 缓存键
// expiration: 过期时间。单位：time.Duration (如 time.Second, time.Hour)
func (m *Manager) Expire(key string, expiration time.Duration) error {
	return m.client.Expire(key, expiration).Err()
}

// SetJSON 将对象序列化为JSON存储
// key: 缓存键
// value: 任意可序列化为 JSON 的结构体或 map
// expiration: 过期时间，0 表示永不过期。单位：time.Duration (如 time.Second, time.Hour)
func (m *Manager) SetJSON(key string, value interface{}, expiration time.Duration) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(key, jsonBytes, expiration).Err()
}

// GetJSON 获取JSON数据并反序列化到对象中
// key: 缓存键
// dest: 目标对象指针，用于接收反序列化后的数据
func (m *Manager) GetJSON(key string, dest interface{}) error {
	val, err := m.client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
