package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const baseURL = "http://127.0.0.1:8080"

func main() {
	fmt.Println("🚀 Starting API Tests...")

	// 1. Test Health
	testHealth()

	// 2. Test Register
	username := fmt.Sprintf("testuser_%d", time.Now().Unix())
	phone := fmt.Sprintf("138%08d", rand.Intn(100000000))
	password := "Test@1234"

	fmt.Printf("\n👤 Registering user: %s / %s\n", username, phone)
	registerUser(username, phone, password)

	// 3. Test Login
	fmt.Println("\n🔑 Logging in...")
	token, userID := loginUser(username, password)
	if token == "" {
		fmt.Println("❌ Login failed, aborting tests")
		return
	}
	fmt.Printf("✅ Login success! Token: %s... UserID: %v\n", token[:20], userID)

	// 4. Test Update Profile
	fmt.Println("\n📝 Testing Update Profile...")
	updateProfile(token, "Cool Nickname", "http://example.com/avatar.jpg", "Hello World Bio", 1)

	// 5. Test Get Profile
	fmt.Println("\n🔍 Testing Get Profile...")
	getProfile(token)

	// 6. Test Public Moments
	fmt.Println("\n📄 Testing Get Moments...")
	getMoments()
}

func testHealth() {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("✅ Health check status: %s\n", resp.Status)
}

func registerUser(username, phone, password string) {
	data := map[string]string{
		"username": username,
		"phone":    phone,
		"password": password,
	}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Register failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))
}

func loginUser(account, password string) (string, interface{}) {
	data := map[string]string{
		"account":  account,
		"password": password,
	}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		return "", nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Code int `json:"code"`
		Data struct {
			Token    string `json:"token"`
			UserInfo struct {
				UserID interface{} `json:"userId"`
			} `json:"userInfo"`
		} `json:"data"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("❌ Parse login response failed: %v\n", err)
		return "", nil
	}

	if result.Code != 200 {
		fmt.Printf("❌ Login failed: %s\n", result.Message)
		return "", nil
	}

	return result.Data.Token, result.Data.UserInfo.UserID
}

func updateProfile(token, nickname, avatar, bio string, gender int) {
	data := map[string]interface{}{
		"nickname": nickname,
		"avatar":   avatar,
		"bio":      bio,
		"gender":   gender,
	}
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("PUT", baseURL+"/api/users/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Update profile failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))
}

func getProfile(token string) {
	req, _ := http.NewRequest("GET", baseURL+"/api/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Get profile failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(body))
}

func getMoments() {
	resp, err := http.Get(baseURL + "/moments")
	if err != nil {
		fmt.Printf("❌ Get moments failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// Truncate if too long
	if len(body) > 500 {
		fmt.Printf("Response: %s...\n", string(body[:500]))
	} else {
		fmt.Printf("Response: %s\n", string(body))
	}
}
