package post

import (
	"strconv"
	"strings"

	"github.com/biuren/hebe/internal/domain/post"
	errors "github.com/biuren/hebe/internal/utils/error"
	"github.com/gin-gonic/gin"
)

// parsePostID 解析并验证日记ID
// 从字符串中解析出有效的日记ID
//
// 参数:
//   - idStr: 待解析的ID字符串
//
// 返回值:
//   - int64: 解析后的ID
//   - error: 解析失败时的错误信息
func (ctrl *PostController) parsePostID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParams.Err()
	}
	return id, nil
}

// parseDiaryID 解析并验证日记本ID
// 从字符串中解析出有效的日记本ID
//
// 参数:
//   - idStr: 待解析的ID字符串
//
// 返回值:
//   - int64: 解析后的ID
//   - error: 解析失败时的错误信息
func (ctrl *PostController) parseDiaryID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidParams.Err()
	}
	return id, nil
}

// validatePostOwnership 验证日记所有权
// 检查指定日记是否存在，以及当前用户是否有权限操作该日记
//
// 参数:
//   - postID: 日记ID
//   - userID: 当前用户ID
//
// 返回值:
//   - *post.Post: 日记对象
//   - error: 验证失败时的错误信息
func (ctrl *PostController) validatePostOwnership(postID, userID int64) (*post.Post, error) {
	postData, err := ctrl.postService.GetPostByID(postID)
	if err != nil || postData == nil {
		return nil, errors.InvalidParams.Err()
	}

	// 检查是否是当前用户的日记
	if postData.AddUserID != userID {
		return nil, errors.PermissionDenied.Err()
	}

	return postData, nil
}

// getAndValidatePost 获取并验证日记权限（用于表单渲染）
// 从URL参数中获取post_id，验证其有效性和用户权限
//
// 参数:
//   - c: Gin上下文
//   - userID: 当前用户ID
//
// 返回值:
//   - *post.Post: 日记对象，如果验证失败返回nil
//   - error: 错误信息
func (ctrl *PostController) getAndValidatePost(c *gin.Context, userID int64) (*post.Post, error) {
	postID, err := ctrl.parsePostID(c.Query("post_id"))
	if err != nil {
		return nil, err
	}

	return ctrl.validatePostOwnership(postID, userID)
}