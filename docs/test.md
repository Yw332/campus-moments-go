# 后端接口测试文档

## 测试说明

本文档模拟前端请求方式，测试后端接口的返回结果。所有测试均基于前端实际调用场景。

**测试环境**：
- 后端地址：`http://106.52.165.122:8080` 或 `http://localhost:8080`
- 请求格式：JSON
- 认证方式：JWT Token (Bearer Token)

**响应格式**：
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

---

## 一、用户认证接口测试

### 1.1 用户注册

**功能说明**：新用户注册账号，支持用户名、手机号、密码注册。

**请求示例**：
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "phone": "13800138000",
    "password": "Test123456"
  }'
```

**预期返回**：
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "userId": "0000000001",
    "username": "testuser",
    "phone": "13800138000"
  }
}
```

**测试要点**：
- ✅ 用户名格式验证（3-20字符）
- ✅ 手机号格式验证（11位，1开头）
- ✅ 密码强度验证（8-20位，包含大小写字母和数字）
- ✅ 用户名和手机号唯一性检查

---

### 1.2 用户登录

**功能说明**：用户登录获取Token，支持用户名或手机号登录。

**请求示例**：
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "account": "testuser",
    "password": "Test123456"
  }'
```

**预期返回**：
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "userInfo": {
      "userId": "0000000001",
      "username": "testuser",
      "phone": "13800138000",
      "role": 0,
      "isAdmin": false
    }
  }
}
```

**测试要点**：
- ✅ 支持用户名登录
- ✅ 支持手机号登录
- ✅ 密码验证
- ✅ 返回JWT Token
- ✅ 返回用户基本信息

**保存Token**：后续测试需要使用此Token
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### 1.3 获取用户资料

**功能说明**：获取当前登录用户的详细资料信息。

**请求示例**：
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "0000000001",
    "username": "testuser",
    "phone": "13800138000",
    "avatarUrl": "",
    "role": 0,
    "status": 0,
    "createdAt": "2024-12-31T10:00:00Z"
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 返回用户完整信息
- ✅ 不包含敏感信息（密码）

---

### 1.4 修改密码

**功能说明**：用户修改登录密码，需要验证旧密码。

**请求示例**：
```bash
curl -X PUT http://localhost:8080/api/users/password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "oldPassword": "Test123456",
    "newPassword": "NewPass123"
  }'
```

**预期返回**：
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 验证旧密码正确性
- ✅ 新密码强度验证
- ✅ 密码加密存储

---

### 1.5 退出登录

**功能说明**：用户退出登录，将Token加入黑名单。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "退出成功",
  "data": {
    "userId": "0000000001",
    "logoutAt": "2024-12-31 16:30:00"
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ Token加入黑名单
- ✅ 退出后Token失效

---

## 二、动态相关接口测试

### 2.1 获取动态列表

**功能说明**：获取首页动态列表，支持分页，返回格式匹配前端瀑布流需求。

**请求示例**：
```bash
curl -X GET "http://localhost:8080/api/moments?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "动态标题",
        "author": "用户名",
        "imageUrl": "http://localhost:8080/static/files/image1.jpg",
        "likeCount": 12,
        "createTime": "2024-12-31 10:30"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 100
    }
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 分页功能正常
- ✅ 返回格式匹配前端（id, title, author, imageUrl, likeCount, createTime）
- ✅ 按创建时间倒序排列
- ✅ 提取第一张图片作为imageUrl

---

### 2.2 发布动态

**功能说明**：用户发布新动态，支持标题、内容、标签、图片。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/moments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "我的第一条动态",
    "content": "这是动态内容，可以包含文字描述",
    "tags": ["校园", "生活"],
    "images": [
      "http://localhost:8080/static/files/image1.jpg",
      "http://localhost:8080/static/files/image2.jpg"
    ]
  }'
```

**预期返回**：
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "id": 1,
    "userId": "0000000001",
    "title": "我的第一条动态",
    "content": "这是动态内容，可以包含文字描述",
    "images": ["http://localhost:8080/static/files/image1.jpg"],
    "tags": ["校园", "生活"],
    "likeCount": 0,
    "commentCount": 0,
    "createdAt": "2024-12-31T16:30:00Z"
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 支持title参数
- ✅ 支持images数组
- ✅ 支持tags数组
- ✅ 自动关联当前用户

---

### 2.3 获取动态详情

**功能说明**：获取单条动态的详细信息，包含图片、评论等。

**请求示例**：
```bash
curl -X GET http://localhost:8080/api/moments/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "动态标题",
    "content": "动态内容",
    "images": ["http://localhost:8080/static/files/image1.jpg"],
    "tags": ["标签1"],
    "likeCount": 12,
    "commentCount": 5,
    "author": {
      "userId": "0000000001",
      "username": "testuser",
      "avatarUrl": ""
    },
    "createdAt": "2024-12-31T10:30:00Z"
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 返回完整动态信息
- ✅ 包含作者信息（author字段）
- ✅ 包含图片数组

---

### 2.4 获取我的动态

**功能说明**：获取当前用户发布的所有动态列表。

**请求示例**：
```bash
curl -X GET "http://localhost:8080/api/moments/my?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "我的动态",
        "content": "内容",
        "likeCount": 5,
        "createdAt": "2024-12-31T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 1
    }
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 只返回当前用户的动态
- ✅ 支持分页

---

## 三、评论相关接口测试

### 3.1 获取评论列表

**功能说明**：获取指定动态的所有评论，公开接口无需Token。

**请求示例**：
```bash
curl -X GET "http://localhost:8080/public/posts/1/comments?page=1&pageSize=20" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "comments": [
      {
        "id": 1,
        "postId": 1,
        "userId": "0000000001",
        "content": "评论内容",
        "likeCount": 0,
        "user": {
          "username": "testuser",
          "avatarUrl": ""
        },
        "createdAt": "2024-12-31T10:30:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  }
}
```

**测试要点**：
- ✅ 公开接口，无需Token
- ✅ 返回评论列表
- ✅ 包含评论者信息
- ✅ 支持分页

---

### 3.2 发布评论

**功能说明**：在指定动态下发布评论。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/comments/post/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这是一条评论"
  }'
```

**预期返回**：
```json
{
  "code": 200,
  "message": "评论成功",
  "data": {
    "id": 1,
    "postId": 1,
    "userId": "0000000001",
    "content": "这是一条评论",
    "likeCount": 0,
    "createdAt": "2024-12-31T16:30:00Z"
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 自动关联当前用户
- ✅ 自动更新动态评论数

---

### 3.3 点赞评论

**功能说明**：点赞或取消点赞评论。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/comments/1/like \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "liked": true,
    "likeCount": 1
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 支持点赞/取消点赞切换
- ✅ 自动更新点赞数

---

## 四、点赞相关接口测试

### 4.1 点赞动态

**功能说明**：点赞或取消点赞动态。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/likes/post/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "liked": true
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 支持点赞/取消点赞切换
- ✅ 自动更新动态点赞数
- ✅ 防止重复点赞

---

## 五、搜索相关接口测试

### 5.1 搜索内容

**功能说明**：搜索动态和用户，支持关键词搜索。

**请求示例**：
```bash
curl -X GET "http://localhost:8080/api/search?keyword=测试&page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "posts": [
      {
        "postId": 1,
        "title": "测试动态",
        "content": "包含测试关键词的内容",
        "author": {
          "userId": "0000000001",
          "username": "testuser"
        },
        "likeCount": 5
      }
    ],
    "users": [],
    "pagination": {
      "page": 1,
      "pageSize": 10,
      "total": 1
    }
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 搜索动态标题和内容
- ✅ 搜索用户名
- ✅ 返回格式匹配前端（postId, author）

---

### 5.2 获取热门关键词

**功能说明**：获取搜索热门关键词列表。

**请求示例**：
```bash
curl -X GET http://localhost:8080/api/search/hot-words \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "keyword": "校园",
      "count": 50
    },
    {
      "keyword": "生活",
      "count": 30
    }
  ]
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 按搜索次数排序
- ✅ 返回热门关键词

---

## 六、文件上传接口测试

### 6.1 通用文件上传

**功能说明**：上传图片文件，用于动态发布。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/upload/file \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/image.jpg"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "fileId": "uuid-string",
    "filename": "20241231163000_uuid.jpg",
    "originalName": "image.jpg",
    "fileSize": 102400,
    "fileType": ".jpg",
    "fileUrl": "http://localhost:8080/static/files/20241231163000_uuid.jpg"
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 支持图片格式（jpg, jpeg, png, gif, webp）
- ✅ 文件大小限制（10MB）
- ✅ 返回文件访问URL

---

### 6.2 头像上传

**功能说明**：上传用户头像。

**请求示例**：
```bash
curl -X POST http://localhost:8080/api/upload/avatar \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@/path/to/avatar.jpg"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatarUrl": "http://localhost:8080/static/avatars/20241231163000_uuid.jpg",
    "filename": "20241231163000_uuid.jpg",
    "size": 51200
  }
}
```

**测试要点**：
- ✅ 需要Token认证
- ✅ 只支持图片格式
- ✅ 文件大小限制（5MB）
- ✅ 自动更新用户头像URL

---

## 七、标签相关接口测试

### 7.1 获取标签列表

**功能说明**：获取所有可用标签列表，用于发布动态时选择。

**请求示例**：
```bash
curl -X GET http://localhost:8080/public/tags \
  -H "Content-Type: application/json"
```

**预期返回**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "name": "校园",
      "postCount": 50
    },
    {
      "id": 2,
      "name": "生活",
      "postCount": 30
    }
  ]
}
```

**测试要点**：
- ✅ 公开接口，无需Token
- ✅ 返回所有标签
- ✅ 包含标签使用次数

---

## 八、完整测试流程

### 测试步骤

1. **注册新用户**
   ```bash
   curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","phone":"13800138000","password":"Test123456"}'
   ```

2. **登录获取Token**
   ```bash
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{"account":"testuser","password":"Test123456"}'
   ```
   保存返回的Token

3. **上传图片**
   ```bash
   curl -X POST http://localhost:8080/api/upload/file \
     -H "Authorization: Bearer $TOKEN" \
     -F "file=@image.jpg"
   ```
   保存返回的fileUrl

4. **发布动态**
   ```bash
   curl -X POST http://localhost:8080/api/moments \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "title":"测试动态",
       "content":"这是测试内容",
       "tags":["测试"],
       "images":["http://localhost:8080/static/files/image.jpg"]
     }'
   ```

5. **获取动态列表**
   ```bash
   curl -X GET "http://localhost:8080/api/moments?page=1&pageSize=10" \
     -H "Authorization: Bearer $TOKEN"
   ```

6. **发布评论**
   ```bash
   curl -X POST http://localhost:8080/api/comments/post/1 \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"content":"这是一条评论"}'
   ```

7. **点赞动态**
   ```bash
   curl -X POST http://localhost:8080/api/likes/post/1 \
     -H "Authorization: Bearer $TOKEN"
   ```

8. **搜索内容**
   ```bash
   curl -X GET "http://localhost:8080/api/search?keyword=测试" \
     -H "Authorization: Bearer $TOKEN"
   ```

---

## 九、错误处理测试

### 9.1 未认证请求

**测试**：不带Token访问需要认证的接口
```bash
curl -X GET http://localhost:8080/api/users/profile
```

**预期返回**：
```json
{
  "code": 401,
  "message": "未认证",
  "data": null
}
```

### 9.2 无效Token

**测试**：使用无效Token
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer invalid_token"
```

**预期返回**：
```json
{
  "code": 401,
  "message": "Token无效或已过期",
  "data": null
}
```

### 9.3 参数验证错误

**测试**：注册时使用无效参数
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"ab","phone":"123","password":"123"}'
```

**预期返回**：
```json
{
  "code": 400,
  "message": "用户名长度必须在3-20个字符之间",
  "data": null
}
```

---

## 十、性能测试

### 10.1 响应时间

- 登录接口：< 500ms
- 获取动态列表：< 300ms
- 搜索接口：< 500ms
- 文件上传：< 2s（取决于文件大小）

### 10.2 并发测试

使用工具测试并发请求：
```bash
# 使用ab工具测试
ab -n 100 -c 10 -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/moments?page=1&pageSize=10
```

---

## 测试检查清单

### 核心功能
- [ ] 用户注册
- [ ] 用户登录（获取Token）
- [ ] 获取用户资料
- [ ] 修改密码
- [ ] 退出登录

### 动态功能
- [ ] 获取动态列表（分页）
- [ ] 发布动态（带图片）
- [ ] 获取动态详情
- [ ] 获取我的动态
- [ ] 点赞动态

### 评论功能
- [ ] 获取评论列表
- [ ] 发布评论
- [ ] 点赞评论

### 搜索功能
- [ ] 搜索内容
- [ ] 获取热门关键词

### 文件上传
- [ ] 上传图片
- [ ] 上传头像

### 标签功能
- [ ] 获取标签列表

### 错误处理
- [ ] 未认证请求
- [ ] 无效Token
- [ ] 参数验证错误

---

## 注意事项

1. **Token管理**：所有需要认证的接口都需要在请求头中携带Token
2. **文件上传**：使用 `multipart/form-data` 格式，字段名为 `file` 或 `avatar`
3. **响应格式**：统一为 `{code, message, data}` 结构
4. **错误处理**：code不为200时表示请求失败，message包含错误信息
5. **分页参数**：page从1开始，pageSize建议10-20

---

## 快速测试脚本

### Linux/Mac
```bash
# 设置变量
API_BASE="http://localhost:8080"

# 注册
curl -X POST $API_BASE/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","phone":"13800138000","password":"Test123456"}'

# 登录并保存Token
TOKEN=$(curl -s -X POST $API_BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"testuser","password":"Test123456"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Token: $TOKEN"

# 获取动态列表
curl -X GET "$API_BASE/api/moments?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

### Windows PowerShell
```powershell
$baseUrl = "http://localhost:8080"

# 注册
$registerBody = @{
    username = "testuser"
    phone = "13800138000"
    password = "Test123456"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method POST -Body $registerBody -ContentType "application/json"

# 登录
$loginBody = @{
    account = "testuser"
    password = "Test123456"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
$token = $response.data.token

# 获取动态列表
$headers = @{
    "Authorization" = "Bearer $token"
}
Invoke-RestMethod -Uri "$baseUrl/api/moments?page=1&pageSize=10" -Method GET -Headers $headers
```

