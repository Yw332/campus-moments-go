package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// 直接提供一个简化的搜索响应
func main() {
	keyword := "数据结构"
	
	// 模拟搜索结果
	if strings.Contains(keyword, "数据结构") {
		result := map[string]interface{}{
			"code": 200,
			"message": "success", 
			"data": map[string]interface{}{
				"moments": []map[string]interface{}{
					{
						"id": 1,
						"title": "数据结构期末复习资料分享",
						"content": "整理了一学期的数据结构笔记，包含所有重要知识点和例题解析，需要的同学可以下载。希望可以帮助到大家！",
						"images": []string{
							"https://oss.example.com/posts/1-1.jpg",
							"https://oss.example.com/posts/1-2.jpg",
						},
						"authorId": "0000000001",
						"likeCount": 5,
						"commentCount": 8,
						"createdAt": time.Now().Format(time.RFC3339),
					},
				},
				"users": []map[string]interface{}{
					{
						"userId": "0000000001",
						"username": "新用户名456",
						"avatar": "http://example.com/avatar.jpg",
					},
				},
				"pagination": map[string]int{
					"page":     1,
					"pageSize": 10,
					"total":    2,
				},
			},
		}
		
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println("修复后的搜索响应示例:")
		fmt.Println(string(jsonData))
	}
}