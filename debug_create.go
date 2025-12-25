package main

import (
	"fmt"
	"database/sql/driver"
	"encoding/json"
)

// Tags 标签数组类型
type Tags []string

// Value 实现 driver.Valuer 接口，用于数据库存储
func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	result, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	fmt.Printf("调试: Tags序列化 - 输入: %v, 输出: %s\n", t, string(result))
	return result, nil
}

// MediaItem 媒体项
type MediaItem struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Size     int64  `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"`
}

// CreateMomentRequest 创建动态请求
type CreateMomentRequest struct {
	Content    string             `json:"content"`
	Tags       []string           `json:"tags"`
	Media      []MediaItem        `json:"media"`
	Visibility int                `json:"visibility"`
}

func main() {
	// 模拟JSON请求
	jsonStr := `{
		"content": "测试带图片的动态发布",
		"tags": ["测试"],
		"media": [
			{
				"type": "image",
				"url": "http://106.52.165.122:8080/static/files/test_image_123.jpg"
			}
		],
		"visibility": 0
	}`

	var req CreateMomentRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		panic(err)
	}

	fmt.Printf("解析后的请求:\n")
	fmt.Printf("Content: %s\n", req.Content)
	fmt.Printf("Tags: %v\n", req.Tags)
	fmt.Printf("Media数量: %d\n", len(req.Media))
	for i, media := range req.Media {
		fmt.Printf("Media[%d]: Type=%s, URL=%s\n", i, media.Type, media.URL)
	}

	// 从media中提取图片URL
	var images []string
	fmt.Printf("\n开始提取图片:\n")
	for _, item := range req.Media {
		fmt.Printf("检查Media项: Type=%s, URL=%s\n", item.Type, item.URL)
		if item.Type == "image" && item.URL != "" {
			images = append(images, item.URL)
			fmt.Printf("✅ 提取图片: %s\n", item.URL)
		} else {
			fmt.Printf("❌ 跳过: Type=%s, URL=%s\n", item.Type, item.URL)
		}
	}

	fmt.Printf("\n最终提取的图片数组: %v\n", images)
	
	// 测试Tags序列化
	tags := Tags(images)
	result, err := tags.Value()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("最终要存储的images字段: %s\n", string(result.([]byte)))
}