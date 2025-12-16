# 校园动态系统 API 文档 (稳定版本)

## 基本信息

- **Base URL**: `http://106.52.165.122:8080`
- **Content-Type**: `application/json`
- **认证方式**: JWT Token (需要在请求头中携带 `Authorization: Bearer <token>`)

---

## 1. 用户认证相关

### 1.1 用户注册
**接口地址**: `POST /auth/register`

**请求参数**:
```json
{
  "username": "testuser",
  "phone": "13800138000", 
  "password": "Test123456"
}
```

**成功响应**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "userId": 100000001,
    "username": "testuser",
    "phone": "13800138000"
  }
}
```

### 1.2 用户登录
**接口地址**: `POST /auth/login`

**请求参数**:
```json
{
  "account": "testuser",
  "password": "Test123456"
}
```

**成功响应**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "userInfo": {
      "userId": 100000001,
      "username": "testuser",
      "phone": "13800138000"
    }
  }
}
```

### 1.3 退出登录
**接口地址**: `POST /api/auth/logout`

**请求头**:
```
Authorization: Bearer <your_token>
```

**成功响应**:
```json
{
  "code": 200,
  "message": "退出成功",
  "data": null
}
```

---

## 2. 用户资料相关

### 2.1 获取用户资料
**接口地址**: `GET /api/users/profile`

**请求头**:
```
Authorization: Bearer <your_token>
```

**成功响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "userId": 100000001,
    "username": "testuser",
    "phone": "13800138000",
    "avatar": "",
    "postCount": 0,
    "likeCount": 0,
    "commentCount": 0
  }
}
```

### 2.2 更新用户资料
**接口地址**: `PUT /api/users/profile`

**请求头**:
```
Authorization: Bearer <your_token>
```

**请求参数**:
```json
{
  "username": "newusername",
  "phone": "13900139000"
}
```

### 2.3 修改密码
**接口地址**: `PUT /api/users/password`

**请求头**:
```
Authorization: Bearer <your_token>
```

**请求参数**:
```json
{
  "oldPassword": "Test123456",
  "newPassword": "NewTest123456"
}
```

---

## 3. 文件上传相关

### 3.1 上传文件
**接口地址**: `POST /api/upload/file`

**请求头**:
```
Authorization: Bearer <your_token>
Content-Type: multipart/form-data
```

**请求参数**:
- file: 文件 (FormData)

**成功响应**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "/uploads/20241216/example.jpg",
    "filename": "example.jpg",
    "size": 1024000
  }
}
```

### 3.2 上传头像
**接口地址**: `POST /api/upload/avatar`

**请求头**:
```
Authorization: Bearer <your_token>
Content-Type: multipart/form-data
```

**请求参数**:
- avatar: 头像文件 (FormData)

**成功响应**:
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatarUrl": "/uploads/avatars/example.jpg",
    "filename": "example.jpg"
  }
}
```

---

## 4. 系统接口

### 4.1 健康检查
**接口地址**: `GET /health`

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok",
    "service": "Campus Moments Go API",
    "database": "connected"
  }
}
```

### 4.2 根路径
**接口地址**: `GET /`

**成功响应**:
```json
{
  "code": 0,
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

## 5. 通用响应格式

所有接口都遵循统一的响应格式：

```json
{
  "code": 200,        // 状态码：200成功，400客户端错误，401未认证，404不存在，500服务器错误
  "message": "操作成功", // 消息说明
  "data": {}          // 响应数据
}
```

---

## 6. 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 7. 测试建议

### 7.1 在 Apifox 中回退接口的方法：

1. **删除现有的 API 集合**：在 Apifox 中删除当前的校园动态系统 API 集合
2. **重新导入**：使用上面的文档重新创建接口
3. **修改 Base URL**：确保所有接口的 Base URL 设置为 `http://localhost:8080`

### 7.2 快速测试步骤：

1. 先测试 `/health` 接口确认服务正常
2. 测试用户注册接口 `/auth/register`
3. 测试用户登录接口 `/auth/login` 获取 token
4. 使用获取的 token 测试需要认证的接口

---

## 注意事项

1. **端口配置**：确保服务器运行在 8080 端口
2. **Token 使用**：需要认证的接口必须在请求头中携带有效的 JWT Token
3. **文件上传**：使用 multipart/form-data 格式
4. **错误处理**：前端需要统一处理不同的错误码