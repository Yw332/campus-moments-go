package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

const BASE_URL = "http://localhost:8080"

// 用户注册
func register() (string, error) {
	userData := map[string]interface{}{
		"username": "testuser" + strconv.FormatInt(time.Now().Unix(), 10),
		"phone":    "1380013" + strconv.FormatInt(time.Now().Unix()%10000, 10),
		"password": "Test123456",
	}

	jsonData, _ := json.Marshal(userData)
	resp, err := http.Post(BASE_URL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("注册响应: %s\n", string(body))
	return "", nil
}

// 用户登录
func login() (string, error) {
	loginData := map[string]interface{}{
		"account":  "testuser",
		"password": "Test123456",
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(BASE_URL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if token, ok := result["data"].(map[string]interface{})["token"].(string); ok {
		return token, nil
	}
	
	return "", fmt.Errorf("获取token失败")
}

// 发布动态
func createMoment(token string) error {
	momentData := map[string]interface{}{
		"content": "这是一条测试动态，包含标签功能 #" + strconv.FormatInt(time.Now().Unix(), 10),
		"tags":    []string{"测试", "动态", "标签"},
		"media": []map[string]interface{}{
			{
				"url":      "https://example.com/test1.jpg",
				"type":     "image",
				"size":     1024000,
				"width":    1920,
				"height":   1080,
				"duration": 0,
			},
		},
		"visibility": 0,
	}

	jsonData, _ := json.Marshal(momentData)
	req, _ := http.NewRequest("POST", BASE_URL+"/api/moments", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("发布动态响应: %s\n", string(body))
	return nil
}

// 获取动态列表
func getMoments() error {
	resp, err := http.Get(BASE_URL + "/moments?page=1&pageSize=5")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("动态列表响应: %s\n", string(body))
	return nil
}

// 获取用户动态列表
func getUserMoments(token string) error {
	req, _ := http.NewRequest("GET", BASE_URL+"/api/moments/my?page=1&pageSize=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("用户动态列表响应: %s\n", string(body))
	return nil
}

func main() {
	fmt.Println("=== 动态功能测试 ===")

	// 1. 注册用户
	fmt.Println("\n1. 注册用户...")
	if err := register(); err != nil {
		log.Printf("注册失败: %v", err)
	}

	// 2. 登录获取token
	fmt.Println("\n2. 用户登录...")
	token, err := login()
	if err != nil {
		log.Printf("登录失败: %v", err)
		return
	}
	fmt.Printf("获取到token: %s...\n", token[:min(len(token), 50)])

	// 3. 发布动态
	fmt.Println("\n3. 发布动态...")
	if err := createMoment(token); err != nil {
		log.Printf("发布动态失败: %v", err)
	}

	// 4. 获取动态列表
	fmt.Println("\n4. 获取动态列表...")
	if err := getMoments(); err != nil {
		log.Printf("获取动态列表失败: %v", err)
	}

	// 5. 获取用户动态列表
	fmt.Println("\n5. 获取用户动态列表...")
	if err := getUserMoments(token); err != nil {
		log.Printf("获取用户动态列表失败: %v", err)
	}

	fmt.Println("\n=== 测试完成 ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}