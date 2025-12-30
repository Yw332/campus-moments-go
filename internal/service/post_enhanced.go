package service

import (
	"encoding/json"
	"github.com/Yw332/campus-moments-go/internal/models"
	"time"
)

// PostResponse 增强的帖子响应结构
type PostResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Images      []string `json:"images"`
	Cover       string    `json:"cover"`       // 封面图片（第一张图）
	Video       string    `json:"video"`
	AuthorID    string    `json:"authorId"`    // 作者ID
	LikeCount   int       `json:"likeCount"`
	CommentCount int      `json:"commentCount"`
	ViewCount   int       `json:"viewCount"`
	CreatedAt   time.Time `json:"createdAt"`
	
	// 用户信息
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// ConvertToPostResponse 将Post模型转换为响应格式
func ConvertToPostResponse(post models.Post) PostResponse {
	response := PostResponse{
		ID:           int64(post.ID),
		Title:        post.Title,
		Content:      post.Content,
		Video:        post.Video,
		AuthorID:     post.UserID,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		ViewCount:    post.ViewCount,
		CreatedAt:    post.CreatedAt,
	}
	
	// 处理图片和封面
	if post.Images != nil && len(post.Images) > 0 {
		var images []string
		if err := json.Unmarshal(post.Images, &images); err == nil {
			response.Images = images
			if len(images) > 0 {
				response.Cover = images[0] // 第一张图作为封面
			}
		}
	}
	
	// 处理用户信息
	if post.User != nil {
		response.Username = post.User.Username
		response.Avatar = post.User.Avatar
	}
	
	return response
}

// GetEnhancedPostList 获取增强版帖子列表
func GetEnhancedPostList(userID string, postType string, page, pageSize int) ([]PostResponse, int64, error) {
	posts, total, err := GetPostList(userID, postType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = ConvertToPostResponse(post)
	}
	
	return responses, total, nil
}

// GetEnhancedPostDetail 获取增强版帖子详情
func GetEnhancedPostDetail(postID int64, userID string) (*PostResponse, error) {
	post, err := GetPostDetail(postID, userID)
	if err != nil {
		return nil, err
	}
	
	response := ConvertToPostResponse(*post)
	return &response, nil
}