package cache

import "time"

// Session helper constants
const (
	SessionPrefix     = "session:"
	DefaultSessionTTL = 24 * time.Hour
)

// SetSession 存储会话数据 (自动添加 session: 前缀)
// sessionID: 会话唯一标识，内部会自动添加 "session:" 前缀
// data: 会话数据，可以是任意结构体，将自动序列化为 JSON 存储
// expiration: 过期时间。单位：time.Duration (如 time.Second, time.Hour)。如果传 0，则使用默认 24 小时
func (m *Manager) SetSession(sessionID string, data interface{}, expiration time.Duration) error {
	if expiration == 0 {
		expiration = DefaultSessionTTL
	}
	return m.SetJSON(SessionPrefix+sessionID, data, expiration)
}

// GetSession 获取会话数据
// sessionID: 会话唯一标识
// dest: 目标对象指针，用于接收反序列化后的会话数据
func (m *Manager) GetSession(sessionID string, dest interface{}) error {
	return m.GetJSON(SessionPrefix+sessionID, dest)
}

// DelSession 删除会话
// sessionID: 会话唯一标识
func (m *Manager) DelSession(sessionID string) error {
	return m.Del(SessionPrefix + sessionID)
}
