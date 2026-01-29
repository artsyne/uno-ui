package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestManager_Session(t *testing.T) {
	manager, s := setupTestRedis(t)
	defer s.Close()

	sessionID := "sess_123"
	type SessionData struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
	}

	data := SessionData{UserID: 1, Username: "testuser"}

	// Test SetSession
	err := manager.SetSession(sessionID, data, time.Hour)
	assert.NoError(t, err)

	// Verify prefix in redis
	exists := s.Exists("session:" + sessionID)
	assert.True(t, exists)

	// Test GetSession
	var got SessionData
	err = manager.GetSession(sessionID, &got)
	assert.NoError(t, err)
	assert.Equal(t, data, got)

	// Test DelSession
	err = manager.DelSession(sessionID)
	assert.NoError(t, err)

	exists = s.Exists("session:" + sessionID)
	assert.False(t, exists)
}
