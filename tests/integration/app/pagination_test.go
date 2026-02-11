package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/app/controller/diary"
	"github.com/biuren/hebe/internal/app/controller/entry"
	"github.com/biuren/hebe/internal/app/view"
	types "github.com/biuren/hebe/internal/contracts"
	"github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/like"
	userDomain "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/tests/framework"
	"github.com/biuren/hebe/tests/framework/helpers"
	testTypes "github.com/biuren/hebe/tests/framework/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestPagination_Integration(t *testing.T) {
	// 1. Setup Environment
	// Define a custom option to modify config before services are initialized
	withCustomConfig := func(tc *framework.TestContext) {
		if cfg, ok := tc.GetConfig().(*testTypes.TestConfig); ok {
			cfg.App.PageSize = 2
			cfg.App.UserJoinedDiariesLimit = 100
		}
	}

	// Pass custom config option BEFORE WithDatabase to ensure services use the modified config
	ctx := framework.NewTestContext(
		func(tc *framework.TestContext) { withCustomConfig(tc) },
		framework.WithDatabase(),
	)
	defer ctx.Cleanup()

	config := ctx.GetConfig().GetAppConfig()

	// Get Services
	userService := ctx.GetServices().GetUserService()
	diaryService := ctx.GetServices().GetDiaryService()
	entryService := ctx.GetServices().GetEntryService()
	
	// Manually create CommentService
	commentRepo := repository.NewCommentRepository()
	commentService := comment.NewService(config, commentRepo, userService, entryService, diaryService)

	// Manually create LikeService
	likeRepo := repository.NewLikeRepository()
	likeService := like.NewLikeService(likeRepo)

	// Create Helpers
	i18nManager := i18n.NewLocalizationManager(config)
	authHelper := userDomain.NewAuthHelper(config)
	renderService := view.NewViewService(config, i18nManager)

	// Initialize Controllers
	diaryCtrl := diary.NewDiaryController(config, diaryService, entryService, userService, likeService, authHelper, renderService)
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

	// Create User u2 (for invitations)
	u2 := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u2")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u2@test.com")),
	)
	err = userService.RegisterUser(u2)
	require.NoError(t, err)

	// Create User u3 (empty state)
	u3 := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u3")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u3@test.com")),
	)
	err = userService.RegisterUser(u3)
	require.NoError(t, err)

	// Create multiple Diaries for u1 (3 diaries, PageSize=2 -> 2 pages)
	var u1DiaryIDs []int64
	for i := 0; i < 3; i++ {
		d := dataGen.CreateDiary(u1.UserId)
		d.Title = fmt.Sprintf("U1 Diary %d", i)
		id, err := diaryService.AddDiary(d)
		require.NoError(t, err)
		u1DiaryIDs = append(u1DiaryIDs, id)
	}

	// Create Joined Diaries for u1 (u2 invites u1)
	// Create 3 diaries for u2 and invite u1
	for i := 0; i < 3; i++ {
		d := dataGen.CreateDiary(u2.UserId)
		d.Title = fmt.Sprintf("U2 Diary %d", i)
		id, err := diaryService.AddDiary(d)
		require.NoError(t, err)
		
		err = diaryService.InviteAuthor(id, u2.UserId, u1.Email, "en")
		require.NoError(t, err)
		err = diaryService.AcceptInvitation(id, u1.UserId)
		require.NoError(t, err)
		// Ensure distinct join times for stable sorting
		time.Sleep(1100 * time.Millisecond)
	}

	// Create Entries for u1's first diary (3 entries, PageSize=2 -> 2 pages)
	targetDiaryID := u1DiaryIDs[0]
	for i := 0; i < 3; i++ {
		e := dataGen.CreateEntry(u1.UserId, targetDiaryID)
		e.Title = fmt.Sprintf("Entry %d", i)
		_, err := entryService.AddEntry(e)
		require.NoError(t, err)
	}

	// Create Drafts for u1 (3 drafts)
	for i := 0; i < 3; i++ {
		e := dataGen.CreateEntry(u1.UserId, targetDiaryID)
		e.Title = fmt.Sprintf("Draft %d", i)
		e.Status = 2 // Draft
		_, err := entryService.AddEntry(e)
		require.NoError(t, err)
	}

	// Mock Auth User Map
	users := map[string]*types.AuthUser{
		fmt.Sprintf("%d", u1.UserId): {
			UserId:      u1.UserId,
			DisplayName: u1.DisplayName,
			UserGroup:   u1.UserGroup,
			Status:      u1.Status,
			Lang:        "en",
		},
		fmt.Sprintf("%d", u2.UserId): {
			UserId:      u2.UserId,
			DisplayName: u2.DisplayName,
			UserGroup:   u2.UserGroup,
			Status:      u2.Status,
			Lang:        "en",
		},
		fmt.Sprintf("%d", u3.UserId): {
			UserId:      u3.UserId,
			DisplayName: u3.DisplayName,
			UserGroup:   u3.UserGroup,
			Status:      u3.Status,
			Lang:        "en",
		},
	}

	// 3. Setup Router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	helpers.LoadTestTemplates(router) // Reuse helper from diary_route_test.go

	router.Use(helpers.MockAuthMiddleware(users))
	diaryCtrl.RegisterRoutes(router)
	entryCtrl.RegisterRoutes(router)

	// Helper to make request
	makeRequest := func(method, path string, isHTMX bool, hxTarget string, userID int64) *httptest.ResponseRecorder {
		req, _ := http.NewRequest(method, path, nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", userID))
		if isHTMX {
			req.Header.Set("HX-Request", "true")
			if hxTarget != "" {
				req.Header.Set("HX-Target", hxTarget)
			}
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}

	// 4. Test Cases

	// --- ShowUserProfile ---
	t.Run("ShowUserProfile_Full", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d", u1.UserId), false, "", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "<html", "Full page should contain html tag")
		require.Contains(t, body, "id=\"diary-container\"", "Should contain diary container")
	})

	t.Run("ShowUserProfile_Partial", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d", u1.UserId), true, "", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.NotContains(t, body, "<html", "Partial page should NOT contain html tag")
		require.Contains(t, body, "id=\"diary-container\"", "Should contain diary container")
	})

	// --- ShowJoinedDiaries ---
	t.Run("ShowJoinedDiaries_Full", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d/joined", u1.UserId), false, "", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "<html", "Full page should contain html tag")
		require.Contains(t, body, "id=\"joined-diaries-content\"", "Should contain joined diaries container")
	})

	t.Run("ShowJoinedDiaries_Partial_Page2", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d/joined/p/2", u1.UserId), true, "#joined-diaries-content", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.NotContains(t, body, "<html", "Partial page should NOT contain html tag")
		require.Contains(t, body, "id=\"joined-diaries-content\"", "Should contain joined diaries container")
		// Check if it contains expected content from page 2
		// Diaries are created: U2 Diary 0, U2 Diary 1, U2 Diary 2
		// Default order is usually by create time desc
		// Page 1: U2 Diary 2, U2 Diary 1
		// Page 2: U2 Diary 0
		require.Contains(t, body, "U2 Diary 0", "Should contain 'U2 Diary 0' on page 2")
		require.NotContains(t, body, "U2 Diary 2", "Should NOT contain 'U2 Diary 2' on page 2")
	})

	t.Run("ShowJoinedDiaries_Empty", func(t *testing.T) {
		// User u3 has no joined diaries (only created ones)
		w := makeRequest("GET", fmt.Sprintf("/u/%d/joined", u3.UserId), true, "#joined-diaries-content", u3.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "id=\"joined-diaries-content\"", "Should contain container")
		require.Contains(t, body, "empty-placeholder", "Should contain empty placeholder class")
		require.Contains(t, body, "mdi-inbox-outline", "Should contain empty icon")
	})

	t.Run("ShowJoinedDiaries_OutOfBounds", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d/joined/p/999", u1.UserId), true, "#joined-diaries-content", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "id=\"joined-diaries-content\"", "Should contain container")
		// Should show empty or just header?
		// Usually redirects to last page or shows empty list.
		// Based on logic, it might show empty list if Offset is large.
	})

	// --- ShowDrafts ---
	t.Run("ShowDrafts_Full", func(t *testing.T) {
		w := makeRequest("GET", "/drafts", false, "", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "<html", "Full page should contain html tag")
		require.Contains(t, body, "id=\"drafts-list\"", "Should contain drafts list container")
	})

	t.Run("ShowDrafts_Partial", func(t *testing.T) {
		w := makeRequest("GET", "/drafts", true, "drafts-list", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.NotContains(t, body, "<html", "Partial page should NOT contain html tag")
		require.Contains(t, body, "id=\"drafts-list\"", "Should contain drafts list container")
	})

	t.Run("ShowDrafts_Partial_Page2", func(t *testing.T) {
		w := makeRequest("GET", "/drafts?page=2", true, "drafts-list", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.NotContains(t, body, "<html", "Partial page should NOT contain html tag")
		require.Contains(t, body, "id=\"drafts-list\"", "Should contain drafts list container")
	})

	// --- ShowDiaryIndex ---
	t.Run("ShowDiaryIndex_Full", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/diary/%d", targetDiaryID), false, "", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "<html", "Full page should contain html tag")
		require.Contains(t, body, "id=\"diary-content\"", "Should contain diary content container")
	})

	t.Run("ShowDiaryIndex_Partial", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/diary/%d", targetDiaryID), true, "", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.NotContains(t, body, "<html", "Partial page should NOT contain html tag")
		require.Contains(t, body, "id=\"diary-content\"", "Should contain diary content container")
	})

	t.Run("ShowDiaryIndex_Partial_Page2", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/diary/%d?page=2", targetDiaryID), true, "", u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.NotContains(t, body, "<html", "Partial page should NOT contain html tag")
		require.Contains(t, body, "id=\"diary-content\"", "Should contain diary content container")
	})
}
