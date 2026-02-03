package controller_test

import (
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/biuren/hebe/internal/app/controller/entry"
	"github.com/biuren/hebe/internal/app/templates"
	"github.com/biuren/hebe/internal/app/view"
	domainEntry "github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
	"github.com/biuren/hebe/internal/utils/i18n"
)

// MockRepo for Real Test
type MockRepoReal struct {
	ListAllEntriesByUserIDFunc func(userID int64, status int8, pageIndex, pageSize int) ([]domainEntry.Entry, int, error)
}

func (m *MockRepoReal) Insert(entry *domainEntry.Entry) (int64, error) { return 0, nil }
func (m *MockRepoReal) Update(entry *domainEntry.Entry) error          { return nil }
func (m *MockRepoReal) ListByUserID(userID int64, diaryID int64, status int8, pageIndex, pageSize int) ([]domainEntry.Entry, int, error) {
	return nil, 0, nil
}
func (m *MockRepoReal) GetCountByUserID(userID int64, diaryID int64, status int8) (int, error) {
	return 0, nil
}
func (m *MockRepoReal) ListAllByUserID(userID int64, status int8, pageIndex, pageSize int) ([]domainEntry.Entry, int, error) {
	if m.ListAllEntriesByUserIDFunc != nil {
		return m.ListAllEntriesByUserIDFunc(userID, status, pageIndex, pageSize)
	}
	return nil, 0, nil
}
func (m *MockRepoReal) GetByID(entryID int64) (*domainEntry.Entry, error) { return nil, nil }
func (m *MockRepoReal) GetByIDForUser(entryID int64, userID int64) (*domainEntry.Entry, error) {
	return nil, nil
}
func (m *MockRepoReal) UpdateByID(entry *domainEntry.Entry) error { return nil }
func (m *MockRepoReal) DeleteByID(entryID int64) error            { return nil }
func (m *MockRepoReal) HardDeleteByID(entryID int64) error        { return nil }
func (m *MockRepoReal) GetPreviousEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
	return 0, nil
}
func (m *MockRepoReal) GetNextEntryByDiaryID(currentEntryID int64, diaryID int64) (int64, error) {
	return 0, nil
}
func (m *MockRepoReal) GetStatsByIDs(entryIDs []int64) ([]domainEntry.Entry, error) { return nil, nil }
func (m *MockRepoReal) IncrementLikeNum(entryID int64) error                        { return nil }
func (m *MockRepoReal) DecrementLikeNum(entryID int64) error                        { return nil }

func TestShowDrafts_PanicReproduction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Setup Config
	absPath, _ := filepath.Abs("../../../../website/locales")
	viewPath, _ := filepath.Abs("../../../../website/views")

	cfg := &config.Config{
		App: config.AppConfig{
			PageSize:    10,
			LocalesDir:  absPath,
			DefaultLang: "zh-cn",
			Domain:      "http://localhost",
		},
	}

	// 2. Setup ViewService & Templates
	i18nManager := i18n.NewLocalizationManager(cfg)
	templates.SetGlobalManager(i18nManager)

	viewSvc := view.NewViewService(cfg, i18nManager)

	// 3. Setup AuthHelper
	authHelper := user.NewAuthHelper(cfg)

	// 4. Setup Service with Mock Repo
	mockRepo := &MockRepoReal{
		ListAllEntriesByUserIDFunc: func(userID int64, status int8, pageIndex, pageSize int) ([]domainEntry.Entry, int, error) {
			entries := []domainEntry.Entry{
				{
					EntryID:       1,
					DiaryID:       1,
					Title:         "Test Draft",
					Status:        2,
					UpdateUTCTime: 1678888888000,
				},
			}
			return entries, 1, nil
		},
	}
	entrySvc := domainEntry.NewService(cfg, nil, nil, nil, mockRepo, i18nManager)

	// 5. Setup Controller
	// We only need entryService, auth, and viewService for ShowDrafts
	ctrl := entry.NewEntryController(cfg, entrySvc, nil, nil, authHelper, viewSvc, nil, nil)

	// 6. Setup Gin Engine
	router := gin.New()

	// Manually load templates to handle recursive directories
	tmpl := template.New("")
	tmpl.Funcs(templates.RegisterAllFuncs())

	err := filepath.WalkDir(viewPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Parse the file content into the template set
			_, err = tmpl.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}
	router.SetHTMLTemplate(tmpl)

	// Setup Route
	router.Use(func(c *gin.Context) {
		// Simulate missing auth data to reproduce panic
		// Do not set IsLogin or AuthUser
		c.Next()
	})

	router.GET("/drafts", ctrl.ShowDrafts)

	// 7. Run Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/drafts", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic: %v", r)
		}
	}()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Logf("Response Body: %s", w.Body.String())
		t.Fatalf("Expected 302 Found (Redirect), got %d", w.Code)
	}
	
	location := w.Header().Get("Location")
	if location != "/login" {
		t.Fatalf("Expected redirect to /login, got %s", location)
	}
}
