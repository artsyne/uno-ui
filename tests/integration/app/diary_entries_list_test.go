package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/biuren/hebe/internal/app/controller/entry"
	"github.com/biuren/hebe/internal/app/view"
	"github.com/biuren/hebe/internal/domain/comment"
	diary "github.com/biuren/hebe/internal/domain/diary"
	entryDomain "github.com/biuren/hebe/internal/domain/entry"
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

func TestDiaryEntriesList_MultipleAuthors(t *testing.T) {
	// 1. Setup Environment
	ctx := framework.NewTestContext(framework.WithDatabase())
	defer ctx.Cleanup()

	config := ctx.GetConfig().GetAppConfig()

	// Get Services from Registry
	userService := ctx.GetServices().GetUserService()
	diaryService := ctx.GetServices().GetDiaryService()
	entryService := ctx.GetServices().GetEntryService()

	// Manually create CommentService (if not in registry)
	commentRepo := repository.NewCommentRepository()
	commentService := comment.NewService(config, commentRepo, userService, entryService, diaryService)

	// Manually create LikeService (not in registry yet)
	likeRepo := repository.NewLikeRepository()
	likeService := like.NewLikeService(likeRepo)

	// Create Helpers
	i18nManager := i18n.NewLocalizationManager(config)
	authHelper := userDomain.NewAuthHelper(config)
	renderService := view.NewViewService(config, i18nManager)

	// Initialize Controller
	entryCtrl := entry.NewEntryController(config, entryService, diaryService, commentService, authHelper, renderService, userService, likeService)

	// 2. Prepare Data
	dataGen := ctx.GetDataGenerator()

	// Create Owner
	owner := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("owner")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("owner@test.com")),
	)
	err := userService.RegisterUser(owner)
	require.NoError(t, err)

	// Create Author
	author := dataGen.CreateUser(
		helpers.WithHandle(ctx.GetIsolator().IsolateHandle("author")),
		helpers.WithEmail(ctx.GetIsolator().IsolateEmail("author@test.com")),
	)
	err = userService.RegisterUser(author)
	require.NoError(t, err)

	// Create Diary
	diaryObj := dataGen.CreateDiary(owner.UserId)
	diaryObj.ViewSetting = diary.ViewSettingPublic // Public
	diaryID, err := diaryService.AddDiary(diaryObj)
	require.NoError(t, err)

	// Add Author to Diary
	// Invite
	err = diaryService.InviteAuthor(diaryID, owner.UserId, author.Handle, "en")
	require.NoError(t, err)
	// Accept
	err = diaryService.AcceptInvitation(diaryID, author.UserId)
	require.NoError(t, err)

	// Create Entry by Owner
	entry1 := dataGen.CreateEntry(owner.UserId, diaryID)
	entry1.Title = "Owner Entry Title"
	entry1.Status = entryDomain.StatusActive // Active
	entry1.AddUTCTime = time.Now().Unix() - 1000 // Older
	_, err = entryService.AddEntry(entry1)
	require.NoError(t, err)

	// Create Entry by Author
	entry2 := dataGen.CreateEntry(author.UserId, diaryID)
	entry2.Title = "Author Entry Title"
	entry2.Status = entryDomain.StatusActive // Active
	entry2.AddUTCTime = time.Now().Unix() // Newer
	_, err = entryService.AddEntry(entry2)
	require.NoError(t, err)

	// 3. Setup Router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	helpers.LoadTestTemplates(router)
	
	// Register Routes
	entryCtrl.RegisterRoutes(router)

	// 4. Test Case: List Diary Entries
	t.Run("ListDiaryEntries_ShouldShowAllAuthors", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/diary/%d", diaryID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		body := w.Body.String()
		// Should contain both titles
		assert.True(t, strings.Contains(body, "Owner Entry Title"), "Response should contain owner entry title")
		assert.True(t, strings.Contains(body, "Author Entry Title"), "Response should contain author entry title")
	})
}
