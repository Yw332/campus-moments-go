-- 添加role字段到users表（管理员角色）
ALTER TABLE users ADD COLUMN role TINYINT DEFAULT 0 COMMENT '0-普通用户 1-管理员' AFTER status;

-- 创建默认管理员账号（密码需要通过应用程序重置）
-- 默认管理员用户ID为0000000001（第一个注册的用户）
-- 可以通过手动更新或使用应用程序接口设置

-- 示例：设置用户为管理员
-- UPDATE users SET role = 1 WHERE id = '0000000001';
