-- 修复 comments 表的外键约束问题
-- 步骤1: 删除现有的外键约束
ALTER TABLE comments DROP FOREIGN KEY IF EXISTS comments_ibfk_1;

-- 步骤2: 修改 post_id 列类型为 bigint 以匹配 posts.id
ALTER TABLE comments MODIFY COLUMN post_id BIGINT NOT NULL;

-- 步骤3: 重新添加外键约束
ALTER TABLE comments ADD FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;

-- 同样修复 likes 表的 target_id 列
ALTER TABLE likes MODIFY COLUMN target_id BIGINT NOT NULL;
