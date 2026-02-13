package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestLatestPublicEntries_Integration(t *testing.T) {
	// 1. Setup Environment
	ctx := framework.NewTestContext(framework.WithDatabase())
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

	// Create Diary for u1 (Public)
	d := dataGen.CreateDiary(u1.UserId)
	d.Title = "Public Diary"
	d.ViewSetting = 1 // Public (matches repository filter)
	d.Status = 1      // Active
	diaryID, err := diaryService.AddDiary(d)
	require.NoError(t, err)

	// Create Public Entry
	e := dataGen.CreateEntry(u1.UserId, diaryID)
	e.Title = "Public Entry"
	e.Status = 1 // Active
	entryID1, err := entryService.AddEntry(e)
	require.NoError(t, err)
	e.EntryID = entryID1

	// Create another Diary for u1 (Public but Social Disabled)
	d2 := dataGen.CreateDiary(u1.UserId)
	d2.Title = "Social Disabled Diary"
	d2.ViewSetting = 1   // Public
	d2.Status = 1        // Active
	d2.SocialSetting = 2 // Disabled
	diaryID2, err := diaryService.AddDiary(d2)
	require.NoError(t, err)

	// Create Entry in Social Disabled Diary
	e2 := dataGen.CreateEntry(u1.UserId, diaryID2)
	e2.Title = "Social Disabled Entry"
	e2.Status = 1
	entryID2, err := entryService.AddEntry(e2)
	require.NoError(t, err)
	e2.EntryID = entryID2

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
	diaryCtrl.RegisterRoutes(router)
	entryCtrl.RegisterRoutes(router)

	// Helper to make request
	makeRequest := func(method, path string, userID int64) *httptest.ResponseRecorder {
		req, _ := http.NewRequest(method, path, nil)
		req.Header.Set("X-Test-User-ID", fmt.Sprintf("%d", userID))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}

	// 4. Test Cases

	t.Run("Owner Views Own Latest Entries", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d/latest", u1.UserId), u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()

		// Verify IsOwner logic (Edit button presence)
		// The edit button is inside <div class="entry-edit">
		require.Contains(t, body, "entry-edit", "Should show edit button for owner")

		// Verify Social Setting Logic
		// With new entryStats component, social setting is controlled by {{if eq .SocialSetting 1}}
		// Check for entryStats component initialization for enabled diary
		require.Contains(t, body, `x-data="entryStats(`, "Should render entryStats component for entries with social enabled")
		require.Contains(t, body, "x-init=\"initStats()\"", "Should call initStats() for entryStats component")

		// Verify entries are present
		require.Contains(t, body, "Public Entry")
		require.Contains(t, body, "Social Disabled Entry")

		// Verify specific "Edit" menu item presence
		// The edit menu item is inside {{if eq .CurrentUserID .Entry.AddUserID}}
		// It contains "mdi-pencil-outline"
		require.Contains(t, body, "mdi-pencil-outline", "Should show edit menu item (pencil icon) for owner")
	})

	t.Run("ShowLatestPublicEntries_Success", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d/latest", u1.UserId), u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "<html", "Full page should contain html tag")
		require.Contains(t, body, "Public Entry", "Should contain public entry title")
		require.Contains(t, body, `x-data="diaryDetail"`, "Should contain Alpine component initialization")
	})

	t.Run("ShowLatestPublicEntries_UserHandleFallback", func(t *testing.T) {
		// Create User u3 with Handle but no DisplayName
		u3 := dataGen.CreateUser(
			helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u3_handle")),
			helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u3@test.com")),
			helpers.WithDisplayName(""), // Explicitly empty display name
		)
		err := userService.RegisterUser(u3)
		require.NoError(t, err)

		// Create a diary and entry for u3
		d3 := dataGen.CreateDiary(u3.UserId, helpers.WithViewSetting(1)) // Public
		d3.Status = 1
		diaryID, err := diaryService.AddDiary(d3)
		require.NoError(t, err)
		
		e3 := dataGen.CreateEntry(u3.UserId, diaryID)
		e3.Status = 1
		_, err = entryService.AddEntry(e3)
		require.NoError(t, err)

		// Mock Auth for u3
		users[fmt.Sprintf("%d", u3.UserId)] = &types.AuthUser{
			UserId:      u3.UserId,
			DisplayName: "",
			UserGroup:   u3.UserGroup,
			Status:      u3.Status,
			Lang:        "en",
		}

		w := makeRequest("GET", fmt.Sprintf("/u/%d/latest", u3.UserId), u3.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, u3.Handle, "Should display handle when display name is empty")
		require.NotContains(t, body, ">User<", "Should NOT fallback to 'User' when handle is present")
	})

	t.Run("ShowLatestPublicEntries_Empty", func(t *testing.T) {
		// Create User u2 with no entries
		u2 := dataGen.CreateUser(
			helpers.WithHandle(ctx.GetIsolator().IsolateHandle("u2")),
			helpers.WithEmail(ctx.GetIsolator().IsolateEmail("u2@test.com")),
		)
		err := userService.RegisterUser(u2)
		require.NoError(t, err)

		// Mock Auth for u2
		users[fmt.Sprintf("%d", u2.UserId)] = &types.AuthUser{
			UserId:      u2.UserId,
			DisplayName: u2.DisplayName,
			UserGroup:   u2.UserGroup,
			Status:      u2.Status,
			Lang:        "en",
		}
		
		w := makeRequest("GET", fmt.Sprintf("/u/%d/latest", u2.UserId), u2.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "empty-placeholder", "Should contain empty placeholder")
		require.Contains(t, body, "mdi-inbox-outline", "Should contain empty icon")
	})

	t.Run("ShowUserProfile_NoLatestPublicEntries", func(t *testing.T) {
		w := makeRequest("GET", fmt.Sprintf("/u/%d", u1.UserId), u1.UserId)
		require.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		require.Contains(t, body, "<html", "Full page should contain html tag")
		// Ideally we check that the section is missing. 
		// The section had "user.profile.latest_public_diary" in the header (Latest Public Diary), 
		// but that string is also used in the link.
		// The link should be present.
		require.Contains(t, body, fmt.Sprintf("/u/%d/latest", u1.UserId), "Should contain link to latest public entries")
		
		// We can check that the entry title is NOT present in the main part of the page, 
		// assuming it's not listed in the "Diaries" list (which it shouldn't be, as that lists Diaries, not Entries)
		// Wait, ShowUserProfile lists "Diaries", not "Entries". 
		// So "Public Entry" string should NOT be present if the "Latest Public Entries" section is removed.
		require.NotContains(t, body, "Public Entry", "Should NOT contain public entry title on profile page")
	})
}
