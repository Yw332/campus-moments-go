# 🌐 Campus Moments Go API 接口文档

## 📋 基础信息

- **Base URL**: `http://106.52.165.122:8080`
- **响应格式**: 统一 JSON 格式
- **状态码说明**: 所有接口成功时返回 `code: 200`

### 🎯 响应格式标准

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": { /* 具体数据 */ }
}
```

#### 错误响应
```json
{
  "code": 400|401|404|500,
  "message": "错误描述",
  "data": null
}
```

---

## 🔐 认证接口

### 1. 用户登录
```
POST /auth/login
```

**请求体**:
```json
{
  "account": "手机号或用户名",
  "password": "密码"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "JWT token",
    "userInfo": {
      "userId": "用户ID",
      "username": "用户名",
      "phone": "手机号"
    }
  }
}
```

---

## 🔍 搜索功能

### 1. 关键词搜索
```
GET /api/search?keyword=关键词&sortBy=排序方式
Authorization: Bearer {token}
```

**参数**:
- `keyword`: 搜索关键词（必填）
- `sortBy`: 排序方式（可选）
  - `latest`: 最新（默认）
  - `hottest`: 最热
  - `comprehensive`: 综合

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "moments": [
      {
        "id": 1,
        "content": "动态内容",
        "authorId": "用户ID",
        "likeCount": 点赞数,
        "commentCount": 评论数,
        "createdAt": "2025-12-23T13:00:00Z"
      }
    ],
    "users": [
      {
        "id": 1,
        "username": "用户名",
        "phone": "手机号"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 19
    }
  }
}
```

### 2. 获取热词
```
GET /search/hot-words
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": ["校园", "活动", "学习", "美食", "运动"]
}
```

### 3. 获取搜索历史
```
GET /api/search/history
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": ["搜索词1", "搜索词2", "搜索词3"]
}
```

### 4. 保存搜索历史
```
POST /api/search/history
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "keyword": "搜索关键词"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "保存成功",
  "data": null
}
```

---

## 📱 用户相关

### 1. 获取用户资料
```
GET /api/users/profile
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "username": "用户名",
    "phone": "手机号",
    "avatar": "头像URL"
  }
}
```

---

## 📝 动态管理

### 1. 获取动态列表
```
GET /moments
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "content": "动态内容",
      "author": {
        "id": 1,
        "username": "用户名",
        "avatar": "头像URL"
      },
      "tags": ["标签1", "标签2"],
      "media": [
        {
          "url": "媒体URL",
          "type": "image|video",
          "size": 文件大小,
          "width": 宽度,
          "height": 高度
        }
      ],
      "likeCount": 点赞数,
      "commentCount": 评论数,
      "createdAt": "2025-12-23T13:00:00Z"
    }
  ]
}
```

### 2. 创建动态
```
POST /api/moments
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "content": "动态内容",
  "tags": ["标签1", "标签2"],
  "visibility": 0, // 0公开, 1好友, 2私密
  "media": [
    {
      "url": "媒体URL",
      "type": "image|video"
    }
  ]
}
```

---

## 💬 互动功能

### 1. 点赞/取消点赞
```
POST /api/likes
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "targetId": 1,
  "targetType": "moment", // moment|comment
  "action": "like" // like|unlike
}
```

### 2. 发表评论
```
POST /api/comments
Authorization: Bearer {token}
```

**请求体**:
```json
{
  "momentId": 1,
  "content": "评论内容",
  "replyToId": null // 回复评论时填评论ID
}
```

---

## 🖼️ 文件上传

### 1. 上传文件
```
POST /api/upload/file
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**参数**:
- `file`: 文件（支持图片和视频）

**响应**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "文件访问URL",
    "type": "image|video",
    "size": 文件大小
  }
}
```

---

## 🔧 系统接口

### 1. 健康检查
```
GET /health
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "ok",
    "service": "Campus Moments Go API",
    "database": "connected"
  }
}
```

### 2. 首页信息
```
GET /
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "Campus Moments Go API 运行中",
    "version": "1.0.0",
    "author": "Yw332",
    "status": "healthy"
  }
}
```

---

## ⚠️ 错误码说明

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未认证或token过期 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 📝 开发注意事项

1. **认证**: 除了 `/health`、`/`、`/search/hot-words`、`/moments`、`/auth/login` 外，其他接口都需要在Header中携带 `Authorization: Bearer {token}`

2. **时间格式**: 所有时间字段使用 ISO 8601 格式: `2025-12-23T13:00:00Z`

3. **分页**: 列表类接口支持 `page` 和 `pageSize` 参数，默认为 `page=1&pageSize=10`

4. **异步操作**: 搜索时会自动保存搜索历史，无需额外调用保存接口

5. **热词**: 热词基于当日用户搜索频率统计，无搜索历史时返回默认热词

---

**🎉 所有接口响应格式已统一为 `code: 200`，可放心对接！**