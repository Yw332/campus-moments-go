package main

import (
	"encoding/json"
	"fmt"
	"database/sql/driver"
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
	fmt.Printf("Tags.Value 输入: %v, 输出: %s\n", t, string(result))
	return result, nil
}

// MediaItem 媒体项
type MediaItem struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	Width int    `json:"width"`
	Height int   `json:"height"`
}

func main() {
	// 模拟CreateMomentRequest
	req := struct {
		Content string        `json:"content"`
		Tags    []string      `json:"tags"`
		Media   []MediaItem  `json:"media"`
	}{
		Content: "测试带图片的动态发布",
		Tags:    []string{"测试"},
		Media: []MediaItem{
			{Type: "image", URL: "http://106.52.165.122:8080/static/files/test_image_123.jpg"},
		},
	}

	// 从media中提取图片URL
	var images []string
	for _, item := range req.Media {
		if item.Type == "image" && item.URL != "" {
			images = append(images, item.URL)
			fmt.Printf("提取图片: %s\n", item.URL)
		}
	}

	fmt.Printf("提取的图片数组: %v\n", images)
	
	// 测试Tags序列化
	tags := Tags(images)
	result, err := tags.Value()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("最终存储到数据库的images字段: %s\n", string(result.([]byte)))
}