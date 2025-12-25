# 校园动态 API 接口文档

## 基础信息
- **服务器地址**: `http://localhost:8080`
- **认证方式**: Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

---

## 🔐 认证接口

### 用户登录
```http
POST /auth/login
Content-Type: application/json
```

**请求参数**:
```json
{
  "account": "13800138001",
  "password": "password"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "userInfo": {
      "userId": "0000000001",
      "username": "新用户名456",
      "phone": "13800138001",
      "avatar": "http://example.com/avatar.jpg"
    }
  }
}
```

---

## 📝 动态相关接口

### 获取动态列表
```http
GET /moments?page=1&pageSize=10
```

**参数**:
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认10，最大100）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 28,
        "title": "测试发布动态功能",
        "content": "测试发布动态功能",
        "images": ["http://localhost:8080/static/files/test.jpg"],
        "authorId": "0000000001",
        "likeCount": 0,
        "commentCount": 0,
        "createdAt": "2025-12-25T10:41:01.179+08:00"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 17
    }
  }
}
```

### 发布动态
```http
POST /api/moments
Authorization: Bearer {token}
Content-Type: application/json
```

**请求参数**:
```json
{
  "content": "这是我的动态内容",
  "tags": ["生活", "学习"],
  "media": [
    {
      "type": "image",
      "url": "http://localhost:8080/static/files/test.jpg"
    }
  ],
  "visibility": 0
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "id": 29,
    "content": "这是我的动态内容",
    "images": ["http://localhost:8080/static/files/test.jpg"],
    "status": 1,
    "createdAt": "2025-12-25T11:03:39.092+08:00"
  }
}
```

### 获取我的动态
```http
GET /api/moments/my?page=1&pageSize=10
Authorization: Bearer {token}
```

---

## 📁 文件上传接口

### 上传普通文件
```http
POST /api/upload/file
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:
- `file`: 文件（表单字段名）

**响应示例**:
```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "fileId": "a8b27dd0-2d09-42",
    "filename": "20251225113508_a8b27dd0.jpg",
    "originalName": "test_image.jpg",
    "fileSize": 1024,
    "fileType": ".jpg",
    "fileUrl": "http://localhost:8080/static/files/20251225113508_a8b27dd0.jpg"
  }
}
```

### 上传头像
```http
POST /api/upload/avatar
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:
- `avatar`: 图片文件

**响应示例**:
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatarUrl": "http://localhost:8080/static/avatars/avatar.jpg",
    "filename": "avatar.jpg",
    "size": 1024
  }
}
```

---

## 🔍 搜索接口

### 搜索内容
```http
GET /search?keyword=学习&page=1&pageSize=10&sortBy=latest
```

**参数**:
- `keyword`: 搜索关键词（必填）
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认10）
- `sortBy`: 排序方式（latest/latest/hottest/comprehensive）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "moments": [
      {
        "id": 1,
        "title": "数据结构期末复习资料分享",
        "content": "整理了一学期的数据结构笔记...",
        "images": ["http://localhost:8080/static/files/1-1.jpg"],
        "authorId": "0000000001",
        "likeCount": 5,
        "commentCount": 8
      }
    ],
    "users": [
      {
        "userId": "0000000001",
        "username": "新用户名456",
        "avatar": "http://example.com/avatar.jpg"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 2
    }
  }
}
```

### 获取热词
```http
GET /search/hot-words
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": ["校园", "活动", "学习", "美食", "运动", "兼职", "考试", "社团", "室友", "考研"]
}
```

### 获取搜索建议
```http
GET /search/suggestions?keyword=学
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": ["新用户名456"]
}
```

---

## 👤 用户信息接口

### 获取用户资料
```http
GET /api/users/profile
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "userId": "0000000001",
    "username": "新用户名456",
    "phone": "13800138001",
    "avatar": "http://example.com/avatar.jpg",
    "signature": "最终测试",
    "status": 1
  }
}
```

### 更新用户资料
```http
PUT /api/users/profile
Authorization: Bearer {token}
Content-Type: application/json
```

**请求参数**:
```json
{
  "username": "新用户名",
  "signature": "个性签名",
  "avatar": "头像URL"
}
```

---

## 💬 互动接口

### 切换点赞状态
```http
POST /api/likes
Authorization: Bearer {token}
Content-Type: application/json
```

**请求参数**:
```json
{
  "momentId": 1,
  "isLiked": true
}
```

### 获取点赞列表
```http
GET /api/likes?momentId=1&page=1&pageSize=10
Authorization: Bearer {token}
```

---

## 📊 响应格式说明

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": { /* 具体数据 */ }
}
```

### 错误响应
```json
{
  "code": 400/401/404/500,
  "message": "错误说明",
  "data": null
}
```

### 分页格式
```json
{
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 100
  }
}
```

---

## 🚨 重要说明

1. **文件访问**: 所有上传的文件都通过 `http://localhost:8080/static/files/` 访问
2. **认证**: 除标明外的接口都需要在 Header 中携带 `Authorization: Bearer {token}`
3. **时间格式**: 统一使用 ISO 8601 格式
4. **排序**: 大部分列表按创建时间倒序排列
5. **状态码**: 200 成功，400 参数错误，401 未认证，404 资源不存在，500 服务器错误

---

## 🔧 开发调试

### 健康检查
```http
GET /health
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "service": "Campus Moments Go API",
    "database": "connected",
    "status": "ok"
  },
  "message": "success"
}
```