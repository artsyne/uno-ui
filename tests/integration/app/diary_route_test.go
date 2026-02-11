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
	"github.com/stretchr/testify/require"
)


func TestDiaryRoutes_Integration(t *testing.T) {
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

	// Create Invitee
	invitee := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u2")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u2@test.com")),
	)
	err = userService.RegisterUser(invitee)
	require.NoError(t, err)

	// Create Stranger
	stranger := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u3")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u3@test.com")),
	)
	err = userService.RegisterUser(stranger)
	require.NoError(t, err)

	// Create Diary (Private)
	diaryObj := dataGen.CreateDiary(owner.UserId)
	diaryID, err := diaryService.AddDiary(diaryObj)
	require.NoError(t, err)

	// Convert to AuthUser for Mock Middleware
	users := map[string]*types.AuthUser{
		fmt.Sprintf("%d", owner.UserId): {
			UserId:      owner.UserId,
			DisplayName: owner.DisplayName,
			UserGroup:   owner.UserGroup,
			Status:      owner.Status,
			Lang:        "en",
		},
		fmt.Sprintf("%d", invitee.UserId): {
			UserId:      invitee.UserId,
			DisplayName: invitee.DisplayName,
			UserGroup:   invitee.UserGroup,
			Status:      invitee.Status,
			Lang:        "en",
		},
		fmt.Sprintf("%d", stranger.UserId): {
			UserId:      stranger.UserId,
			DisplayName: stranger.DisplayName,
			UserGroup:   stranger.UserGroup,
			Status:      stranger.Status,
			Lang:        "en",
		},
	}

	// 3. Setup Router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	helpers.LoadTestTemplates(router)

	router.Use(helpers.MockAuthMiddleware(users))
	diaryCtrl.RegisterRoutes(router)

	// 4. Test Cases

	t.Run("InviteAuthor_Success", func(t *testing.T) {
		// Owner invites Invitee
		data := url.Values{}
		data.Set("invitee_identity", invitee.Email)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/diary/%d/authors", diaryID), bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", owner.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, true, resp["success"], "Expected success: true, got: %v", resp)
	})

	t.Run("InviteAuthor_Forbidden_NotOwner", func(t *testing.T) {
		// Stranger tries to invite Invitee
		data := url.Values{}
		data.Set("invitee_identity", invitee.Email)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/diary/%d/authors", diaryID), bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", stranger.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, false, resp["success"], "Expected success: false")
	})

	t.Run("AcceptInvitation_Success", func(t *testing.T) {
		// Invitee accepts
		data := url.Values{}
		data.Set("diary_id", fmt.Sprintf("%d", diaryID))

		req, _ := http.NewRequest("POST", "/invitations/diary/accept", bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", invitee.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, true, resp["success"])
	})

	t.Run("UpdateAuthor_Success", func(t *testing.T) {
		// Owner updates Invitee to Admin (Role 2)
		data := url.Values{}
		data.Set("role", "2")

		req, _ := http.NewRequest("PUT", fmt.Sprintf("/diary/%d/authors/%d", diaryID, invitee.UserId), bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", owner.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, true, resp["success"])
	})

	t.Run("RemoveAuthor_Success", func(t *testing.T) {
		// Owner removes Invitee
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/diary/%d/authors/%d", diaryID, invitee.UserId), nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", owner.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, true, resp["success"])
	})

	t.Run("QuitDiary_Success", func(t *testing.T) {
		// Re-invite and accept first via Service
		err := diaryService.InviteAuthor(diaryID, owner.UserId, invitee.Email, "en")
		require.NoError(t, err)
		err = diaryService.AcceptInvitation(diaryID, invitee.UserId)
		require.NoError(t, err)

		// Invitee quits
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/diary/%d/authors/self", diaryID), nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", invitee.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, true, resp["success"])
	})

	t.Run("GetAuthors_Owner_Success", func(t *testing.T) {
		// Owner accesses authors list
		req, _ := http.NewRequest("GET", fmt.Sprintf("/diary/%d/authors", diaryID), nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", owner.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		// Verify content type is HTML
		require.Contains(t, w.Header().Get("Content-Type"), "text/html")
	})

	t.Run("ShowInviteForm_Success", func(t *testing.T) {
		// Owner accesses invite form
		req, _ := http.NewRequest("GET", fmt.Sprintf("/diary/%d/authors/invite", diaryID), nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", owner.UserId))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		// Verify content type is HTML
		require.Contains(t, w.Header().Get("Content-Type"), "text/html")
		// Verify Lang is rendered (by checking placeholder or label if possible, or just 200 OK)
		// Assuming template renders successfully means Lang was present
	})
}
