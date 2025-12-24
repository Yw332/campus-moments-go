INSERT INTO posts (user_id, title, content, images, visibility, status, tags, liked_users, comments_summary, like_count, comment_count, view_count, created_at, updated_at) 
VALUES (
    '1000000001', 
    '测试动态标题', 
    '这是一条测试动态内容，用于验证前端对接是否正常！🎉', 
    '[]', 
    0, 1, 
    '["测试", "前端"]', 
    '[]', 
    '{}', 
    10, 3, 25, 
    NOW(), NOW()
);