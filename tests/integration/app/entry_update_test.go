package app_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/biuren/hebe/internal/app/controller/entry"
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/app/view"
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/like"
	userDomain "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/pkg"
	"github.com/biuren/hebe/tests/framework"
	"github.com/biuren/hebe/tests/framework/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUpdateEntry_HTMX_Integration(t *testing.T) {
	// 1. Setup Environment
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	config := ctx.GetConfig().GetAppConfig()

	// Get Services
	userService := ctx.GetServices().GetUserService()
	diaryService := ctx.GetServices().GetDiaryService()
	entryService := ctx.GetServices().GetEntryService()

	// Manually create CommentService & LikeService
	commentRepo := repository.NewCommentRepository()
	commentService := comment.NewService(config, commentRepo, userService, entryService, diaryService)
	likeRepo := repository.NewLikeRepository()
	likeService := like.NewLikeService(likeRepo)

	// Create Helpers
	i18nManager := i18n.NewLocalizationManager(config)
	authHelper := userDomain.NewAuthHelper(config)
	renderService := view.NewViewService(config, i18nManager)

	// Initialize Controllers
	// Note: We need EntryController for UpdateEntry
	entryCtrl := entry.NewEntryController(config, entryService, diaryService, commentService, authHelper, renderService, userService, likeService)

	// 2. Prepare Data
	dataGen := ctx.GetDataGenerator()

	// Create User u1
	u1 := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u1")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u1@test.com")),
	)
	err := userService.RegisterUser(u1)
	require.NoError(t, err)

	// Create Diary for u1 (Public)
	d := dataGen.CreateDiary(u1.UserId)
	d.Title = "Public Diary"
	d.ViewSetting = 1
	d.Status = 1
	d.SocialSetting = 1 // Enabled
	diaryID, err := diaryService.AddDiary(d)
	require.NoError(t, err)

	// Create Entry
	e := dataGen.CreateEntry(u1.UserId, diaryID)
	e.Title = "Original Title"
	e.Body = "Original Body"
	e.Status = 1
	entryID, err := entryService.AddEntry(e)
	require.NoError(t, err)

	// Mock Auth User Map
	users := map[string]*types.AuthUser{
		fmt.Sprintf("%d", u1.UserId): {
			UserId:      u1.UserId,
			DisplayName: u1.DisplayName,
			UserGroup:   u1.UserGroup,
			Status:      u1.Status,
			Lang:        "en",
		},
	}

	// 3. Setup Router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	helpers.LoadTestTemplates(router)

	router.Use(helpers.MockAuthMiddleware(users))
	// Register routes
	// We need to register entry routes. UpdateEntry is at PUT /entry/:id
	router.PUT("/entry/:id", middleware.Validate[entry.UpdateEntryRequest](), entryCtrl.UpdateEntry)

	// 4. Perform Request
	// Simulate HTMX request to Update Entry
	updateReq := map[string]interface{}{
		"title":     "Updated Title",
		"info_type": 1, // Required field
		"body":      "Updated Body",
		"diaryID":   diaryID,
		"entryID":   entryID,
	}
	jsonBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/entry/%d?source=drafts", entryID), strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true") // Trigger HTMX logic
	req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", u1.UserId)) // Trigger MockAuthMiddleware

	// Set Auth Cookie
	authUser := users[fmt.Sprintf("%d", u1.UserId)]
	userDataBytes, _ := json.Marshal(authUser)
	cookieVal, _ := pkg.AESEncrypt(string(userDataBytes), config.Cookie.Secret)
	req.AddCookie(&http.Cookie{Name: "ud", Value: cookieVal})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 5. Assertions
	// Note: When source=drafts or source=box, UpdateEntry returns 204 No Content
	// The frontend uses HX-Trigger header to refresh the UI
	require.Equal(t, http.StatusNoContent, w.Code)

	// Verify HX-Trigger header contains entryUpdated event
	hxTrigger := w.Header().Get("HX-Trigger")
	require.Contains(t, hxTrigger, "refreshDrafts", "HX-Trigger should contain refreshDrafts for drafts source")

	// For drafts source, no HTML body is returned
	// The entryUpdated trigger is handled by frontend event listener
	body := w.Body.String()
	require.Empty(t, body, "Body should be empty for 204 response")
}
