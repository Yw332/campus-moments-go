# 🌐 Campus Moments API 完整对接文档

## 📡 服务器信息
- **服务器地址**: `http://106.52.165.122:8080`
- **健康检查**: `GET /health`
- **API版本**: v1.0.0

## 🔐 认证方式
所有 `/api/*` 接口都需要在请求头中携带JWT Token：
```javascript
headers: {
  'Content-Type': 'application/json',
  'Authorization': 'Bearer <your_token_here>'
}
```

---

## 📱 用户认证模块

### 1. 用户注册
**POST** `/auth/register`

**请求参数**:
```json
{
  "phone": "17875242005",
  "password": "yourpassword",
  "username": "testuser",
  "nickname": "测试用户"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "token": "jwt_token_here",
    "userInfo": {
      "id": 1,
      "phone": "17875242005",
      "nickname": "测试用户"
    }
  }
}
```

### 2. 用户登录
**POST** `/auth/login`

**请求参数**:
```json
{
  "account": "17875242005",
  "password": "yourpassword"
}
```

### 3. 🆕 忘记密码 - 发送验证码
**POST** `/auth/send-verification`

**请求参数**:
```json
{
  "phone": "17875242005"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "验证码发送成功",
  "data": {
    "phone": "17875242005",
    "expiresIn": 300,
    "resendAfter": 60
  }
}
```

### 4. 🆕 忘记密码 - 验证并重置密码
**POST** `/auth/verify-and-reset`

**请求参数**:
```json
{
  "phone": "17875242005",
  "verificationCode": "123456",
  "newPassword": "newpassword123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "密码重置成功",
  "data": null
}
```

---

## 👤 用户资料模块

### 1. 获取用户资料
**GET** `/api/users/profile`

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "username": "test_user",
    "nickname": "测试用户",
    "avatar": "http://106.52.165.122:8080/static/avatars/xxx.jpg",
    "bio": "这个人很懒，什么都没留下",
    "phone": "17875242005"
  }
}
```

### 2. 更新用户资料
**PUT** `/api/users/profile`

**请求参数**:
```json
{
  "nickname": "新昵称",
  "bio": "个人简介",
  "avatar": "头像URL"
}
```

### 3. 修改密码
**PUT** `/api/users/password`

**请求参数**:
```json
{
  "oldPassword": "oldpassword",
  "newPassword": "newpassword123"
}
```

---

## 📤 文件上传模块

### 1. 📸 上传头像
**POST** `/api/upload/avatar`

**请求格式**: `multipart/form-data`
- **字段名**: `avatar`
- **文件类型**: jpg, jpeg, png, gif, webp
- **文件大小**: 最大5MB

**JavaScript示例**:
```javascript
async function uploadAvatar(file) {
  const formData = new FormData();
  formData.append('avatar', file);

  const response = await fetch('http://106.52.165.122:8080/api/upload/avatar', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    },
    body: formData
  });

  const result = await response.json();
  if (result.code === 200) {
    console.log('头像上传成功:', result.data.avatarUrl);
  }
}
```

**响应**:
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatarUrl": "http://106.52.165.122:8080/static/avatars/20231216_12345678_uuid.jpg",
    "filename": "20231216_12345678_uuid.jpg",
    "size": 1024000
  }
}
```

### 2. 通用文件上传
**POST** `/api/upload/file`

---

## 📝 动态内容管理模块

### 1. 发布动态
**POST** `/api/moments`

**请求参数**:
```json
{
  "content": "这是我的第一条动态",
  "tags": ["学习", "生活", "校园"],
  "media": [
    {
      "type": "image",
      "url": "http://example.com/image.jpg",
      "description": "图片描述"
    }
  ],
  "visibility": 0
}
```

**参数说明**:
- `content`: 动态内容（必填）
- `tags`: 标签数组（可选，支持多标签）
- `media`: 媒体文件数组（可选）
- `visibility`: 可见性（0=公开，1=好友，2=私密）

### 2. 获取动态列表
**GET** `/moments`

**查询参数**:
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认10）
- `userId`: 指定用户ID（可选）

### 3. 📋 获取我的动态列表
**GET** `/api/moments/my`

**查询参数**:
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认10）

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "content": "我的动态内容",
        "tags": ["学习", "生活"],
        "author": {
          "id": 1,
          "nickname": "测试用户",
          "avatar": "头像URL"
        },
        "likeCount": 10,
        "commentCount": 5,
        "createdAt": "2023-12-16T14:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 25
    }
  }
}
```

### 4. 获取动态详情
**GET** `/api/moments/:id`

### 5. 编辑动态
**PATCH** `/api/moments/:id`

### 6. 🗑️ 删除动态
**DELETE** `/api/moments/:id`

**响应**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {
    "postId": 123
  }
}
```

---

## 🔍 搜索功能模块

### 1. 搜索内容
**GET** `/api/search?keyword=关键词`

### 2. 获取热门关键词
**GET** `/api/search/hot-words`

### 3. 搜索历史
**GET** `/api/search/history`

---

## 🧪 完整测试流程

### 1. 测试用户注册
```javascript
// 注册
const registerResponse = await fetch('http://106.52.165.122:8080/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    phone: '17875242005',
    password: 'password123',
    username: 'testuser',
    nickname: '测试用户'
  })
});
```

### 2. 测试忘记密码流程
```javascript
// 1. 发送验证码
await fetch('http://106.52.165.122:8080/auth/send-verification', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ phone: '17875242005' })
});

// 2. 验证并重置密码
await fetch('http://106.52.165.122:8080/auth/verify-and-reset', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    phone: '17875242005',
    verificationCode: '123456',
    newPassword: 'newpassword123'
  })
});
```

### 3. 测试动态管理
```javascript
// 1. 发布带标签的动态
const createResponse = await fetch('http://106.52.165.122:8080/api/moments', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer ' + token
  },
  body: JSON.stringify({
    content: '今天天气真好！',
    tags: ['日常', '心情', '校园'],
    visibility: 0
  })
});

// 2. 获取我的动态列表
const myMoments = await fetch('http://106.52.165.122:8080/api/moments/my?page=1&pageSize=10', {
  headers: {
    'Authorization': 'Bearer ' + token
  }
});

// 3. 删除指定动态
await fetch('http://106.52.165.122:8080/api/moments/123', {
  method: 'DELETE',
  headers: {
    'Authorization': 'Bearer ' + token
  }
});
```

---

## 📋 功能完成状态

| 功能模块 | 状态 | 接口数量 |
|---------|------|----------|
| ✅ 用户注册登录 | 完成 | 2个 |
| ✅ 🆕 忘记密码（手机验证码） | **新增完成** | 2个 |
| ✅ 📸 上传头像 | 完成 | 1个 |
| ✅ 用户资料管理 | 完成 | 3个 |
| ✅ 📝 动态内容管理 | 完成 | 6个 |
| ✅ 🏷️ 多标签支持 | 完成 | 集成在动态中 |
| ✅ 📋 我的内容列表 | 完成 | 1个 |
| ✅ 🗑️ 删除我的内容 | 完成 | 1个 |
| ✅ 搜索功能 | 完成 | 4个 |
| **总计** | **20个接口** | |

---

## 🚀 部署说明

### 1. 环境要求
- Go 1.25.3+
- MySQL 8.0+
- Redis (可选，用于验证码存储)

### 2. 部署步骤
```bash
# 1. 拉取代码
git clone https://github.com/Yw332/campus-moments-go.git
cd campus-moments-go

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库等信息

# 3. 安装依赖
go mod tidy

# 4. 启动服务
go run cmd/api/main.go
```

### 3. 生产环境部署
```bash
# 编译
go build -o campus-moments cmd/api/main.go

# 运行
./campus-moments
```

---

## 📞 技术支持

如有问题，请联系开发团队或在GitHub Issues中反馈。

**所有功能已完整实现，可直接用于生产环境对接！** 🎉