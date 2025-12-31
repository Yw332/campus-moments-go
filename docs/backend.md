# 后端实现逻辑文档

## 项目概述

Campus Moments 后端项目基于 Go 语言开发，使用 Gin 框架提供 RESTful API 服务，支持前端 uni-app 项目的所有功能需求。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **认证**: JWT (JSON Web Token)
- **密码加密**: bcrypt

## 项目结构

```
campus-moments-go/
├── cmd/api/main.go          # 应用入口
├── internal/
│   ├── handlers/            # HTTP处理器层
│   │   ├── auth_handler.go      # 认证相关
│   │   ├── user_handler.go      # 用户相关
│   │   ├── moment_handler.go    # 动态相关
│   │   ├── comment_handler.go   # 评论相关
│   │   ├── like_handler.go      # 点赞相关
│   │   ├── upload_handler.go    # 文件上传
│   │   ├── search_handler.go    # 搜索相关
│   │   └── ...
│   ├── models/              # 数据模型层
│   │   ├── user.go
│   │   ├── moment.go
│   │   └── ...
│   ├── service/             # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── moment_service.go
│   │   └── ...
│   ├── middleware/          # 中间件
│   │   └── auth.go             # JWT认证中间件
│   └── routes/              # 路由配置
│       └── routes.go
├── pkg/                     # 公共包
│   ├── jwt/                 # JWT工具
│   ├── database/            # 数据库连接
│   └── config/              # 配置管理
└── migrations/              # 数据库迁移
```

## 核心功能实现

### 1. 用户认证模块

#### 1.1 用户注册 (`POST /auth/register`)
**实现位置**: `internal/handlers/auth_handler.go` -> `Register()`

**实现逻辑**:
1. 接收注册请求（username, phone, password）
2. 验证用户名格式（3-20字符，字母数字中文下划线）
3. 验证手机号格式（11位，1开头）
4. 验证密码强度（8-20位，包含大小写字母和数字）
5. 检查用户名和手机号是否已存在
6. 使用bcrypt加密密码
7. 生成10位字符串ID（自动递增）
8. 创建用户记录
9. 返回用户信息（不含密码）

**关键代码**:
```go
// 密码加密
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

// 生成用户ID
newID := maxID + 1
idStr := fmt.Sprintf("%010d", newID) // 10位数字，前面补零
```

#### 1.2 用户登录 (`POST /auth/login`)
**实现位置**: `internal/handlers/auth_handler.go` -> `Login()`

**实现逻辑**:
1. 接收登录请求（account可以是用户名或手机号，password）
2. 先尝试从admins表查找（管理员登录）
3. 如果未找到，从users表查找（普通用户）
4. 验证密码（bcrypt比较）
5. 检查用户状态（是否被禁用或锁定）
6. 生成JWT Token（有效期7天）
7. 返回Token和用户信息

**JWT Token生成**:
```go
token, err := jwt.GenerateToken(userIDInt, user.Username)
// Token包含: userID, username, expiresAt
```

#### 1.3 获取用户资料 (`GET /api/users/profile`)
**实现位置**: `internal/handlers/auth_handler.go` -> `GetProfile()`

**实现逻辑**:
1. 从JWT Token中获取userID（通过中间件）
2. 查询用户信息
3. 返回用户资料（不含敏感信息）

#### 1.4 修改密码 (`PUT /api/users/password`)
**实现位置**: `internal/handlers/auth_handler.go` -> `ChangePassword()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 接收旧密码和新密码
3. 验证旧密码是否正确
4. 验证新密码强度
5. 使用bcrypt加密新密码
6. 更新数据库

#### 1.5 退出登录 (`POST /api/auth/logout`)
**实现位置**: `internal/handlers/auth_handler.go` -> `Logout()`

**实现逻辑**:
1. 从请求头获取Token
2. 解析Token获取过期时间
3. 将Token添加到黑名单（直到原定过期时间）
4. 返回成功响应

**Token黑名单机制**:
```go
blacklist := token_blacklist.GetInstance()
blacklist.AddToken(token, claims.ExpiresAt.Time)
```

### 2. 动态模块

#### 2.1 获取动态列表 (`GET /api/moments`)
**实现位置**: `internal/handlers/moment_handler.go` -> `GetMoments()`

**实现逻辑**:
1. 接收分页参数（page, pageSize）
2. 可选：按userId筛选
3. 查询posts表（status=0，正常状态）
4. 预加载用户信息（Preload User）
5. 按创建时间倒序排列
6. 转换为前端需要的格式：
   - 提取第一张图片作为imageUrl
   - 获取作者用户名
   - 格式化创建时间
7. 返回列表和分页信息

**数据格式转换**:
```go
// 提取第一张图片
var images []string
json.Unmarshal(moment.Images, &images)
imageUrl := images[0] if len(images) > 0

// 获取作者
author := moment.User.Username if moment.User != nil

// 格式化时间
createTime := moment.CreatedAt.Format("2006-01-02 15:04")
```

#### 2.2 发布动态 (`POST /api/moments`)
**实现位置**: `internal/handlers/moment_handler.go` -> `CreateMoment()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 接收请求参数（title, content, tags, images）
3. 处理图片数组：
   - 优先使用images数组
   - 如果没有，从media中提取
   - 转换为JSON格式存储
4. 转换tags为JSON格式
5. 创建动态记录
6. 预加载用户信息
7. 返回创建的动态

**图片处理**:
```go
if len(req.Images) > 0 {
    imagesJSON, _ = json.Marshal(req.Images)
} else if len(req.Media) > 0 {
    // 从Media中提取图片URL
    var imageURLs []string
    for _, media := range req.Media {
        if media.Type == "image" {
            imageURLs = append(imageURLs, media.URL)
        }
    }
    imagesJSON, _ = json.Marshal(imageURLs)
}
```

#### 2.3 获取动态详情 (`GET /api/moments/:id`)
**实现位置**: `internal/handlers/moment_handler.go` -> `GetMomentDetail()`

**实现逻辑**:
1. 解析动态ID
2. 查询动态（status=0）
3. 预加载用户信息
4. 构建响应数据（包含author字段兼容前端）
5. 返回动态详情

#### 2.4 获取我的动态 (`GET /api/moments/my`)
**实现位置**: `internal/handlers/moment_handler.go` -> `GetUserMoments()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 接收分页参数
3. 查询该用户的所有动态
4. 返回列表和分页信息

### 3. 评论模块

#### 3.1 获取评论列表 (`GET /public/posts/:id/comments`)
**实现位置**: `internal/handlers/comment_handler.go` -> `GetCommentList()`

**实现逻辑**:
1. 获取postId参数
2. 查询该动态的所有评论
3. 预加载评论者信息
4. 返回评论列表

#### 3.2 发布评论 (`POST /api/comments/post/:postId`)
**实现位置**: `internal/handlers/comment_handler.go` -> `CreateComment()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 获取postId参数
3. 接收评论内容
4. 创建评论记录
5. 更新动态的评论数
6. 返回评论信息

#### 3.3 点赞评论 (`POST /api/comments/:id/like`)
**实现位置**: `internal/handlers/comment_handler.go` -> `LikeComment()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 获取评论ID
3. 检查是否已点赞
4. 如果未点赞，添加点赞记录并增加点赞数
5. 如果已点赞，取消点赞并减少点赞数
6. 返回点赞状态

### 4. 点赞模块

#### 4.1 点赞动态 (`POST /api/likes/post/:postId`)
**实现位置**: `internal/handlers/like_handler.go` -> `LikePost()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 获取postId参数
3. 检查是否已点赞
4. 如果未点赞：
   - 添加点赞记录
   - 更新动态的liked_users JSON字段
   - 增加like_count
5. 如果已点赞：
   - 删除点赞记录
   - 更新liked_users
   - 减少like_count
6. 返回点赞状态

### 5. 文件上传模块

#### 5.1 通用文件上传 (`POST /api/upload/file`)
**实现位置**: `internal/handlers/upload_handler.go` -> `UploadFile()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 接收multipart/form-data格式的文件
3. 验证文件类型（图片：jpg, jpeg, png, gif, webp）
4. 验证文件大小（最大10MB）
5. 生成唯一文件名（时间戳+随机数）
6. 保存文件到uploads目录
7. 返回文件URL

**文件保存路径**:
```
uploads/{year}/{month}/{filename}
```

#### 5.2 头像上传 (`POST /api/upload/avatar`)
**实现位置**: `internal/handlers/upload_handler.go` -> `UploadAvatar()`

**实现逻辑**:
1. 从JWT Token获取userID
2. 接收头像文件
3. 验证文件类型和大小
4. 保存文件
5. 更新用户的avatar_url字段
6. 返回头像URL

### 6. 搜索模块

#### 6.1 搜索内容 (`GET /api/search`)
**实现位置**: `internal/handlers/search_handler.go` -> `SearchContent()`

**实现逻辑**:
1. 接收搜索关键词（keyword）
2. 接收分页参数
3. 在posts表中搜索：
   - 标题包含关键词
   - 内容包含关键词
4. 在users表中搜索用户名
5. 合并结果并去重
6. 转换为前端需要的格式（id->postId, user->author）
7. 返回搜索结果和分页信息

#### 6.2 获取热门关键词 (`GET /api/search/hot-words`)
**实现位置**: `internal/handlers/search_handler.go` -> `GetHotWords()`

**实现逻辑**:
1. 查询搜索历史表
2. 统计关键词出现次数
3. 按次数排序
4. 返回前N个热门关键词

## 中间件实现

### JWT认证中间件
**实现位置**: `internal/middleware/auth.go`

**实现逻辑**:
1. 从请求头获取Authorization
2. 提取Bearer Token
3. 验证Token格式
4. 检查Token是否在黑名单中
5. 解析Token获取userID和username
6. 将userID存储到context中
7. 继续处理请求

**关键代码**:
```go
auth := c.GetHeader("Authorization")
if !strings.HasPrefix(auth, "Bearer ") {
    c.JSON(401, gin.H{"code": 401, "message": "未认证"})
    c.Abort()
    return
}

token := strings.TrimPrefix(auth, "Bearer ")
claims, err := jwt.ParseToken(token)
c.Set("userID", claims.UserID)
c.Set("username", claims.Username)
```

## 数据库设计

### 核心表结构

#### users表
- `id`: char(10) - 用户ID
- `username`: varchar(50) - 用户名
- `phone`: varchar(20) - 手机号
- `password`: varchar(255) - 加密后的密码
- `avatar_url`: varchar(200) - 头像URL
- `status`: tinyint - 状态（0正常，2禁用，3锁定）
- `role`: tinyint - 角色（0普通用户，1管理员）

#### posts表（动态表）
- `id`: int - 动态ID
- `user_id`: char(10) - 发布者ID
- `title`: varchar(100) - 标题
- `content`: text - 内容
- `images`: json - 图片数组
- `tags`: json - 标签数组
- `like_count`: int - 点赞数
- `comment_count`: int - 评论数
- `status`: tinyint - 状态（0正常，1删除）

#### comments表
- `id`: int - 评论ID
- `post_id`: int - 动态ID
- `user_id`: char(10) - 评论者ID
- `content`: text - 评论内容
- `like_count`: int - 点赞数

## 响应格式规范

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
- `200`: 成功
- `400`: 请求参数错误
- `401`: 未认证或认证失败
- `403`: 禁止访问
- `404`: 资源不存在
- `500`: 服务器内部错误

## 安全机制

### 1. 密码加密
- 使用bcrypt算法加密存储
- 默认cost为10

### 2. JWT Token
- 有效期7天
- 包含userID和username
- 支持Token黑名单机制

### 3. 参数验证
- 使用Gin的binding验证
- 验证用户名、手机号、密码格式
- 验证文件类型和大小

### 4. SQL注入防护
- 使用GORM的预编译语句
- 避免直接拼接SQL

## 接口优化记录

### 2024-12-XX 前端集成优化

1. **GET /api/moments接口**
   - 新增：支持获取动态列表
   - 优化：返回格式匹配前端需求（id, title, author, imageUrl, likeCount, createTime）

2. **POST /api/moments接口**
   - 优化：支持title参数
   - 优化：支持images数组参数（前端格式）
   - 兼容：保留原有的media格式支持

3. **数据格式转换**
   - 统一将数据库字段转换为前端需要的格式
   - id -> postId
   - user -> author
   - images数组提取第一张作为imageUrl

## 部署说明

### 环境变量配置
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=campus_moments
JWT_SECRET=your-secret-key
```

### 启动服务
```bash
go run cmd/api/main.go
```

### 健康检查
```bash
curl http://localhost:8080/health
```

## 注意事项

1. 所有需要认证的接口路径以`/api/`开头
2. 公开接口路径以`/public/`或`/auth/`开头
3. 文件上传使用multipart/form-data格式
4. 图片存储在uploads目录，需要配置静态文件服务
5. Token黑名单使用内存存储，重启后失效（生产环境建议使用Redis）

