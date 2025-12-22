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
	fmt.Println("🚀 Starting Interaction Tests...")

	// 1. Setup Users
	userA := registerAndLogin("UserA")
	userB := registerAndLogin("UserB")

	if userA.Token == "" || userB.Token == "" {
		fmt.Println("❌ User setup failed")
		return
	}

	// 2. UserA creates a moment
	momentID := createMoment(userA.Token, "Hello World Moment")
	if momentID == 0 {
		fmt.Println("❌ Create moment failed")
		return
	}
	fmt.Printf("📝 UserA created moment %d\n", momentID)

	// 3. UserB likes the moment
	if toggleLike(userB.Token, momentID, 1) { // 1 = Moment
		fmt.Println("❤️ UserB liked the moment")
	} else {
		fmt.Println("❌ UserB like failed")
	}

	// 4. UserB comments on the moment
	commentID := createComment(userB.Token, momentID, "Nice moment!", nil, nil)
	if commentID == 0 {
		fmt.Println("❌ UserB comment failed")
	} else {
		fmt.Printf("💬 UserB commented: %d\n", commentID)
	}

	// 5. UserA replies to UserB's comment
	replyID := createComment(userA.Token, momentID, "Thanks!", &commentID, &userB.UserID)
	if replyID == 0 {
		fmt.Println("❌ UserA reply failed")
	} else {
		fmt.Printf("↩️ UserA replied: %d\n", replyID)
	}

	// 6. UserA likes UserB's comment
	if toggleLike(userA.Token, commentID, 2) { // 2 = Comment
		fmt.Println("❤️ UserA liked UserB's comment")
	} else {
		fmt.Println("❌ UserA like comment failed")
	}

	// 7. Verify comments
	verifyComments(userA.Token, momentID)
}

type UserSession struct {
	Token  string
	UserID string
}

func registerAndLogin(prefix string) UserSession {
	rand.Seed(time.Now().UnixNano())
	username := fmt.Sprintf("%s_%d", prefix, rand.Intn(100000))
	phone := fmt.Sprintf("139%08d", rand.Intn(100000000))
	password := "Test@1234"

	// Register
	data := map[string]string{"username": username, "phone": phone, "password": password}
	jsonData, _ := json.Marshal(data)
	http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))

	// Login
	loginData := map[string]string{"account": username, "password": password}
	jsonLogin, _ := json.Marshal(loginData)
	resp, _ := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonLogin))
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			Token    string `json:"token"`
			UserInfo struct {
				UserID string `json:"userId"`
			} `json:"userInfo"`
		} `json:"data"`
	}
	json.Unmarshal(body, &result)
	return UserSession{Token: result.Data.Token, UserID: result.Data.UserInfo.UserID}
}

func createMoment(token, content string) int64 {
	data := map[string]interface{}{
		"content":    content,
		"visibility": 0,
	}
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", baseURL+"/api/moments", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Create moment request failed: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("Create moment failed with status %d: %s\n", resp.StatusCode, string(body))
		return 0
	}

	var result struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(body, &result)
	return result.Data.ID
}

func toggleLike(token string, targetID int64, targetType int) bool {
	data := map[string]interface{}{
		"targetId":   targetID,
		"targetType": targetType,
	}
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", baseURL+"/api/likes", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func createComment(token string, momentID int64, content string, parentID *int64, replyToUserID *string) int64 {
	data := map[string]interface{}{
		"momentId": momentID,
		"content":  content,
	}
	if parentID != nil {
		data["parentId"] = *parentID
	}
	if replyToUserID != nil {
		data["replyToUserId"] = *replyToUserID
	}

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", baseURL+"/api/comments", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(body, &result)
	return result.Data.ID
}

func verifyComments(token string, momentID int64) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/moments/%d/comments", baseURL, momentID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("📄 Comments List:\n%s\n", string(body))
}
