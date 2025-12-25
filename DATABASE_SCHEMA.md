# Campus Moments 数据库结构设计

## 用户表 (users)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | CHAR(10) | PRIMARY KEY | 用户ID |
| username | VARCHAR(20) | | 用户名 |
| password | VARCHAR(80) | | 密码（加密存储） |
| phone | VARCHAR(15) | | 手机号 |
| avatar | VARCHAR(500) | | 头像URL |
| avatar_type | TINYINT | | 头像类型 |
| avatar_updated_at | DATETIME | | 头像更新时间 |
| post_count | INT | | 发帖数量 |
| like_count | INT | | 获赞数量 |
| comment_count | INT | | 评论数量 |
| status | BIGINT | | 用户状态 |
| last_login_at | DATETIME | | 最后登录时间 |
| last_login_ip | VARCHAR(45) | | 最后登录IP |
| login_count | INT | | 登录次数 |
| created_at | DATETIME(3) | | 创建时间 |
| updated_at | DATETIME(3) | | 更新时间 |
| openid | VARCHAR(100) | | 微信openid |
| unionid | VARCHAR(100) | | 微信unionid |
| wechat_nickname | VARCHAR(100) | | 微信昵称 |
| wechat_avatar | VARCHAR(500) | | 微信头像 |
| signature | VARCHAR(200) | | 个人签名 |
| login_type | TINYINT | | 登录类型 |
| last_active_at | DATETIME | | 最后活跃时间 |

## 动态表 (posts)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 动态ID |
| user_id | CHAR(10) | | 发帖用户ID |
| title | VARCHAR(100) | | 标题 |
| content | TEXT | | 内容 |
| images | JSON | | 图片URL数组 |
| video | VARCHAR(200) | | 视频URL |
| visibility | TINYINT | | 可见性（0公开 1好友 2仅自己） |
| status | TINYINT | | 状态（0正常 1删除） |
| tags | JSON | | 标签数组 |
| liked_users | JSON | | 点赞用户数组 |
| comments_summary | JSON | | 评论摘要 |
| like_count | INT | | 点赞数量 |
| comment_count | INT | | 评论数量 |
| view_count | INT | | 浏览数量 |
| created_at | DATETIME | | 创建时间 |
| updated_at | DATETIME | | 更新时间 |

## 评论表 (comments)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 评论ID |
| post_id | INT | | 所属帖子ID |
| user_id | CHAR(10) | | 评论用户ID |
| content | VARCHAR(1000) | | 评论内容 |
| replies | JSON | | 二级评论 |
| like_count | INT | | 点赞数量 |
| is_author | TINYINT(1) | | 是否作者评论 |
| status | TINYINT | | 状态 |
| created_at | DATETIME | | 创建时间 |
| updated_at | DATETIME | | 更新时间 |

## 点赞表 (likes)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 点赞ID |
| user_id | CHAR(10) | | 用户ID |
| target_type | TINYINT | | 点赞类型（1帖子 2评论） |
| target_id | INT | | 目标ID |
| created_at | DATETIME | | 创建时间 |

## 消息表 (messages)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PRIMARY KEY | 消息ID |
| sender_id | CHAR(10) | | 发送者ID |
| receiver_id | CHAR(10) | | 接收者ID |
| msg_type | TINYINT | | 消息类型 |
| content_preview | VARCHAR(255) | | 内容预览 |
| file_url | VARCHAR(500) | | 文件URL |
| file_size | INT | | 文件大小 |
| is_encrypted | TINYINT(1) | | 是否加密 |
| is_read | TINYINT(1) | | 是否已读 |
| device_id | VARCHAR(64) | | 设备ID |
| server_msg_id | VARCHAR(64) | | 服务端消息ID |
| created_at | DATETIME | | 创建时间 |

## 对话表 (conversations)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 对话ID |
| user_id | CHAR(10) | | 用户ID |
| peer_id | CHAR(10) | | 对方用户ID |
| last_msg_id | BIGINT | | 最后消息ID |
| last_msg_preview | VARCHAR(255) | | 最后消息预览 |
| unread_count | INT | | 未读数量 |
| is_pinned | TINYINT(1) | | 是否置顶 |
| is_muted | TINYINT(1) | | 是否静音 |
| updated_at | DATETIME | | 更新时间 |

## 好友关系表 (friend_relations)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 关系ID |
| user_id | CHAR(10) | | 用户ID |
| friend_id | CHAR(10) | | 好友ID |
| relation_type | TINYINT | | 关系类型 |
| remark_name | VARCHAR(50) | | 备注名称 |
| status | TINYINT | | 状态 |
| created_at | DATETIME | | 创建时间 |
| updated_at | DATETIME | | 更新时间 |

## 好友申请表 (friend_requests)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 申请ID |
| from_user_id | CHAR(10) | | 申请人ID |
| to_user_id | CHAR(10) | | 接收人ID |
| message | VARCHAR(200) | | 申请消息 |
| status | TINYINT | | 状态 |
| expires_at | DATETIME | | 过期时间 |
| created_at | DATETIME | | 创建时间 |
| updated_at | DATETIME | | 更新时间 |

## 标签表 (tags)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 标签ID |
| name | VARCHAR(20) | | 标签名称 |
| color | VARCHAR(7) | | 标签颜色 |
| icon | VARCHAR(50) | | 标签图标 |
| description | VARCHAR(200) | | 标签描述 |
| usage_count | INT | | 使用次数 |
| last_used_at | DATETIME | | 最后使用时间 |
| status | TINYINT | | 状态 |
| created_at | DATETIME | | 创建时间 |

## 搜索历史表 (search_history)
| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | INT | PRIMARY KEY | 历史ID |
| user_id | CHAR(10) | | 用户ID |
| keyword | VARCHAR(100) | | 搜索关键词 |
| created_at | DATETIME | | 创建时间 |

---

**注意：**
1. 所有涉及用户的ID都使用CHAR(10)类型
2. posts表的status字段：0正常，1删除
3. posts表的visibility字段：0公开，1好友，2仅自己
4. 时间字段都使用DATETIME类型，部分支持毫秒精度