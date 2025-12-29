package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 测试数据结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=20"`
}

// 设置测试环境
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// 初始化测试数据库
	models.AutoMigrate()
	
	// 设置路由
	routes.SetupRoutes(router)
	
	return router
}

// TestUserRegistration 测试用户注册
func TestUserRegistration(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		request        RegisterRequest
		expectedCode   int
		expectedFields []string
	}{
		{
			name: "有效注册",
			request: RegisterRequest{
				Username: "testuser123",
				Phone:    "13800138000",
				Password: "TestPass123",
			},
			expectedCode:   200,
			expectedFields: []string{"userId", "username", "phone"},
		},
		{
			name: "用户名太短",
			request: RegisterRequest{
				Username: "ab",
				Phone:    "13800138001",
				Password: "TestPass123",
			},
			expectedCode: 400,
		},
		{
			name: "密码格式错误",
			request: RegisterRequest{
				Username: "testuser456",
				Phone:    "13800138002",
				Password: "12345678",
			},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedCode == 200 {
				assert.Equal(t, float64(200), response["code"])
				data := response["data"].(map[string]interface{})
				for _, field := range tt.expectedFields {
					assert.Contains(t, data, field)
				}
			}
		})
	}
}

// TestUserLogin 测试用户登录
func TestUserLogin(t *testing.T) {
	router := setupTestRouter()

	// 先注册一个测试用户
	registerReq := RegisterRequest{
		Username: "logintest",
		Phone:    "13800138003",
		Password: "LoginTest123",
	}
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	tests := []struct {
		name           string
		request        LoginRequest
		expectedCode   int
		expectToken    bool
	}{
		{
			name: "用户名登录成功",
			request: LoginRequest{
				Account:  "logintest",
				Password: "LoginTest123",
			},
			expectedCode: 200,
			expectToken:  true,
		},
		{
			name: "手机号登录成功",
			request: LoginRequest{
				Account:  "13800138003",
				Password: "LoginTest123",
			},
			expectedCode: 200,
			expectToken:  true,
		},
		{
			name: "密码错误",
			request: LoginRequest{
				Account:  "logintest",
				Password: "wrongpassword",
			},
			expectedCode: 401,
			expectToken:  false,
		},
		{
			name: "用户不存在",
			request: LoginRequest{
				Account:  "nonexistent",
				Password: "password123",
			},
			expectedCode: 404,
			expectToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectToken {
				data := response["data"].(map[string]interface{})
				assert.Contains(t, data, "token")
				assert.Contains(t, data, "userInfo")
			}
		})
	}
}

// TestGetUserProfile 测试获取用户资料
func TestGetUserProfile(t *testing.T) {
	router := setupTestRouter()

	// 注册并登录获取token
	registerReq := RegisterRequest{
		Username: "profiletest",
		Phone:    "13800138004",
		Password: "ProfileTest123",
	}
	
	// 注册
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 登录
	loginReq := LoginRequest{
		Account:  "profiletest",
		Password: "ProfileTest123",
	}
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 提取token
	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "有效token获取资料",
			token:        token,
			expectedCode: 200,
		},
		{
			name:         "无效token",
			token:        "invalid_token",
			expectedCode: 401,
		},
		{
			name:         "缺失token",
			token:        "",
			expectedCode: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/users/profile", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == 200 {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, float64(200), response["code"])
				data := response["data"].(map[string]interface{})
				assert.Contains(t, data, "userId")
				assert.Contains(t, data, "username")
				assert.Contains(t, data, "phone")
			}
		})
	}
}

// TestChangePassword 测试修改密码
func TestChangePassword(t *testing.T) {
	router := setupTestRouter()

	// 注册并登录获取token
	registerReq := RegisterRequest{
		Username: "changepass",
		Phone:    "13800138005",
		Password: "OldPassword123",
	}
	
	// 注册
	body, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 登录
	loginReq := LoginRequest{
		Account:  "changepass",
		Password: "OldPassword123",
	}
	body, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 提取token
	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	tests := []struct {
		name           string
		request        ChangePasswordRequest
		token          string
		expectedCode   int
	}{
		{
			name: "成功修改密码",
			request: ChangePasswordRequest{
				OldPassword: "OldPassword123",
				NewPassword: "NewPassword123",
			},
			token:        token,
			expectedCode: 200,
		},
		{
			name: "原密码错误",
			request: ChangePasswordRequest{
				OldPassword: "WrongPassword123",
				NewPassword: "NewPassword123",
			},
			token:        token,
			expectedCode: 400,
		},
		{
			name: "新密码格式错误",
			request: ChangePasswordRequest{
				OldPassword: "OldPassword123",
				NewPassword: "12345678",
			},
			token:        token,
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("PUT", "/api/users/password", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
		})
	}
}

// TestHealthCheck 测试健康检查
func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "status")
	assert.Contains(t, data, "service")
	assert.Contains(t, data, "database")
}