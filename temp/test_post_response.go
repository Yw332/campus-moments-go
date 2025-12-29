package main

import (
	"encoding/json"
	"fmt"
)

// 模拟一个完整的帖子响应
type TestPostResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Images      []string `json:"images"`
	Cover       string   `json:"cover"`
	AuthorID    string   `json:"authorId"`
	LikeCount   int      `json:"likeCount"`
	CommentCount int     `json:"commentCount"`
	CreatedAt   string   `json:"createdAt"`
	Username    string   `json:"username"`
	Avatar      string   `json:"avatar"`
}

func main() {
	post := TestPostResponse{
		ID:           4,
		Title:        "动漫社招新啦!",
		Content:      "本周五下午3点，学生活动中心302室，动漫社招新活动。有cosplay展示、动漫放映、游戏体验等环节，欢迎所有",
		Images:       []string{"https://oss.example.com/posts/5-1.jpg", "https://oss.example.com/posts/5-2.jpg"},
		Cover:        "https://oss.example.com/posts/5-1.jpg", // 第一张图作为封面
		AuthorID:     "0000000005",
		LikeCount:    6,
		CommentCount: 15,
		CreatedAt:    "2025-12-28T18:30:00Z",
		Username:     "动漫社官方",
		Avatar:       "https://oss.example.com/avatars/anime-club.jpg",
	}
	
	data, _ := json.MarshalIndent(post, "", "  ")
	fmt.Println(string(data))
}