package services

import (
	"github.com/biuren/hebe/internal/configs"
	"github.com/biuren/hebe/internal/errors"
	"github.com/biuren/hebe/internal/models"
	"github.com/biuren/hebe/pkg"
)

// 增加新的日记（post）
// 返回新增日记ID
func AddPost(post *models.Post) (int64, error) {

	//截取配置长度字符串，末尾增加 suffix
	post.Summary = pkg.TruncateString(post.Body, configs.GetConfig().App.SummaryMax, configs.GetConfig().App.Suffix)

	// todo：需判断空间是否有权限发布日记

	//插入数据库
	postID, err := post.Insert()
	if err != nil {
		return 0, errors.DbOperationErr
	}
	post.PostID = postID

	//日记本信息数量加一
	err = models.IncrementPostNum(post.DiaryID, post.AddUserID, post.AddUserName, post.AddUTCTime)
	if err != nil {
		return 0, errors.DbOperationErr
	}

	return post.PostID, nil
}

// GetPostById 根据id获取post
func GetPostByID(postId int64) (*models.Post,error){
	post,error :=  models.GetPostByID(postId)
	return post,error
}

// UpdatePost 更新 post
func UpdatePost(post *models.Post)(error){
	//校验参数：todo


	// 更新日记本
	err := post.Update()
	if err != nil {
		return err
	}

	return nil
}