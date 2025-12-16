# 用户认证接口文档

## 基本信息

- **Base URL**: `http://localhost:8081`
- **Content-Type**: `application/json`
- **认证方式**: JWT Token (需要在请求头中携带 `Authorization: Bearer <token>`)

---

## 1. 用户注册

**接口地址**: `POST /auth/register`

**请求参数**:

```json
{
  "username": "Yw166332",
  "phone": "17875242005", 
  "password": "JiangCan030"
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，3-20个字符，支持字母、数字、中文、下划线 |
| phone | string | 是 | 手机号，11位数字，1开头 |
| password | string | 是 | 密码，8-20位，必须包含大小写字母和数字 |

**成功响应**:

```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "userId": 4,
    "username": "Yw166332",
    "phone": "17875242005"
  }
}
```

**错误响应**:

```json
{
  "code": 400,
  "message": "密码必须包含大小写字母和数字",
  "data": null
}
```

**可能的错误信息**:
- `用户名长度必须在3-20个字符之间`
- `用户名只能包含字母、数字、中文和下划线`
- `手机号格式不正确`
- `密码长度必须在8-20位之间`
- `密码必须包含大小写字母和数字`
- `用户名已存在`
- `手机号已被注册`

---

## 2. 用户登录

**接口地址**: `POST /auth/login`

**请求参数**:

```json
{
  "account": "Yw166332",
  "password": "JiangCan030"
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| account | string | 是 | 用户名或手机号 |
| password | string | 是 | 用户密码 |

**成功响应**:

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "userInfo": {
      "userId": 4,
      "username": "Yw166332",
      "phone": "17875242005"
    }
  }
}
```

**错误响应**:

```json
{
  "code": 404,
  "message": "账户不存在",
  "data": null
}
```

**可能的错误信息**:
- `账户不存在` (404)
- `密码错误` (401)
- `账户已被锁定，请联系管理员` (403)

---

## 3. 获取用户资料

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
    "userId": 4,
    "username": "Yw166332",
    "phone": "17875242005"
  }
}
```

**错误响应**:

```json
{
  "code": 401,
  "message": "未认证",
  "data": null
}
```

---

## 4. 修改密码

**接口地址**: `PUT /api/users/password`

**请求头**:

```
Authorization: Bearer <your_token>
```

**请求参数**:

```json
{
  "oldPassword": "JiangCan030",
  "newPassword": "NewPassword123"
}
```

**参数说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| oldPassword | string | 是 | 原密码 |
| newPassword | string | 是 | 新密码，需符合密码强度要求 |

**成功响应**:

```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

**错误响应**:

```json
{
  "code": 400,
  "message": "原密码错误",
  "data": null
}
```

---

## 5. 退出登录

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
  "data": 4
}
```

---

## 通用响应格式

所有接口都遵循统一的响应格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

**状态码说明**:

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 测试账号

已注册的测试账号：

| 用户名 | 手机号 | 密码 |
|--------|--------|------|
| Yw166332 | 17875242005 | JiangCan030 |

---

## 注意事项

1. **Token有效期**: 7天 (168小时)
2. **密码强度**: 必须包含大小写字母和数字，长度8-20位
3. **用户名规则**: 3-20个字符，支持字母、数字、中文、下划线
4. **手机号格式**: 中国大陆11位手机号
5. **认证方式**: JWT Token，需要在请求头中携带
6. **接口前缀**: 
   - 无需认证的接口: `/auth/*`
   - 需要认证的接口: `/api/*`

---

## 前端集成建议

### 1. 登录流程
```javascript
// 登录
const login = async (account, password) => {
  const response = await fetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ account, password })
  });
  const data = await response.json();
  
  if (data.code === 200) {
    // 保存token到本地存储
    localStorage.setItem('token', data.data.token);
    localStorage.setItem('userInfo', JSON.stringify(data.data.userInfo));
  }
  return data;
};

// 请求拦截器（添加认证头）
const authFetch = (url, options = {}) => {
  const token = localStorage.getItem('token');
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  };
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  return fetch(url, { ...options, headers });
};
```

### 2. 错误处理
```javascript
// 统一错误处理
const handleResponse = (response) => {
  if (response.code === 401) {
    // token过期或无效，跳转到登录页
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    window.location.href = '/login';
  }
  return response;
};
```

---

## 更新日志

- **v1.0.0** (2024-12-13): 初始版本，包含基础认证功能
  - 用户注册/登录/退出
  - 密码修改
  - 用户资料获取
  - JWT Token认证