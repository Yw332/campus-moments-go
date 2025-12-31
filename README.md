# Campus Moments Go API

## 🚀 项目简介

Campus Moments 是一个基于 Go 语言开发的校园动态分享平台后端API服务，提供用户认证、动态发布、文件上传、内容搜索等核心功能。

## 📋 技术栈

### 后端框架
- **Go 1.21+** - 主要编程语言
- **Gin** - Web框架，提供高性能HTTP服务
- **GORM** - ORM框架，支持MySQL数据库操作

### 数据库
- **MySQL 8.0** - 主数据库
- **原生SQL + GORM** - 双重数据库访问模式

### 认证与安全
- **JWT (JSON Web Token)** - 用户认证
- **bcrypt** - 密码加密
- **中间件** - 请求拦截和权限验证

### 文档
- **OpenAPI 3.0** - API接口文档
- **Swagger** - 交互式API文档

## 🏗️ 项目结构

```
campus-moments-go/
├── cmd/                    # 应用入口
│   └── api/
│       └── main.go        # 主程序入口
├── internal/               # 内部包
│   ├── handlers/          # HTTP处理器
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── moment_handler.go
│   │   └── upload_handler.go
│   ├── models/            # 数据模型
│   │   ├── user.go
│   │   ├── moment.go
│   │   └── db.go
│   ├── middleware/        # 中间件
│   │   └── auth.go
│   ├── routes/            # 路由配置
│   │   └── routes.go
│   └── service/           # 业务逻辑层
│       ├── auth_service.go
│       ├── user_service.go
│       └── moment_service.go
├── pkg/                   # 公共包
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   └── jwt/              # JWT工具
├── docs/                 # 文档目录
│   ├── api.md           # 前端接口需求文档
│   ├── backend.md       # 后端实现逻辑文档
│   ├── test.md          # 接口测试文档
│   └── test.sh          # 自动化测试脚本
├── go.mod               # Go模块文件
├── go.sum               # 依赖锁定文件
├── .env                 # 环境配置
├── .env.example         # 配置模板
└── openapi.json         # OpenAPI规范
```

## 🔧 核心模块

### 1. 认证模块 (`internal/handlers/auth_handler.go`)
- 用户注册/登录/退出
- JWT Token管理
- 密码修改
- 用户资料管理

**主要接口：**
- `POST /auth/register` - 用户注册
- `POST /auth/login` - 用户登录
- `GET /api/users/profile` - 获取用户资料
- `PUT /api/users/password` - 修改密码

### 2. 动态模块 (`internal/handlers/moment_handler.go`)
- 动态发布/编辑/删除
- 动态列表获取
- 动态详情查看

**主要接口：**
- `GET /api/moments` - 获取动态列表（支持分页）
- `POST /api/moments` - 发布动态
- `GET /api/moments/:id` - 获取动态详情
- `GET /api/moments/my` - 获取我的动态
- `PUT /api/moments/:id` - 编辑动态
- `DELETE /api/moments/:id` - 删除动态

### 3. 文件上传模块 (`internal/handlers/upload_handler.go`)
- 图片上传
- 头像上传
- 文件类型验证
- 大小限制

**主要接口：**
- `POST /api/upload/file` - 通用文件上传
- `POST /api/upload/avatar` - 头像上传

### 4. 搜索模块 (`internal/handlers/search_handler.go`)
- 内容搜索
- 热词统计
- 搜索建议
- 内容筛选

**主要接口：**
- `GET /api/search` - 搜索内容
- `GET /api/search/hot-words` - 热门关键词
- `GET /api/search/suggestions` - 搜索建议

### 5. 评论模块 (`internal/handlers/comment_handler.go`)
- 评论发布/编辑/删除
- 评论列表获取
- 评论点赞

**主要接口：**
- `GET /public/posts/:id/comments` - 获取评论列表（公开）
- `POST /api/comments/post/:postId` - 发布评论
- `PUT /api/comments/:id` - 编辑评论
- `DELETE /api/comments/:id` - 删除评论
- `POST /api/comments/:id/like` - 点赞评论

### 6. 点赞模块 (`internal/handlers/like_handler.go`)
- 动态点赞/取消点赞
- 点赞列表获取

**主要接口：**
- `POST /api/likes/post/:postId` - 点赞动态
- `GET /api/posts/:id/likes` - 获取点赞列表

### 7. 标签模块 (`internal/handlers/tag_handler.go`)
- 标签列表获取
- 热门标签
- 标签搜索

**主要接口：**
- `GET /public/tags` - 获取标签列表
- `GET /public/tags/hot` - 获取热门标签
- `GET /public/tags/search` - 搜索标签

## 🛠️ 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Git

### 安装步骤

1. **克隆项目**
```bash
git clone https://github.com/Yw332/campus-moments-go.git
cd campus-moments-go
```

2. **安装依赖**
```bash
go mod download
```

3. **配置环境变量**
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库等信息
```

4. **启动应用**
```bash
go run cmd/api/main.go
```

5. **访问服务**
- API服务：`http://localhost:8080`
- 健康检查：`http://localhost:8080/health`
- 首页：`http://localhost:8080/`

6. **运行测试**
```bash
# 运行自动化测试脚本
cd docs
./test.sh

# 或指定后端地址
API_BASE=http://localhost:8080 ./test.sh
```

### 部署

#### CNB云原生平台
项目支持CNB云原生平台部署，使用Procfile配置：

```procfile
web: ./start.sh
```

#### Docker部署（未来支持）
```bash
# 构建镜像
docker build -t campus-moments-go .

# 运行容器
docker run -p 8080:8080 campus-moments-go
```

## 📚 API文档

### 文档目录 (`docs/`)

- **[前端接口需求文档](./docs/api.md)** - 完整的前端接口需求说明，包含18个已实现的接口
- **[后端实现逻辑文档](./docs/backend.md)** - 详细的后端架构和实现逻辑
- **[接口测试文档](./docs/test.md)** - 完整的接口测试指南和示例
- **[自动化测试脚本](./docs/test.sh)** - 一键测试所有接口的Shell脚本

### 其他文档

- **OpenAPI规范**: [openapi.json](./openapi.json)
- **在线文档**: 启动服务后访问 `/swagger/index.html`

### 快速测试

运行自动化测试脚本：
```bash
cd docs
./test.sh
```

或指定后端地址：
```bash
API_BASE=http://106.52.165.122:8080 ./test.sh
```

## 🎯 功能特性

### ✅ 已实现（18个核心接口）

#### 用户认证模块（5个）
- [x] 用户注册/登录认证
- [x] JWT Token认证机制（7天有效期）
- [x] 用户资料管理
- [x] 密码修改
- [x] 退出登录（Token黑名单）

#### 动态模块（4个）
- [x] 动态发布和浏览
- [x] 动态详情查看
- [x] 我的动态列表
- [x] 动态编辑和删除

#### 评论模块（3个）
- [x] 评论发布/编辑/删除
- [x] 评论列表获取
- [x] 评论点赞功能

#### 点赞模块（1个）
- [x] 动态点赞/取消点赞
- [x] 点赞列表获取

#### 搜索模块（2个）
- [x] 内容搜索（动态和用户）
- [x] 热门关键词获取

#### 文件上传模块（2个）
- [x] 通用文件上传（图片）
- [x] 头像上传

#### 标签模块（1个）
- [x] 标签列表获取
- [x] 热门标签
- [x] 标签搜索

#### 其他功能
- [x] 响应式错误处理
- [x] 完整的API文档
- [x] 自动化测试脚本
- [x] 管理员功能

### 🚧 计划中
- [ ] 消息通知系统
- [ ] 用户关注系统
- [ ] 图片压缩处理
- [ ] 缓存机制优化
- [ ] 接口限流
- [ ] 日志系统完善
- [ ] 单元测试覆盖

## 🔒 安全特性

- **密码加密**: 使用bcrypt加密存储用户密码
- **JWT认证**: 7天有效期的安全Token认证
- **Token黑名单**: 退出登录时Token加入黑名单，防止Token复用
- **参数验证**: 严格的输入参数验证（用户名、手机号、密码格式）
- **CORS支持**: 跨域请求安全控制
- **SQL注入防护**: 使用GORM防止SQL注入
- **权限控制**: 管理员中间件，保护敏感操作

## 📊 性能优化

- **数据库连接池**: 优化数据库连接管理
- **Gin框架**: 高性能HTTP服务
- **响应缓存**: 静态资源缓存机制
- **优雅降级**: 数据库连接失败时的容错处理

## 🤝 贡献指南

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📝 更新日志

### v1.1.0 (2024-12-31)
- ✨ 完成所有18个前端需要的接口
- ✨ 实现评论系统（发布、编辑、删除、点赞）
- ✨ 实现点赞功能（动态点赞、评论点赞）
- ✨ 优化动态列表接口，匹配前端瀑布流需求
- ✨ 修复评论和点赞服务的模型映射问题
- ✨ 添加自动化测试脚本（`docs/test.sh`）
- 📚 完善文档（api.md、backend.md、test.md）
- 🐛 修复多个接口bug（评论列表、点赞功能等）

### v1.0.0 (2024-12-13)
- ✨ 初始版本发布
- ✨ 完成用户认证模块
- ✨ 实现动态发布功能
- ✨ 添加文件上传功能
- ✨ 集成内容搜索功能
- 📚 完善API文档

## 📞 联系方式

- **开发者**: Yw332
- **GitHub**: https://github.com/Yw332/campus-moments-go
- **Issues**: https://github.com/Yw332/campus-moments-go/issues

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。