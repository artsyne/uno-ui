package controller_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/mock"

	"github.com/biuren/hebe/internal/app/controller/entry"
	"github.com/biuren/hebe/internal/app/view"
	contracts "github.com/biuren/hebe/internal/contracts"
	domainEntry "github.com/biuren/hebe/internal/domain/entry"
	"github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/utils/config"
)

// --- Mocks ---

type MockRepoRepro struct {
	mock.Mock
}

func (m *MockRepoRepro) Insert(entry *domainEntry.Entry) (int64, error) {
	args := m.Called(entry)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockRepoRepro) Update(entry *domainEntry.Entry) error {
	args := m.Called(entry)
	return args.Error(0)
}
func (m *MockRepoRepro) GetByID(id int64) (*domainEntry.Entry, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainEntry.Entry), args.Error(1)
}
func (m *MockRepoRepro) GetByIDRaw(id int64) (*domainEntry.Entry, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainEntry.Entry), args.Error(1)
}
func (m *MockRepoRepro) DeleteByID(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepoRepro) HardDeleteByID(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepoRepro) UpdateByID(entry *domainEntry.Entry) error {
	args := m.Called(entry)
	return args.Error(0)
}
func (m *MockRepoRepro) GetByIDForUser(entryID, userID int64) (*domainEntry.Entry, error) {
	args := m.Called(entryID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainEntry.Entry), args.Error(1)
}
func (m *MockRepoRepro) ListByUserID(userID, diaryID int64, status int8, page, size int) ([]domainEntry.Entry, int, error) {
	args := m.Called(userID, diaryID, status, page, size)
	return args.Get(0).([]domainEntry.Entry), args.Int(1), args.Error(2)
}
func (m *MockRepoRepro) ListAllByUserID(userID int64, status int8, page, size int) ([]domainEntry.Entry, int, error) {
	args := m.Called(userID, status, page, size)
	return args.Get(0).([]domainEntry.Entry), args.Int(1), args.Error(2)
}
func (m *MockRepoRepro) GetCountByUserID(userID, diaryID int64, status int8) (int, error) {
	args := m.Called(userID, diaryID, status)
	return args.Int(0), args.Error(1)
}
func (m *MockRepoRepro) ListByDiaryID(diaryID int64, status int8, pageIndex, pageSize int) ([]domainEntry.Entry, int, error) {
	args := m.Called(diaryID, status, pageIndex, pageSize)
	return args.Get(0).([]domainEntry.Entry), args.Int(1), args.Error(2)
}
func (m *MockRepoRepro) GetCountByDiaryID(diaryID int64, status int8) (int, error) {
	args := m.Called(diaryID, status)
	return args.Int(0), args.Error(1)
}
func (m *MockRepoRepro) IncrementViewNum(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepoRepro) IncrementLikeNum(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepoRepro) DecrementLikeNum(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockRepoRepro) GetPreviousEntryByDiaryID(currentEntryID, diaryID int64) (int64, error) {
	args := m.Called(currentEntryID, diaryID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockRepoRepro) GetNextEntryByDiaryID(currentEntryID, diaryID int64) (int64, error) {
	args := m.Called(currentEntryID, diaryID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockRepoRepro) GetStatsByIDs(entryIDs []int64) ([]domainEntry.Entry, error) {
	args := m.Called(entryIDs)
	return args.Get(0).([]domainEntry.Entry), args.Error(1)
}

type MockHTMLRenderRepro struct {
	Templates map[string]*template.Template
}

func (r *MockHTMLRenderRepro) Instance(name string, data interface{}) render.Render {
	return &MockRenderRepro{
		Template: r.Templates[name],
		Data:     data,
		Name:     name,
	}
}

type MockRenderRepro struct {
	Template *template.Template
	Data     interface{}
	Name     string
}

func (r *MockRenderRepro) Render(w http.ResponseWriter) error {
	if r.Template == nil {
		panic("Template not found: " + r.Name)
	}
	// Use ExecuteTemplate to render the specific defined template
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

func (r *MockRenderRepro) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
	}
}

// MockI18nManagerRepro mocks i18n manager
type MockI18nManagerRepro struct{}

func (m *MockI18nManagerRepro) GetMessage(key, lang string) string {
	return key
}
func (m *MockI18nManagerRepro) GetMessageWithParams(key, locale string, params map[string]string) string {
	return key
}
func (m *MockI18nManagerRepro) GetDefaultLanguage() string {
	return "zh-cn"
}
func (m *MockI18nManagerRepro) GetSupportedLanguages() []config.SupportLang {
	return []config.SupportLang{
		{Code: "en", Label: "English", LocalText: "English"},
		{Code: "zh-cn", Label: "简体中文", LocalText: "简体中文"},
	}
}
func (m *MockI18nManagerRepro) GetSupportedLanguageCodes() []string {
	return []string{"en", "zh-cn"}
}
func (m *MockI18nManagerRepro) IsLanguageSupported(locale string) bool {
	return true
}
func (m *MockI18nManagerRepro) GetLanguageInfo(locale string) (config.SupportLang, bool) {
	return config.SupportLang{Code: "en", Label: "English", LocalText: "English"}, true
}
func (m *MockI18nManagerRepro) GetMessageCount() int {
	return 0
}
func (m *MockI18nManagerRepro) GetMessagesByLanguage(locale string) map[string]string {
	return map[string]string{}
}
func (m *MockI18nManagerRepro) GetLocaleFromContext(c *gin.Context) string {
	return "en"
}
func (m *MockI18nManagerRepro) ReloadMessages() error {
	return nil
}

// --- Test ---

func TestShowDrafts_Reproduction_WithRealTemplates(t *testing.T) {
	// 1. Setup Config
	cfg := config.GetConfig()
	if cfg == nil {
		cfg = &config.Config{}
	}
	// Force PageSize to 10 for this test
	cfg.App.PageSize = 10

	// 2. Setup Context
	gin.SetMode(gin.TestMode)

	// 3. Setup Mocks
	mockRepo := new(MockRepoRepro)
	entrySvc := domainEntry.NewService(cfg, nil, nil, nil, mockRepo, nil)

	// Mock ListAllByUserID
	entries := []domainEntry.Entry{
		{EntryID: 100, DiaryID: 1, Title: "Draft 1", UpdateUTCTime: 1672531200000, Status: 2},
		{EntryID: 101, DiaryID: 1, Title: "Draft 2", UpdateUTCTime: 1672531200000, Status: 2},
	}
	mockRepo.On("ListAllByUserID", int64(123), int8(2), 1, 10).Return(entries, 20, nil)

	// 4. Load Real Templates
	tmpl := template.New("")
	// We MUST register functions BEFORE parsing
	tmpl.Funcs(template.FuncMap{
		"getMessage": func(key, lang string) string { return key },
		"formatLocalTime": func(ts int64, tz string) string {
			return time.Unix(ts/1000, 0).UTC().Format("2006-01-02 15:04")
		},
		"getTimeZone":    func(data interface{}) string { return "UTC" },
		"getUserID":      func(data interface{}) int64 { return 123 },
		"getDisplayName": func(data interface{}) string { return "Test User" },
		"getHandle":      func(data interface{}) string { return "testuser" },
	})

	// Define dummy partials to avoid missing template errors
	dummyPartials := `
{{define "header"}}<!-- header -->{{end}}
{{define "navbar"}}<!-- navbar -->{{end}}
{{define "footer"}}<!-- footer -->{{end}}
`
	_, err := tmpl.Parse(dummyPartials)
	if err != nil {
		t.Fatalf("Failed to parse dummy partials: %v", err)
	}

	// Read Real Template Files
	projectRoot := "/Users/yasong.cy/go/src/idiary.ai/hebe"
	files := []string{
		"website/views/diary/drafts.html",
		"website/views/diary/drafts_list.html",
		"website/views/pagination.html",
	}

	for _, f := range files {
		content, err := os.ReadFile(filepath.Join(projectRoot, f))
		if err != nil {
			t.Fatalf("Failed to read template file %s: %v", f, err)
		}
		// We parse them into the same template instance
		_, err = tmpl.Parse(string(content))
		if err != nil {
			t.Fatalf("Failed to parse template %s: %v", f, err)
		}
	}

	// 5. Setup Router and Services
	router := gin.New()
	router.HTMLRender = &MockHTMLRenderRepro{
		Templates: map[string]*template.Template{
			"diary/drafts": tmpl,
		},
	}

	// Real ViewService with Mock I18n
	viewSvc := view.NewViewService(cfg, &MockI18nManagerRepro{})

	// Real AuthHelper
	authHelper := user.NewAuthHelper(cfg)

	ctrl := entry.NewEntryController(cfg, entrySvc, nil, nil, authHelper, viewSvc, nil, nil)

	// Register route with middleware simulation
	router.GET("/drafts", func(c *gin.Context) {
		// Simulate SharedData middleware
		testUser := &contracts.AuthUser{
			UserId:      123,
			DisplayName: "Test User",
			Lang:        "en",
			TimeZone:    "UTC",
		}
		c.Set("IsLogin", true)
		c.Set("AuthUser", testUser)
		c.Set("Lang", "en")

		ctrl.ShowDrafts(c)
	})

	req, _ := http.NewRequest("GET", "/drafts", nil)
	w := httptest.NewRecorder()

	// Perform Request
	router.ServeHTTP(w, req)

	// Check Result
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify Draft 1 is present
	body := w.Body.String()
	if !contains(body, "Draft 1") {
		t.Errorf("Response body should contain 'Draft 1'. Got:\n%s", body)
	}
	if !contains(body, "2023-01-01 00:00") { // 1672531200000 -> 2023-01-01
		t.Errorf("Response body should contain formatted date")
	}

	// Verify Pagination is present
	// The URL generated is /drafts?status=2&page=2
	if !contains(body, "page=2") {
		t.Errorf("Response body should contain pagination link to page 2. Body snippet: %s", body[:min(len(body), 500)])
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr ||
		(len(s) > len(substr) && contains(s[1:], substr))
}
