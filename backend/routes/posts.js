module.exports = (posts) => {
  const express = require('express');
  const router = express.Router();

  // 获取所有动态
  router.get('/', (req, res) => {
    const sortedPosts = posts.sort((a, b) => 
      new Date(b.createTime) - new Date(a.createTime)
    );
    
    res.json({ 
      code: 200, 
      message: 'success', 
      data: sortedPosts 
    });
  });

  // 发布新动态
  router.post('/', (req, res) => {
    const { content } = req.body;
    
    const newPost = {
      id: posts.length + 1,
      content,
      userId: 1,
      username: 'admin',
      createTime: new Date().toISOString()
    };
    
    posts.push(newPost);
    res.json({ code: 200, message: '发布成功', data: { postId: newPost.id } });
  });

  return router;
};