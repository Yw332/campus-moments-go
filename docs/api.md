# 前端接口需求文档

> 📖 **文档说明**: 本文档详细列出了前端项目所需的所有API接口，包含接口状态、请求参数、返回格式等信息。
> 
> 📚 **相关文档**:
> - [后端实现逻辑文档](./backend.md) - 了解后端实现细节
> - [接口测试文档](./test.md) - 查看测试示例
> - [自动化测试脚本](./test.sh) - 一键测试所有接口

## 基础信息

- **后端地址**: `http://106.52.165.122:8080`
- **认证方式**: JWT Token (Bearer Token)
- **请求头**: `Authorization: Bearer <token>`
- **响应格式**: 
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

---

## 一、用户认证接口

### 1.1 用户登录
- **接口**: `POST /auth/login`
- **状态**: ✅ 已实现
- **使用位置**: `pages/login/login.vue`
- **请求参数**:
```json
{
  "account": "手机号或用户名",
  "password": "密码"
}
```
- **完整响应**:
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

### 1.2 用户注册
- **接口**: `POST /auth/register`
- **状态**: ✅ 已实现
- **使用位置**: `pages/register/register.vue`
- **请求参数**:
```json
{
  "username": "用户名",
  "phone": "手机号",
  "password": "密码"
}
```
- **完整响应**:
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

### 1.3 获取用户资料
- **接口**: `GET /api/users/profile`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: 未直接调用，但已在api.js中定义
- **完整响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "0000000001",
    "username": "testuser",
    "phone": "13800138000",
    "avatarUrl": "",
    "avatarType": 0,
    "avatarUpdatedAt": null,
    "role": 0,
    "status": 0,
    "createdAt": "2024-12-31T10:00:00Z"
  }
}
```

### 1.4 修改密码
- **接口**: `PUT /api/users/password`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/modify-password/modify-password.vue`
- **请求参数**:
```json
{
  "oldPassword": "旧密码",
  "newPassword": "新密码"
}
```
- **完整响应**:
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

### 1.5 退出登录
- **接口**: `POST /api/auth/logout`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/my/my.vue`
- **完整响应**:
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

---

## 二、动态相关接口

### 2.1 获取动态列表
- **接口**: `GET /api/moments`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/home/home.vue` (目前使用mockData，需要对接)
- **请求参数**:
  - `page`: 页码（默认1）
  - `pageSize`: 每页数量（默认10）
- **完整响应**:
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
        "createTime": "2024-01-15 10:30"
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

### 2.2 发布动态
- **接口**: `POST /api/moments`
- **状态**: ✅ 已实现（已优化支持title和images）
- **认证**: 需要Token
- **使用位置**: `pages/issue/issue.vue` (有TODO注释，需要对接)
- **请求参数**:
```json
{
  "title": "动态标题",
  "content": "动态内容",
  "tags": ["标签1", "标签2"],
  "images": ["图片URL1", "图片URL2"]
}
```
- **完整响应**:
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "id": 1,
    "userId": "0000000001",
    "title": "动态标题",
    "content": "动态内容",
    "images": ["图片URL1", "图片URL2"],
    "tags": ["标签1", "标签2"],
    "likeCount": 0,
    "commentCount": 0,
    "createdAt": "2024-12-31T16:30:00Z"
  }
}
```

### 2.3 获取动态详情
- **接口**: `GET /api/moments/:id`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/detail/detail.vue` (目前使用mockData)
- **完整响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "userId": "0000000001",
    "title": "动态标题",
    "content": "动态内容",
    "images": ["http://localhost:8080/static/files/image1.jpg"],
    "tags": ["标签1"],
    "likeCount": 12,
    "commentCount": 5,
    "viewCount": 100,
    "author": {
      "userId": "0000000001",
      "username": "testuser",
      "avatarUrl": ""
    },
    "createdAt": "2024-12-31T10:30:00Z"
  }
}
```

### 2.4 获取我的动态
- **接口**: `GET /api/moments/my`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/my/my.vue` (目前使用mockData)
- **请求参数**:
  - `page`: 页码
  - `pageSize`: 每页数量
- **完整响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "userId": "0000000001",
        "title": "我的动态",
        "content": "内容",
        "likeCount": 5,
        "commentCount": 2,
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

### 2.5 点赞动态
- **接口**: `POST /api/likes/post/:postId`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/detail/detail.vue` (需要对接)
- **完整响应**:
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "liked": true
  }
}
```

---

## 三、评论相关接口

### 3.1 获取评论列表
- **接口**: `GET /public/posts/:id/comments`
- **状态**: ✅ 已实现
- **认证**: 不需要Token（公开接口）
- **使用位置**: `pages/detail/detail.vue` (目前使用mockData)
- **请求参数**:
  - `page`: 页码（默认1）
  - `pageSize`: 每页数量（默认20）
- **完整响应**:
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

### 3.2 发布评论
- **接口**: `POST /api/comments/post/:postId`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/detail/detail.vue` (需要对接)
- **请求参数**:
```json
{
  "content": "评论内容"
}
```
- **完整响应**:
```json
{
  "code": 200,
  "message": "评论成功",
  "data": {
    "id": 1,
    "postId": 1,
    "userId": "0000000001",
    "content": "评论内容",
    "likeCount": 0,
    "createdAt": "2024-12-31T16:30:00Z"
  }
}
```

### 3.3 点赞评论
- **接口**: `POST /api/comments/:id/like`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/detail/detail.vue` (需要对接)
- **完整响应**:
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "liked": true
  }
}
```

---

## 四、搜索相关接口

### 4.1 搜索内容
- **接口**: `GET /api/search`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/searchresult/searchresult.vue` (目前使用mockData)
- **请求参数**:
  - `keyword`: 搜索关键词
  - `page`: 页码
  - `pageSize`: 每页数量
- **完整响应**:
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

### 4.2 获取热门关键词
- **接口**: `GET /api/search/hot-words`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/search/search.vue` (目前使用mockData)
- **完整响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    "校园",
    "活动",
    "学习",
    "美食",
    "运动"
  ]
}
```

---

## 五、文件上传接口

### 5.1 通用文件上传
- **接口**: `POST /api/upload/file`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/issue/issue.vue` (发布页面需要上传图片)
- **请求格式**: `multipart/form-data`
- **参数**:
  - `file`: 文件
- **完整响应**:
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

### 5.2 头像上传
- **接口**: `POST /api/upload/avatar`
- **状态**: ✅ 已实现
- **认证**: 需要Token
- **使用位置**: `pages/register/register.vue` (注册页面可选)
- **请求格式**: `multipart/form-data`
- **参数**:
  - `avatar`: 头像文件
- **完整响应**:
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

---

## 六、标签相关接口

### 6.1 获取标签列表
- **接口**: `GET /public/tags`
- **状态**: ✅ 已实现
- **认证**: 不需要Token
- **使用位置**: `pages/issue/issue.vue` (标签选择功能)
- **请求参数**:
  - `page`: 页码（默认1）
  - `pageSize`: 每页数量（默认50）
- **完整响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "tags": [
      {
        "id": 1,
        "name": "校园",
        "color": "#DDA0DD",
        "icon": "",
        "description": "",
        "usageCount": 50,
        "lastUsedAt": "2024-12-31T10:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "pageSize": 50
  }
}
```

---

## 接口完成情况总结

### ✅ 已完全实现的接口（18个）

#### 用户认证模块（5个）
1. ✅ POST /auth/login - 用户登录
2. ✅ POST /auth/register - 用户注册
3. ✅ GET /api/users/profile - 获取用户资料
4. ✅ PUT /api/users/password - 修改密码
5. ✅ POST /api/auth/logout - 退出登录

#### 动态模块（4个）
6. ✅ GET /api/moments - 获取动态列表
7. ✅ POST /api/moments - 发布动态
8. ✅ GET /api/moments/:id - 获取动态详情
9. ✅ GET /api/moments/my - 获取我的动态

#### 评论模块（3个）
10. ✅ GET /public/posts/:id/comments - 获取评论列表
11. ✅ POST /api/comments/post/:postId - 发布评论
12. ✅ POST /api/comments/:id/like - 点赞评论

#### 点赞模块（1个）
13. ✅ POST /api/likes/post/:postId - 点赞动态

#### 搜索模块（2个）
14. ✅ GET /api/search - 搜索内容
15. ✅ GET /api/search/hot-words - 获取热门关键词

#### 文件上传模块（2个）
16. ✅ POST /api/upload/file - 通用文件上传
17. ✅ POST /api/upload/avatar - 头像上传

#### 标签模块（1个）
18. ✅ GET /public/tags - 获取标签列表

### 📝 接口实现状态

**所有前端需要的接口已全部实现！**

- ✅ 核心功能接口：100%完成
- ✅ 交互功能接口：100%完成
- ✅ 辅助功能接口：100%完成

---

## 前端需要对接的接口

### 高优先级（核心功能）
1. **GET /api/moments** - 首页动态列表（目前使用mockData）
2. **POST /api/moments** - 发布动态（有TODO注释）
3. **GET /api/moments/:id** - 动态详情（目前使用mockData）
4. **POST /api/upload/file** - 图片上传（发布页面需要）

### 中优先级（交互功能）
5. **GET /api/search** - 搜索功能（搜索结果页）
6. **GET /api/search/hot-words** - 热门关键词（搜索页）
7. **POST /api/likes/post/:postId** - 点赞动态
8. **GET /public/posts/:id/comments** - 获取评论列表
9. **POST /api/comments/post/:postId** - 发布评论

### 低优先级（辅助功能）
10. **GET /api/moments/my** - 我的动态列表
11. **GET /public/tags** - 标签列表

---

## 注意事项

1. 所有 `/api/*` 开头的接口都需要在请求头中携带Token
2. 文件上传接口需要使用 `multipart/form-data` 格式
3. 前端需要处理401错误，自动清除token并跳转登录页
4. 响应格式统一为 `{code, message, data}` 结构

