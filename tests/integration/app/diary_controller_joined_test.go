package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/app/controller/diary"
	"github.com/biuren/hebe/internal/app/view"
	types "github.com/biuren/hebe/internal/contracts"
	diaryDomain "github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/like"
	userDomain "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/tests/framework"
	"github.com/biuren/hebe/tests/framework/helpers"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockHTMLRender for testing
type MockHTMLRender struct {
	TemplateName string
	RenderFunc   func(name string, data interface{}) string
}

func (r *MockHTMLRender) Instance(name string, data interface{}) render.Render {
	return &MockRender{
		HTML: r.RenderFunc(name, data),
	}
}

type MockRender struct {
	HTML string
}

func (r *MockRender) Render(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(r.HTML))
	return err
}

func (r *MockRender) WriteContentType(w http.ResponseWriter) {}

func TestShowUserProfile_JoinedDiaries_Optimization(t *testing.T) {
	// 1. Setup Environment
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	config := ctx.GetConfig().GetAppConfig()
	// Explicitly set config for test stability
	config.App.UserJoinedDiariesLimit = 10
	config.App.DefaultLang = "en"

	// Get Services
	userService := ctx.GetServices().GetUserService()
	diaryService := ctx.GetServices().GetDiaryService()
	entryService := ctx.GetServices().GetEntryService()
	
	likeRepo := repository.NewLikeRepository()
	likeService := like.NewLikeService(likeRepo)

	// Create Helpers
	i18nManager := i18n.NewLocalizationManager(config)
	authHelper := userDomain.NewAuthHelper(config)
	renderService := view.NewViewService(config, i18nManager)

	// Initialize Controller
	diaryCtrl := diary.NewDiaryController(config, diaryService, entryService, userService, likeService, authHelper, renderService)

	// 2. Prepare Data
	// Create User A (The Profile Owner)
	userA := createTestUser(t, userService, ctx.GetDataGenerator().GenerateUniqueHandle("userA"), ctx.GetDataGenerator().GenerateUniqueEmail("usera"))
	
	// Create User B (The Other User)
	userB := createTestUser(t, userService, ctx.GetDataGenerator().GenerateUniqueHandle("userB"), ctx.GetDataGenerator().GenerateUniqueEmail("userb"))

	// Create a Public Diary for User B
	diaryB := &diaryDomain.Diary{
		Title:       "User B Public Diary",
		Remark:      "Public stuff",
		AddUserID:   userB.UserId,
		ViewSetting: diaryDomain.ViewSettingPublic,
		Status:      diaryDomain.StatusActive,
		AddUTCTime:  time.Now().Unix(),
	}
	diaryBID, err := diaryService.AddDiary(diaryB)
	require.NoError(t, err)

	// User A joins User B's diary
	// 1. User B invites User A
	err = diaryService.InviteAuthor(diaryBID, userB.UserId, userA.Email, "en")
	require.NoError(t, err)
	// 2. User A accepts
	err = diaryService.AcceptInvitation(diaryBID, userA.UserId)
	require.NoError(t, err)

	// Verify User A has joined diaries
	joined, _ := diaryService.ListJoinedDiariesByUserID(userA.UserId)
	require.Equal(t, 1, len(joined), "User A should have 1 joined diary")

	// 3. Setup Router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Mock templates with minimal content to verify data presence
	router.HTMLRender = &MockHTMLRender{
		TemplateName: "user/index",
		RenderFunc: func(name string, data interface{}) string {
			pageData, ok := data.(*types.PageData[map[string]interface{}])
			if !ok {
				return fmt.Sprintf("Error: Invalid Data Type: %T", data)
			}
			
			// Check if JoinedDiaries is present in data
			hasJoined := false
			if jd, ok := pageData.Data["JoinedDiaries"]; ok && jd != nil {
				diaries := jd.([]diaryDomain.Diary)
				if len(diaries) > 0 {
					hasJoined = true
				}
			}
			
			isOwner := pageData.Data["IsOwner"].(bool)
			
			return fmt.Sprintf("IsOwner: %v, HasJoined: %v", isOwner, hasJoined)
		},
	}

	// Setup Auth Middleware Mock
	users := map[string]*types.AuthUser{
		fmt.Sprintf("%d", userA.UserId): {
			UserId:      userA.UserId,
			DisplayName: userA.DisplayName,
			UserGroup:   userA.UserGroup,
			Status:      userA.Status,
		},
		fmt.Sprintf("%d", userB.UserId): {
			UserId:      userB.UserId,
			DisplayName: userB.DisplayName,
			UserGroup:   userB.UserGroup,
			Status:      userB.Status,
		},
	}
	router.Use(helpers.MockAuthMiddleware(users))

	// Register Route
	router.GET("/u/:uid", diaryCtrl.ShowUserProfile)

	// 4. Test Case 1: Owner (User A) visits their own profile
	t.Run("Owner sees Joined Diaries", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/u/%d", userA.UserId), nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", userA.UserId))
		
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "IsOwner: true")
		assert.Contains(t, w.Body.String(), "HasJoined: true")
	})

	// 5. Test Case 2: Non-Owner (User B) visits User A's profile
	t.Run("Non-Owner does NOT see Joined Diaries", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/u/%d", userA.UserId), nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", userB.UserId))
		
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "IsOwner: false")
		assert.Contains(t, w.Body.String(), "HasJoined: false")
	})
}

func createTestUser(t *testing.T, userService *userDomain.Service, handle, email string) *userDomain.User {
	u := &userDomain.User{
		Email:       email,
		Pwd:         "Password123!",
		Handle:      handle,
		DisplayName: handle,
		Status:      userDomain.StatusActive,
		UserGroup:   1,
		Lang:        "en",
		TimeZone:    "UTC",
	}
	err := userService.RegisterUser(u)
	require.NoError(t, err)
	return u
}
