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
	return json.Marshal(t)
}

// Scan 实现 sql.Scanner 接口，用于数据库读取
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	}
	return nil
}

func main() {
	// 测试Tags序列化
	images := []string{"http://106.52.165.122:8080/static/files/test_image_123.jpg"}
	tags := Tags(images)
	
	result, err := tags.Value()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Tags序列化结果: %s\n", string(result.([]byte)))
	
	// 测试创建Moment
	moment := struct {
		Images Tags `json:"images"`
	}{
		Images: tags,
	}
	
	jsonData, _ := json.MarshalIndent(moment, "", "  ")
	fmt.Printf("Moment JSON:\n%s\n", string(jsonData))
}