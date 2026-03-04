package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biuren/hebe/internal/app/controller/comment"
	"github.com/biuren/hebe/internal/app/middleware"
	"github.com/biuren/hebe/internal/app/view"
	commentDomain "github.com/biuren/hebe/internal/domain/comment"
	"github.com/biuren/hebe/internal/domain/diary"
	"github.com/biuren/hebe/internal/domain/entry"
	userDomain "github.com/biuren/hebe/internal/domain/user"
	"github.com/biuren/hebe/internal/infra/repository"
	"github.com/biuren/hebe/internal/utils/i18n"
	"github.com/biuren/hebe/tests/framework"
	"github.com/biuren/hebe/tests/framework/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetComments_Pagination_Button_Integration(t *testing.T) {
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
	commentService := commentDomain.NewService(config, commentRepo, userService, entryService, diaryService)

	// Create Helpers
	i18nManager := i18n.NewLocalizationManager(config)
	authHelper := userDomain.NewAuthHelper(config)
	renderService := view.NewViewService(config, i18nManager)

	// Initialize Controller
	commentCtrl := comment.NewCommentController(
		config,
		commentService,
		diaryService,
		entryService,
		authHelper,
		renderService,
	)

	// Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	helpers.LoadTestTemplates(r)

	// Register routes
	api := r.Group("/api")
	{
		commentGroup := api.Group("/comment")
		commentGroup.GET("/list",
			middleware.Validate[comment.GetCommentsRequest](),
			commentCtrl.GetComments,
		)
	}

	// 2. Prepare Data
	isolator := ctx.GetIsolator()

	// Create User
	u := &userDomain.User{
		Handle:      isolator.IsolateHandle("paginator"),
		Email:       isolator.IsolateEmail("paginator@example.com"),
		Pwd:         "Password123!",
		DisplayName: "Paginator User",
		Status:      3,
	}
	err := userService.RegisterUser(u)
	require.NoError(t, err)

	// Create Diary
	d := &diary.Diary{
		AddUserID:     u.UserId,
		Title:         isolator.IsolateTitle("Pagination Diary"),
		Remark:        "Test Pagination",
		ViewSetting:   diary.ViewSettingPublic, // Public
		SocialSetting: 1, // Enable social (comments)
		Status:        1,
	}
	_, err = diaryService.AddDiary(d)
	require.NoError(t, err)

	// Create Entry
	e := &entry.Entry{
		DiaryID:   d.DiaryID,
		AddUserID: u.UserId,
		Title:     isolator.IsolateTitle("Entry with comments"),
		Body:      "Content",
		Status:    1, // Published
		InfoType:  1,
	}
	_, err = entryService.AddEntry(e)
	require.NoError(t, err)

	// Determine page size from config
	pageSize := config.App.PageSize
	if pageSize <= 0 {
		pageSize = 10 // Default fallback if not set in config
	}

	// Create enough comments to have 2 pages
	totalComments := pageSize + 1
	for i := 0; i < totalComments; i++ {
		c := &commentDomain.Comment{
			EntryID:            e.EntryID,
			Content:            fmt.Sprintf("Comment %d", i+1),
			AddUserID:          u.UserId,
			AddUserHandle:      u.Handle,
			AddUserDisplayName: u.DisplayName,
			AddUserIP:          "127.0.0.1",
		}
		_, err := commentService.AddComment(c)
		require.NoError(t, err)
	}

	// 3. Run Tests

	// Case 1: First page, should have "Next" button
	req1 := httptest.NewRequest("GET", fmt.Sprintf("/api/comment/list?entry_id=%d&page_index=1", e.EntryID), nil)
	req1.Header.Set("HX-Request", "true")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	body1 := w1.Body.String()

	// Verify "Next" button exists with HTMX pagination and MDI icon
	assert.Contains(t, body1, "page_index=2")
	assert.Contains(t, body1, "hx-get=\"/api/comment/list")
	assert.Contains(t, body1, "hx-target=\"#comments-list-container\"")
	assert.Contains(t, body1, "mdi-chevron-right")

	// Verify "Prev" button does NOT exist on first page
	assert.NotContains(t, body1, "page_index=0")

	// Case 2: Second page, should have "<<" button
	req2 := httptest.NewRequest("GET", fmt.Sprintf("/api/comment/list?entry_id=%d&page_index=2", e.EntryID), nil)
	req2.Header.Set("HX-Request", "true")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	body2 := w2.Body.String()

	// Verify "Prev" button exists with HTMX pagination and MDI icon
	assert.Contains(t, body2, "page_index=1")
	assert.Contains(t, body2, "hx-get=\"/api/comment/list")
	assert.Contains(t, body2, "hx-target=\"#comments-list-container\"")
	assert.Contains(t, body2, "mdi-chevron-left")

	// Verify "Next" button does NOT exist on last page
	assert.NotContains(t, body2, "page_index=3")
}
