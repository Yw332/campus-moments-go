package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUploadAvatar 测试头像上传功能
func TestUploadAvatar(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	models.AutoMigrate()
	routes.SetupRoutes(router)

	// 添加静态文件服务
	router.Static("/uploads", "./uploads")

	// 注册并登录测试用户
	registerReq := map[string]string{
		"username": "avatartest",
		"phone":    "13800138010",
		"password": "AvatarTest123",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := getAuthToken(router, "avatartest", "AvatarTest123")

	// 创建测试图片文件
	testImageContent := []byte("fake image content for testing")
	
	tests := []struct {
		name           string
		filename       string
		fileContent    []byte
		token          string
		expectedCode   int
		expectedFields []string
	}{
		{
			name:           "成功上传JPG头像",
			filename:       "test_avatar.jpg",
			fileContent:    testImageContent,
			token:          token,
			expectedCode:   200,
			expectedFields: []string{"avatarUrl", "userId"},
		},
		{
			name:           "成功上传PNG头像",
			filename:       "test_avatar.png",
			fileContent:    testImageContent,
			token:          token,
			expectedCode:   200,
			expectedFields: []string{"avatarUrl", "userId"},
		},
		{
			name:         "未提供文件",
			filename:     "",
			fileContent:  []byte{},
			token:        token,
			expectedCode: 400,
		},
		{
			name:           "未认证上传",
			filename:       "test_avatar.jpg",
			fileContent:    testImageContent,
			token:          "",
			expectedCode:   401,
			expectedFields: []string{},
		},
		{
			name:         "不支持的文件格式",
			filename:     "test_avatar.txt",
			fileContent:  []byte("text content"),
			token:        token,
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备multipart表单数据
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if tt.filename != "" {
				part, err := writer.CreateFormFile("avatar", tt.filename)
				assert.NoError(t, err)
				_, err = part.Write(tt.fileContent)
				assert.NoError(t, err)
			}

			err := writer.Close()
			assert.NoError(t, err)

			req, _ := http.NewRequest("POST", "/api/upload/avatar", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedCode == 200 {
				assert.Equal(t, float64(200), response["code"])
				data := response["data"].(map[string]interface{})
				for _, field := range tt.expectedFields {
					assert.Contains(t, data, field)
				}
				// 验证头像URL格式
				avatarURL := data["avatarUrl"].(string)
				assert.Contains(t, avatarURL, "uploads/avatars")
				assert.Contains(t, avatarURL, "test_avatar")
			}
		})
	}
}

// TestGetUserProfileWithAvatar 测试获取用户资料时包含头像信息
func TestGetUserProfileWithAvatar(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	models.AutoMigrate()
	routes.SetupRoutes(router)

	// 添加静态文件服务
	router.Static("/uploads", "./uploads")

	// 注册测试用户
	registerReq := map[string]string{
		"username": "avatarprofiletest",
		"phone":    "13800138011",
		"password": "AvatarProfile123",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 登录获取token
	token := getAuthToken(router, "avatarprofiletest", "AvatarProfile123")

	// 先上传头像
	testImageContent := []byte("fake avatar image content")
	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("avatar", "profile_avatar.jpg")
	assert.NoError(t, err)
	_, err = part.Write(testImageContent)
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	req, _ = http.NewRequest("POST", "/api/upload/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 获取用户资料验证头像
	req, _ = http.NewRequest("GET", "/api/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])

	data := response["data"].(map[string]interface{})
	// 验证头像字段存在
	assert.Contains(t, data, "avatar")
	// 验证其他必要字段
	assert.Contains(t, data, "userId")
	assert.Contains(t, data, "username")
	assert.Contains(t, data, "phone")
}

// TestAvatarFileAccess 测试头像文件访问
func TestAvatarFileAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Static("/uploads", "./uploads")

	// 创建测试目录和文件
	uploadDir := "uploads/avatars"
	os.MkdirAll(uploadDir, 0755)
	defer os.RemoveAll("uploads")

	testFileContent := []byte("test avatar image")
	testFilePath := filepath.Join(uploadDir, "test_access.jpg")
	err := os.WriteFile(testFilePath, testFileContent, 0644)
	assert.NoError(t, err)

	// 测试文件访问
	req, _ := http.NewRequest("GET", "/uploads/avatars/test_access.jpg", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
	assert.Equal(t, testFileContent, w.Body.Bytes())

	// 测试不存在的文件
	req, _ = http.NewRequest("GET", "/uploads/avatars/nonexistent.jpg", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

// BenchmarkUploadAvatar 性能测试
func BenchmarkUploadAvatar(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	models.AutoMigrate()
	routes.SetupRoutes(router)
	router.Static("/uploads", "./uploads")

	// 注册并登录测试用户
	registerReq := map[string]string{
		"username": "benchtest",
		"phone":    "13800138012",
		"password": "BenchTest123",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	token := getAuthToken(router, "benchtest", "BenchTest123")

	testImageContent := []byte("benchmark test image content")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 准备multipart表单数据
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("avatar", fmt.Sprintf("bench_avatar_%d.jpg", i))
		if err != nil {
			b.Fatal(err)
		}
		_, err = part.Write(testImageContent)
		if err != nil {
			b.Fatal(err)
		}

		err = writer.Close()
		if err != nil {
			b.Fatal(err)
		}

		req, _ := http.NewRequest("POST", "/api/upload/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

// getAuthToken 辅助函数：注册并登录用户，返回token
func getAuthToken(router *gin.Engine, username, password string) string {
	// 登录请求
	loginReq := map[string]string{
		"account":  username,
		"password": password,
	}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	data := response["data"].(map[string]interface{})
	return data["token"].(string)
}