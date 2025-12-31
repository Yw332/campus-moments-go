# Campus Moments Go API 完整接口文档

## 📋 基本信息

| 项目 | 说明 |
|------|------|
| **Base URL** | `http://106.52.165.122:8080` |
| **Content-Type** | `application/json` |
| **认证方式** | JWT Token (需要在请求头中携带 `Authorization: Bearer <token>`) |
| **响应格式** | JSON |

---

## 🔑 认证方式

所有 `/api/*` 开头的接口都需要在请求头中携带JWT Token：

```
Authorization: Bearer <your_token_here>
```

---

## 📡 统一响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误描述",
  "data": null
}
```

### 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 🗂️ 接口分类

### 1. 系统接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/` | 首页 | ❌ |
| GET | `/health` | 健康检查 | ❌ |

**健康检查响应示例**：
```json
{
  "status": "ok",
  "message": "Campus Moments API is running"
}
```

---

### 2. 用户认证接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/auth/register` | 用户注册 | ❌ |
| POST | `/auth/login` | 用户登录 | ❌ |
| POST | `/auth/send-verification` | 发送验证码 | ❌ |
| POST | `/auth/verify-and-reset` | 验证并重置密码 | ❌ |
| POST | `/api/auth/logout` | 用户登出 | ✅ |

#### 2.1 用户注册

**请求参数**：
```json
{
  "username": "Yw166332",
  "phone": "17875242005",
  "password": "JiangCan030"
}
```

**参数说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，3-20个字符，支持字母、数字、中文、下划线 |
| phone | string | 是 | 手机号，11位数字，1开头 |
| password | string | 是 | 密码，8-20位，必须包含大小写字母和数字 |

**成功响应**：
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

#### 2.2 用户登录

**请求参数**：
```json
{
  "account": "Yw166332",
  "password": "JiangCan030"
}
```

**参数说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| account | string | 是 | 用户名或手机号 |
| password | string | 是 | 用户密码 |

**成功响应**：
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

#### 2.3 发送验证码

**请求参数**：
```json
{
  "phone": "17875242005"
}
```

#### 2.4 验证并重置密码

**请求参数**：
```json
{
  "phone": "17875242005",
  "verificationCode": "123456",
  "newPassword": "NewPassword123"
}
```

#### 2.5 用户登出

**成功响应**：
```json
{
  "code": 200,
  "message": "退出成功",
  "data": 4
}
```

---

### 3. 用户信息接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/users/profile` | 获取当前用户资料 | ✅ |
| PUT | `/api/users/profile` | 更新用户资料 | ✅ |
| PUT | `/api/users/password` | 修改密码 | ✅ |
| PUT | `/api/users/avatar` | 更新头像 | ✅ |
| PUT | `/api/users/signature` | 更新个性签名 | ✅ |
| POST | `/api/users/active` | 更新最后活跃时间 | ✅ |
| GET | `/api/users/:userId` | 获取指定用户信息 | ✅ |
| GET | `/api/users/search` | 搜索用户 | ✅ |

#### 3.1 获取用户资料

**请求头**：
```
Authorization: Bearer <your_token>
```

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "userId": 4,
    "username": "Yw166332",
    "phone": "17875242005",
    "avatarUrl": "",
    "signature": "",
    "wechatNickname": ""
  }
}
```

#### 3.2 更新用户资料

**请求参数**：
```json
{
  "username": "新用户名",
  "phone": "13800138000",
  "avatarUrl": "头像URL",
  "signature": "个性签名",
  "wechatNickname": "微信昵称"
}
```

#### 3.3 修改密码

**请求参数**：
```json
{
  "oldPassword": "JiangCan030",
  "newPassword": "NewPassword123"
}
```

#### 3.4 更新头像

**请求参数**：
```json
{
  "avatarUrl": "头像URL"
}
```

#### 3.5 更新个性签名

**请求参数**：
```json
{
  "signature": "新的个性签名"
}
```

#### 3.6 更新最后活跃时间

**请求参数**：
```json
{
  "lastActive": "2024-12-30T10:00:00Z"
}
```

#### 3.7 获取指定用户信息

**路径参数**：
- `userId`: 用户ID

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "userId": 4,
    "username": "Yw166332",
    "avatarUrl": "头像URL",
    "signature": "个性签名"
  }
}
```

#### 3.8 搜索用户

**查询参数**：
- `keyword`: 搜索关键词

**成功响应**：
```json
{
  "code": 200,
  "message": "搜索成功",
  "data": [
    {
      "userId": 4,
      "username": "Yw166332",
      "avatarUrl": "头像URL"
    }
  ]
}
```

---

### 4. 帖子接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/public/posts` | 获取公开帖子列表 | ❌ |
| GET | `/public/posts/:id` | 获取帖子详情（公开） | ❌ |
| GET | `/home` | 获取首页帖子（支持公开和好友） | ❌ |
| POST | `/api/posts` | 创建帖子 | ✅ |
| PUT | `/api/posts/:id` | 更新帖子 | ✅ |
| DELETE | `/api/posts/:id` | 删除帖子 | ✅ |
| GET | `/api/posts/my` | 获取我的帖子 | ✅ |
| GET | `/api/posts/user/:userId` | 获取用户帖子 | ✅ |

#### 4.1 创建帖子

**请求参数**：
```json
{
  "title": "标题",
  "content": "内容",
  "images": ["url1", "url2"],
  "video": "video_url",
  "visibility": 0,
  "tags": ["标签1", "标签2"]
}
```

**参数说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 否 | 帖子标题 |
| content | string | 是 | 帖子内容 |
| images | array | 否 | 图片URL数组 |
| video | string | 否 | 视频URL |
| visibility | int | 否 | 可见性：0-公开，1-好友可见，2-仅自己可见 |
| tags | array | 否 | 标签数组 |

**成功响应**：
```json
{
  "code": 201,
  "message": "发布成功",
  "data": {
    "postId": 1,
    "title": "标题",
    "content": "内容"
  }
}
```

#### 4.2 更新帖子

**路径参数**：
- `id`: 帖子ID

**请求参数**：
```json
{
  "title": "新标题",
  "content": "新内容",
  "images": ["url1"],
  "video": "",
  "visibility": 0,
  "tags": ["标签1"]
}
```

#### 4.3 删除帖子

**路径参数**：
- `id`: 帖子ID

**成功响应**：
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

#### 4.4 获取帖子列表

**查询参数**：
- `page`: 页码（可选，默认1）
- `pageSize`: 每页数量（可选，默认20）

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "posts": [
      {
        "postId": 1,
        "title": "标题",
        "content": "内容",
        "author": {
          "userId": 1,
          "username": "用户名",
          "avatarUrl": "头像URL"
        },
        "createdAt": "2024-12-30T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 5. 评论接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/public/posts/:postId/comments` | 获取评论列表 | ❌ |
| POST | `/api/comments/post/:postId` | 创建评论 | ✅ |
| PUT | `/api/comments/:id` | 更新评论 | ✅ |
| DELETE | `/api/comments/:id` | 删除评论 | ✅ |
| POST | `/api/comments/:id/like` | 点赞评论 | ✅ |
| POST | `/api/comments/:id/reply` | 回复评论 | ✅ |
| GET | `/api/comments/:id/likes` | 获取评论点赞列表 | ✅ |

#### 5.1 创建评论

**路径参数**：
- `postId`: 帖子ID

**请求参数**：
```json
{
  "content": "评论内容",
  "parentCommentId": 0
}
```

**参数说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 评论内容 |
| parentCommentId | int | 否 | 父评论ID（回复时使用） |

#### 5.2 更新评论

**路径参数**：
- `id`: 评论ID

**请求参数**：
```json
{
  "content": "新的评论内容"
}
```

#### 5.3 删除评论

**路径参数**：
- `id`: 评论ID

#### 5.4 回复评论

**路径参数**：
- `id`: 被回复的评论ID

**请求参数**：
```json
{
  "content": "回复内容"
}
```

---

### 6. 点赞接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/likes/post/:postId` | 点赞/取消点赞帖子 | ✅ |
| GET | `/api/likes/posts/:postId` | 获取帖子点赞列表 | ✅ |
| GET | `/api/likes/comments/:commentId` | 获取评论点赞列表 | ✅ |
| GET | `/api/likes/users/:userId` | 获取用户点赞列表 | ✅ |

#### 6.1 点赞/取消点赞帖子

**路径参数**：
- `postId`: 帖子ID

**成功响应**：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "liked": true,
    "likeCount": 10
  }
}
```

#### 6.2 获取点赞列表

**路径参数**：
- `postId` 或 `commentId`: 目标ID

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "userId": 1,
      "username": "用户名",
      "avatarUrl": "头像URL"
    }
  ]
}
```

---

### 7. 好友接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/friends/request` | 发送好友请求 | ✅ |
| GET | `/api/friends/requests` | 获取好友请求列表 | ✅ |
| PUT | `/api/friends/requests/:id` | 处理好友请求 | ✅ |
| GET | `/api/friends` | 获取好友列表 | ✅ |
| DELETE | `/api/friends/:friendId` | 删除好友 | ✅ |
| PUT | `/api/friends/:friendId/remark` | 更新好友备注 | ✅ |
| GET | `/api/friends/search` | 搜索好友 | ✅ |

#### 7.1 发送好友请求

**请求参数**：
```json
{
  "toUserId": 2,
  "message": "我是张三，想加您为好友"
}
```

#### 7.2 获取好友请求列表

**查询参数**：
- `type`: 类型（sent-发送的，received-收到的，默认received）

#### 7.3 处理好友请求

**路径参数**：
- `id`: 好友请求ID

**请求参数**：
```json
{
  "action": "accept"
}
```

**action** 值：
- `accept`: 接受
- `reject`: 拒绝

#### 7.4 获取好友列表

**查询参数**：
- `keyword`: 搜索关键词（可选）
- `page`: 页码
- `pageSize`: 每页数量

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "friends": [
      {
        "userId": 2,
        "username": "好友名",
        "avatarUrl": "头像URL",
        "remarkName": "备注名"
      }
    ],
    "total": 10,
    "page": 1,
    "pageSize": 50
  }
}
```

#### 7.5 更新好友备注

**路径参数**：
- `friendId`: 好友ID

**请求参数**：
```json
{
  "remarkName": "新备注"
}
```

**成功响应**：
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

#### 7.6 搜索好友

**查询参数**：
- `keyword`: 搜索关键词（必填）
- `page`: 页码
- `pageSize`: 每页数量

**成功响应**：
```json
{
  "code": 200,
  "message": "搜索成功",
  "data": {
    "friends": [],
    "total": 0,
    "page": 1,
    "pageSize": 20
  }
}
```

---

### 8. 管理员接口

所有管理员接口都需要管理员权限。

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/admin/users` | 获取所有用户列表 | ✅ 管理员 |
| GET | `/api/admin/users/:userId/posts` | 查看用户动态 | ✅ 管理员 |
| GET | `/api/admin/users/:userId/friends` | 查看用户好友 | ✅ 管理员 |
| PUT | `/api/admin/users/:userId/password` | 重置用户密码 | ✅ 管理员 |
| PUT | `/api/admin/users/:userId/ban` | 封禁用户 | ✅ 管理员 |
| PUT | `/api/admin/users/:userId/unban` | 解封用户 | ✅ 管理员 |
| DELETE | `/api/admin/users/:userId` | 删除用户 | ✅ 管理员 |
| DELETE | `/api/admin/posts/:id` | 删除用户动态 | ✅ 管理员 |
| DELETE | `/api/admin/comments/:id` | 删除评论 | ✅ 管理员 |

#### 8.1 获取所有用户列表

**查询参数**：
- `keyword`: 搜索关键词（可选）
- `page`: 页码
- `pageSize`: 每页数量

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "users": [],
    "total": 25,
    "page": 1,
    "pageSize": 20
  }
}
```

#### 8.2 查看用户动态

**路径参数**：
- `userId`: 用户ID

**查询参数**：
- `page`: 页码
- `pageSize`: 每页数量

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "posts": [],
    "total": 0,
    "page": 1,
    "pageSize": 20
  }
}
```

#### 8.3 查看用户好友

**路径参数**：
- `userId`: 用户ID

**查询参数**：
- `page`: 页码
- `pageSize`: 每页数量

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "friends": [],
    "total": 0,
    "page": 1,
    "pageSize": 20
  }
}
```

#### 8.4 重置用户密码

**路径参数**：
- `userId`: 用户ID

**请求参数**：
```json
{
  "newPassword": "newpassword123"
}
```

**成功响应**：
```json
{
  "code": 200,
  "message": "密码重置成功",
  "data": {
    "userId": "0000000001"
  }
}
```

#### 8.5 封禁用户

**路径参数**：
- `userId`: 用户ID

**成功响应**：
```json
{
  "code": 200,
  "message": "封禁成功",
  "data": {
    "userId": "0000000001"
  }
}
```

#### 8.6 解封用户

**路径参数**：
- `userId`: 用户ID

**成功响应**：
```json
{
  "code": 200,
  "message": "解封成功",
  "data": {
    "userId": "0000000001"
  }
}
```

#### 8.7 删除用户

**路径参数**：
- `userId`: 用户ID

**成功响应**：
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

#### 8.8 删除用户动态

**路径参数**：
- `id`: 动态ID

**成功响应**：
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {
    "postId": 46
  }
}
```

#### 8.9 删除评论

**路径参数**：
- `id`: 评论ID

**成功响应**：
```json
{
  "code": 200,
  "message": "删除成功",
  "data": {
    "commentId": 1
  }
}
```

---

### 9. 消息接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/messages` | 发送消息 | ✅ |
| GET | `/api/messages/:peerId` | 获取消息列表 | ✅ |
| PUT | `/api/messages/:peerId/read` | 标记消息已读 | ✅ |

#### 8.1 发送消息

**请求参数**：
```json
{
  "receiverId": 2,
  "msgType": 1,
  "contentPreview": "消息内容",
  "fileUrl": "",
  "fileSize": 0,
  "isEncrypted": false,
  "deviceId": "",
  "serverMsgId": ""
}
```

**消息类型**：
- 1: 文本消息
- 2: 图片消息
- 3: 视频消息
- 4: 文件消息

#### 8.2 获取消息列表

**路径参数**：
- `peerId`: 对方用户ID

**查询参数**：
- `page`: 页码
- `pageSize`: 每页数量

#### 8.3 标记消息已读

**路径参数**：
- `peerId`: 对方用户ID

---

### 9. 会话接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/conversations` | 获取会话列表 | ✅ |
| PUT | `/api/conversations/:peerId/pin` | 置顶会话 | ✅ |
| DELETE | `/api/conversations/:peerId/pin` | 取消置顶会话 | ✅ |
| PUT | `/api/conversations/:peerId/mute` | 静音会话 | ✅ |
| DELETE | `/api/conversations/:peerId/mute` | 取消静音会话 | ✅ |
| DELETE | `/api/conversations/:peerId` | 删除会话 | ✅ |
| GET | `/api/conversations/unread` | 获取未读消息数 | ✅ |

#### 9.1 获取会话列表

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "peerId": 2,
      "peerUsername": "好友名",
      "peerAvatar": "头像URL",
      "lastMessage": "最后一条消息",
      "unreadCount": 5,
      "isPinned": false,
      "isMuted": false,
      "updatedAt": "2024-12-30T10:00:00Z"
    }
  ]
}
```

#### 9.2 获取未读消息数

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "total": 10,
    "conversations": 3
  }
}
```

---

### 10. 标签接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/public/tags` | 获取标签列表 | ❌ |
| GET | `/public/tags/hot` | 获取热门标签 | ❌ |
| GET | `/public/tags/search` | 搜索标签 | ❌ |
| GET | `/public/tags/by-name/:name/posts` | 获取标签相关帖子 | ❌ |
| GET | `/public/tags/by-id/:id` | 获取标签详情 | ❌ |
| POST | `/api/tags` | 创建标签 | ✅ |
| PUT | `/api/tags/:id` | 更新标签 | ✅ |
| DELETE | `/api/tags/:id` | 删除标签 | ✅ |

#### 10.1 获取标签列表

**查询参数**：
- `page`: 页码（可选，默认1）
- `pageSize`: 每页数量（可选，默认50）
- `status`: 状态（可选，默认0-正常，1-禁用）

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "tags": [
      {
        "tagId": 1,
        "name": "美食",
        "color": "#FF6B6B",
        "postCount": 100
      }
    ],
    "total": 50
  }
}
```

#### 10.2 创建标签

**请求参数**：
```json
{
  "name": "美食",
  "color": "#FF6B6B",
  "icon": "food-icon",
  "description": "美食相关内容"
}
```

---

### 11. 搜索接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/search` | 搜索内容（公开） | ❌ |
| GET | `/search/hot-words` | 获取热词 | ❌ |
| GET | `/search/suggestions` | 获取搜索建议 | ❌ |
| GET | `/api/search` | 搜索内容 | ✅ |
| GET | `/api/search/hot-words` | 获取热词 | ✅ |
| GET | `/api/search/history` | 获取搜索历史 | ✅ |
| GET | `/api/search/filter` | 过滤内容 | ✅ |
| POST | `/api/search/history` | 保存搜索历史 | ✅ |
| GET | `/api/search/suggestions` | 获取搜索建议 | ✅ |

#### 11.1 搜索内容

**查询参数**：
- `keyword`: 搜索关键词
- `page`: 页码（可选，默认1）
- `pageSize`: 每页数量（可选，默认10）

**成功响应**：
```json
{
  "code": 200,
  "message": "搜索成功",
  "data": {
    "posts": [],
    "users": [],
    "tags": [],
    "total": 50
  }
}
```

#### 11.2 获取热门关键词

**成功响应**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "word": "美食",
      "searchCount": 1000
    }
  ]
}
```

#### 11.3 保存搜索历史

**请求参数**：
```json
{
  "keyword": "搜索关键词"
}
```

---

### 12. 文件上传接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/upload/file` | 上传文件 | ✅ |
| POST | `/api/upload/avatar` | 上传头像 | ✅ |

#### 12.1 通用文件上传

**请求方式**: `multipart/form-data`

**请求参数**：
```
file: File (必需) - 上传的文件
```

**文件要求**：
- 支持格式: jpg, jpeg, png, gif
- 文件大小: 最大 10MB

**响应示例**：
```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "fileId": "uuid前16位",
    "filename": "文件名",
    "originalName": "原始文件名",
    "fileSize": 文件大小,
    "fileType": ".jpg",
    "fileUrl": "http://localhost:8080/static/files/文件名"
  }
}
```

#### 12.2 头像上传

**请求方式**: `multipart/form-data`

**请求参数**：
```
avatar: File (必需) - 头像文件
```

**文件要求**：
- 支持格式: jpg, jpeg, png, gif, webp
- 文件大小: 最大 5MB
- 推荐尺寸: 正方形 (如 200x200px)

**响应示例**：
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatarUrl": "http://localhost:8080/static/avatars/文件名",
    "filename": "文件名",
    "size": 文件大小
  }
}
```

**JavaScript 示例**：
```javascript
// 上传头像
async function uploadAvatar(file) {
  const formData = new FormData();
  formData.append('avatar', file);

  const response = await fetch('http://localhost:8080/api/upload/avatar', {
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
  return result;
}
```

---

### 13. 兼容旧版接口（动态相关）

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/moments` | 创建动态 | ✅ |
| GET | `/api/moments/:id` | 获取动态详情 | ✅ |
| PATCH | `/api/moments/:id` | 编辑动态 | ✅ |
| DELETE | `/api/moments/:id` | 删除动态 | ✅ |
| GET | `/api/moments/my` | 获取我的动态 | ✅ |

> ⚠️ 注意：这些接口为保持向后兼容而保留，建议使用新的 `/api/posts` 接口。

---

## 🧪 测试账号

| 用户名 | 手机号 | 密码 |
|--------|--------|------|
| Yw166332 | 17875242005 | JiangCan030 |

---

## 🚀 快速开始示例

### JavaScript 示例

```javascript
// 登录获取token
const login = async (account, password) => {
  const response = await fetch('http://106.52.165.122:8080/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ account, password })
  });
  const data = await response.json();
  if (data.code === 200) {
    localStorage.setItem('token', data.data.token);
  }
  return data;
};

// 带认证的请求
const authRequest = async (url, options = {}) => {
  const token = localStorage.getItem('token');
  const response = await fetch(`http://106.52.165.122:8080${url}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers
    }
  });
  return await response.json();
};

// 上传头像
const uploadAvatar = async (file) => {
  const formData = new FormData();
  formData.append('avatar', file);
  const response = await fetch('http://localhost:8080/api/upload/avatar', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    },
    body: formData
  });
  return await response.json();
};

// 发布帖子
const createPost = async (postData) => {
  return await authRequest('/api/posts', {
    method: 'POST',
    body: JSON.stringify(postData)
  });
};
```

---

## ⚠️ 注意事项

1. **Token有效期**: 7天 (168小时)
2. **密码强度**: 必须包含大小写字母和数字，长度8-20位
3. **用户名规则**: 3-20个字符，支持字母、数字、中文、下划线
4. **手机号格式**: 中国大陆11位手机号
5. **文件上传限制**:
   - 通用文件: jpg, jpeg, png, gif，最大10MB
   - 头像: jpg, jpeg, png, gif, webp，最大5MB
6. **时间格式**: RFC3339 格式 `2024-12-30T10:00:00Z`
7. **可见性参数**:
   - 0: 公开
   - 1: 好友可见
   - 2: 仅自己可见

---

## 📞 常见问题

### Q: 提示"missing authorization header"
A: 请在请求头中添加 `Authorization: Bearer <token>`

### Q: Token过期了怎么办
A: 需要重新登录获取新token

### Q: 文件上传失败
A: 检查文件格式和大小是否符合要求

### Q: 如何处理401错误
A: 清除本地存储的token，重新登录

---

## 📅 更新日志

- **v1.0.0** (2024-12-30): 统一接口文档
  - 整合所有接口到一个文档
  - 标准化响应格式
  - 完善接口说明和示例
