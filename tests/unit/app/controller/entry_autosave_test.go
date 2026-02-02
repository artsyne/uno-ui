package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/biuren/hebe/internal/app/controller/entry"
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/app/templates"
	"github.com/biuren/hebe/internal/app/view"
	contracts "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/diary"
	domainEntry "github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/errors"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockUserRepo for testing
type MockUserRepo struct {
	user.RepositoryInterface
}

func (m *MockUserRepo) GetUserByUserId(userId int64) (*user.User, error) {
	return &user.User{UserId: userId, Handle: "testuser", DisplayName: "Test User"}, nil
}

// MockDiaryRepo for testing
type MockDiaryRepo struct {
	diary.RepositoryInterface
}

func (m *MockDiaryRepo) GetByDiaryID(diaryID int64) (*diary.Diary, error) {
	if diaryID == 999 {
		return nil, errors.DiaryNotFound
	}
	if diaryID == 888 {
		// Simulate BUG: return nil diary and nil error
		return nil, nil
	}
	return &diary.Diary{DiaryID: diaryID, AddUserID: 1, Title: "Test Diary", Status: 1}, nil
}
func (m *MockDiaryRepo) IsDiaryOwner(diaryID, userID int64) (bool, error) {
	return true, nil
}
func (m *MockDiaryRepo) CheckDiaryLimit(userID int64) error {
	return nil
}

// MockEntryRepo for testing
type MockEntryRepo struct {
	domainEntry.RepositoryInterface
}

func (m *MockEntryRepo) Insert(e *domainEntry.Entry) (int64, error) {
	return 100, nil
}
func (m *MockEntryRepo) Update(e *domainEntry.Entry) error {
	return nil
}
func (m *MockEntryRepo) GetByIDForUser(entryID int64, userID int64) (*domainEntry.Entry, error) {
	return &domainEntry.Entry{
		EntryID:   entryID,
		DiaryID:   1,
		AddUserID: userID,
		Status:    domainEntry.StatusDraft,
	}, nil
}

func TestAutoSaveEntry_PanicReproduction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Config
	absPath, _ := filepath.Abs("../../../../website/locales")
	cfg := &config.Config{
		App: config.AppConfig{
			PageSize:       10,
			LocalesDir:     absPath,
			DefaultLang:    "zh-cn",
			Domain:         "http://localhost",
			MaxEntryLength: 10000,
			SummaryMax:     200,
		},
	}

	// Setup Services
	i18nManager := i18n.NewLocalizationManager(cfg)
	templates.SetGlobalManager(i18nManager)
	viewSvc := view.NewViewService(cfg, i18nManager)
	authHelper := user.NewAuthHelper(cfg)

	mockUserRepo := &MockUserRepo{}
	userService := user.NewService(cfg, mockUserRepo, nil, i18nManager, nil)

	mockDiaryRepo := &MockDiaryRepo{}
	mockEntryRepo := &MockEntryRepo{}
	entrySvc := domainEntry.NewService(cfg, userService, nil, mockDiaryRepo, mockEntryRepo, i18nManager)

	// Setup Controller
	ctrl := entry.NewEntryController(cfg, entrySvc, nil, nil, authHelper, viewSvc, userService, nil)

	// Setup Router
	router := gin.New()
	router.Use(middleware.PanicRecovery())

	// Mock Auth Middleware
	router.Use(func(c *gin.Context) {
		c.Set("AuthUser", &contracts.AuthUser{
			UserId:      1,
			Handle:      "testuser",
			DisplayName: "Test User",
		})
		c.Set("IsLogin", true)
		c.Next()
	})

	// Register Route
	router.POST("/entry/auto-save", middleware.Validate[entry.AutoSaveEntryRequest](), ctrl.AutoSaveEntry)

	// Test Case: Create New Draft
	t.Run("Create New Draft", func(t *testing.T) {
		reqBody := entry.AutoSaveEntryRequest{
			DiaryID: 1,
			Body:    "Test Draft Content",
			Title:   "Test Draft",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/entry/auto-save", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"success":true`)
	})

	// Test Case: Diary Not Found (Simulate 9999 error cause)
	t.Run("Create New Draft - Diary Not Found", func(t *testing.T) {
		reqBody := entry.AutoSaveEntryRequest{
			DiaryID: 999, // Non-existent/Inactive diary
			Body:    "Test Draft Content",
			Title:   "Test Draft",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/entry/auto-save", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		t.Logf("Response Code: %d, Body: %s", w.Code, w.Body.String())

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1303), response["code"]) // Should be 1303 (DiaryNotFound)
	})

	t.Run("Create New Draft - Repo Returns Nil Nil (Panic Check)", func(t *testing.T) {
		reqBody := entry.AutoSaveEntryRequest{
			DiaryID: 888, // Triggers nil, nil in Mock
			Body:    "Panic Draft Content",
			Title:   "Panic Draft",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/entry/auto-save", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		t.Logf("Response Code: %d, Body: %s", w.Code, w.Body.String())

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// If panic caught, code should be 9999.
		// BUT now we fixed it, so it should be 1303 (DiaryNotFound)
		assert.Equal(t, float64(1303), response["code"])
	})
}
