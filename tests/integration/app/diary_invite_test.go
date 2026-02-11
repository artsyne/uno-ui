package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/biuren/hebe/internal/app/controller/diary"
	"github.com/biuren/hebe/internal/app/view"
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/like"
	userDomain "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/tests/framework"
	"github.com/biuren/hebe/tests/framework/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiaryInviteAuthor_InvalidUser_NoPanic(t *testing.T) {
	// 1. Setup Environment
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	config := ctx.GetConfig().GetAppConfig()

	// Get Services from Registry
	userService := ctx.GetServices().GetUserService()
	diaryService := ctx.GetServices().GetDiaryService()
	entryService := ctx.GetServices().GetEntryService()

	// Manually create LikeService (not in registry yet)
	likeRepo := repository.NewLikeRepository()
	likeService := like.NewLikeService(likeRepo)

	// Create Helpers
	i18nManager := i18n.NewLocalizationManager(config)
	authHelper := userDomain.NewAuthHelper(config)
	renderService := view.NewViewService(config, i18nManager)

	// Initialize Controller
	diaryCtrl := diary.NewDiaryController(config, diaryService, entryService, userService, likeService, authHelper, renderService)

	// 2. Prepare Data
	dataGen := ctx.GetDataGenerator()

	// Create Owner
	owner := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u1")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u1@test.com")),
	)
	err := userService.RegisterUser(owner)
	require.NoError(t, err)

	// Create Diary
	diaryObj := dataGen.CreateDiary(owner.UserId)
	diaryID, err := diaryService.AddDiary(diaryObj)
	require.NoError(t, err)

	// Mock Auth Middleware
	userMap := map[string]*types.AuthUser{
		fmt.Sprintf("%d", owner.UserId): {
			UserId:      owner.UserId,
			DisplayName: owner.DisplayName,
			UserGroup:   owner.UserGroup,
			Status:      owner.Status,
			Lang:        "en",
		},
	}

	// Setup Router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	helpers.LoadTestTemplates(router)
	router.Use(gin.Recovery())
	router.Use(helpers.MockAuthMiddleware(userMap))
	
	// Register Routes
	diaryCtrl.RegisterRoutes(router)

	// 3. Test Case: Invite non-existent user
	t.Run("InviteNonExistentUser_ShouldNotPanic", func(t *testing.T) {
		// Prepare Request
		// Form data: invitee_identity=nonexistent
		formData := url.Values{}
		formData.Set("invitee_identity", "nonexistentuser12345")

		req, _ := http.NewRequest("POST", fmt.Sprintf("/diary/%d/authors", diaryID), bytes.NewBufferString(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", owner.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert Response
		// Should return 404 (User Not Found) instead of 500 (Panic)
		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, false, resp["success"])
		assert.Equal(t, "user_not_found", resp["message"])
	})
}

