# Campus Moments Go API 接口文档

基于数据库结构设计的完整后端接口文档。

## 数据库表结构

### 核心表
- **users** - 用户表
- **posts** - 帖子表  
- **comments** - 评论表
- **likes** - 点赞表
- **messages** - 消息表
- **conversations** - 会话表
- **friend_relations** - 好友关系表
- **friend_requests** - 好友请求表
- **tags** - 标签表
- **search_history** - 搜索历史表

## 接口设计

### 1. 认证相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/auth/register` | 用户注册 | ❌ |
| POST | `/auth/login` | 用户登录 | ❌ |
| POST | `/auth/logout` | 用户登出 | ✅ |
| POST | `/auth/send-verification` | 发送验证码 | ❌ |
| POST | `/auth/verify-and-reset` | 验证并重置密码 | ❌ |

### 2. 帖子相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/posts` | 获取帖子列表 | ❌ |
| GET | `/posts/:id` | 获取帖子详情 | ❌ |
| POST | `/api/posts` | 创建帖子 | ✅ |
| PUT | `/api/posts/:id` | 更新帖子 | ✅ |
| DELETE | `/api/posts/:id` | 删除帖子 | ✅ |
| GET | `/api/posts/my` | 获取我的帖子 | ✅ |
| GET | `/api/posts/user/:userId` | 获取用户帖子 | ✅ |
| GET | `/api/posts/:id/likes` | 获取帖子点赞列表 | ✅ |

#### 创建帖子请求示例
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

#### 可见性参数
- 0: 公开
- 1: 好友可见
- 2: 仅自己可见

### 3. 评论相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/posts/:postId/comments` | 获取评论列表 | ❌ |
| POST | `/api/comments/post/:postId` | 创建评论 | ✅ |
| PUT | `/api/comments/:id` | 更新评论 | ✅ |
| DELETE | `/api/comments/:id` | 删除评论 | ✅ |
| POST | `/api/comments/:id/like` | 点赞评论 | ✅ |
| POST | `/api/comments/:id/reply` | 回复评论 | ✅ |
| GET | `/api/comments/:id/likes` | 获取评论点赞列表 | ✅ |

#### 创建评论请求示例
```json
{
  "content": "评论内容",
  "replies": []
}
```

### 4. 点赞相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/likes/post/:postId` | 点赞/取消点赞帖子 | ✅ |
| GET | `/api/likes/posts/:postId` | 获取帖子点赞列表 | ✅ |
| GET | `/api/likes/comments/:commentId` | 获取评论点赞列表 | ✅ |
| GET | `/api/likes/users/:userId` | 获取用户点赞列表 | ✅ |

### 5. 好友相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/friends/request` | 发送好友请求 | ✅ |
| GET | `/api/friends/requests` | 获取好友请求列表 | ✅ |
| PUT | `/api/friends/requests/:id` | 处理好友请求 | ✅ |
| GET | `/api/friends` | 获取好友列表 | ✅ |
| DELETE | `/api/friends/:friendId` | 删除好友 | ✅ |
| PUT | `/api/friends/:friendId/remark` | 更新好友备注 | ✅ |
| GET | `/api/friends/search` | 搜索好友 | ✅ |

#### 发送好友请求示例
```json
{
  "toUserId": "user123",
  "message": "我是张三，想加您为好友"
}
```

#### 处理好友请求示例
```json
{
  "action": "accept"  // accept 或 reject
}
```

### 6. 消息相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/messages` | 发送消息 | ✅ |
| GET | `/api/messages/:peerId` | 获取消息列表 | ✅ |
| PUT | `/api/messages/:peerId/read` | 标记消息已读 | ✅ |

#### 发送消息示例
```json
{
  "receiverId": "user123",
  "msgType": 1,
  "contentPreview": "消息内容",
  "fileUrl": "",
  "fileSize": 0,
  "isEncrypted": false,
  "deviceId": "",
  "serverMsgId": ""
}
```

#### 消息类型
- 1: 文本消息
- 2: 图片消息
- 3: 视频消息
- 4: 文件消息

### 7. 会话相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/conversations` | 获取会话列表 | ✅ |
| PUT | `/api/conversations/:peerId/pin` | 置顶会话 | ✅ |
| DELETE | `/api/conversations/:peerId/pin` | 取消置顶会话 | ✅ |
| PUT | `/api/conversations/:peerId/mute` | 静音会话 | ✅ |
| DELETE | `/api/conversations/:peerId/mute` | 取消静音会话 | ✅ |
| DELETE | `/api/conversations/:peerId` | 删除会话 | ✅ |
| GET | `/api/conversations/unread` | 获取未读消息数 | ✅ |

### 8. 标签相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/tags` | 获取标签列表 | ❌ |
| GET | `/tags/hot` | 获取热门标签 | ❌ |
| GET | `/tags/search` | 搜索标签 | ❌ |
| GET | `/tags/:name/posts` | 获取标签相关帖子 | ❌ |
| GET | `/tags/:id` | 获取标签详情 | ❌ |
| POST | `/api/tags` | 创建标签 | ✅ |
| PUT | `/api/tags/:id` | 更新标签 | ✅ |
| DELETE | `/api/tags/:id` | 删除标签 | ✅ |

#### 创建标签示例
```json
{
  "name": "美食",
  "color": "#FF6B6B",
  "icon": "food-icon",
  "description": "美食相关内容"
}
```

### 9. 用户相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/users/profile` | 获取用户信息 | ✅ |
| PUT | `/api/users/profile` | 更新用户信息 | ✅ |
| PUT | `/api/users/password` | 修改密码 | ✅ |

#### 更新用户信息示例
```json
{
  "username": "新用户名",
  "phone": "13800138000",
  "avatar": "头像URL",
  "signature": "个性签名",
  "wechatNickname": "微信昵称"
}
```

### 10. 搜索相关接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/search` | 搜索内容 | ✅ |
| GET | `/api/search/hot-words` | 获取热词 | ✅ |
| GET | `/api/search/history` | 获取搜索历史 | ✅ |
| POST | `/api/search/history` | 保存搜索历史 | ✅ |
| GET | `/api/search/suggestions` | 获取搜索建议 | ✅ |
| GET | `/api/search/filter` | 过滤内容 | ✅ |

### 11. 文件上传接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/upload/file` | 上传文件 | ✅ |
| POST | `/api/upload/avatar` | 上传头像 | ✅ |

### 12. 系统接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/` | 首页 | ❌ |
| GET | `/health` | 健康检查 | ❌ |

## 响应格式

### 成功响应
```json
{
  "code": 0,
  "msg": "操作成功",
  "data": {}
}
```

### 分页响应
```json
{
  "code": 0,
  "msg": "获取成功",
  "data": {
    "posts": [],
    "total": 100,
    "page": 1,
    "pageSize": 20
  }
}
```

### 错误响应
```json
{
  "error": "错误信息"
}
```

## 权限说明

### 可见性规则
- 公开帖子：所有人可见
- 好友可见：仅好友可见
- 私密帖子：仅作者可见

### 操作权限
- 帖子：作者可以编辑、删除
- 评论：作者可以编辑、删除
- 好友关系：双向关系，双方都可删除

## 数据类型

### 时间格式
所有时间字段使用 RFC3339 格式：`2023-12-28T15:30:00Z`

### JSON字段
- images: 图片URL数组
- tags: 标签数组
- replies: 评论回复数组
- likedUsers: 点赞用户ID数组

## 注意事项

1. 所有需要认证的接口都需要在请求头中携带JWT token：`Authorization: Bearer <token>`
2. 分页参数：`page`（页码，从1开始）、`pageSize`（每页数量）
3. 搜索接口支持关键词过滤和分页
4. 文件上传需要处理文件大小限制和格式验证
5. 敏感操作需要记录日志和进行权限验证

## 兼容性说明

为保持向后兼容，保留了原有的 `/api/moments` 相关接口，与新的 `/api/posts` 接口并存。