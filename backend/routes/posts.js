module.exports = (posts) => {
  const express = require('express');
  const router = express.Router();

  const authMiddleware = require('../middleware/auth');

  // 发布动态 - 需要登录
  router.post('/', authMiddleware, (req, res) => {
    // 现在可以通过 req.user 获取当前用户信息
    const currentUser = req.user;

    const { content, tags = [] } = req.body;

    const newPost = {
      id: posts.length + 1,
      content,
      userId: currentUser.userId,      // 从token中获取用户ID
      username: currentUser.username,  // 从token中获取用户名
      createTime: new Date().toISOString(),
      tags
    };

    posts.push(newPost);

    res.json({
      code: 200,
      message: '发布成功',
      data: { postId: newPost.id }
    });
  });

  // 获取动态列表 - 可以不需要登录
  router.get('/', (req, res) => {
    // 保持原有代码不变
    const sortedPosts = posts.sort((a, b) =>
      new Date(b.createTime) - new Date(a.createTime)
    );

    res.json({
      code: 200,
      message: 'success',
      data: {
        list: sortedPosts,  // 按照文档要求，数据放在 list 字段
        total: sortedPosts.length
      }
    });
  });

  // 获取单条动态详情 - 可选登录（不强制）
  router.get('/:id', (req, res) => {
    const postId = parseInt(req.params.id, 10);
    if (Number.isNaN(postId)) {
      return res.json({ code: 400, message: '无效的动态 ID', data: null });
    }

    const post = posts.find(p => p.id === postId);
    if (!post) {
      return res.json({ code: 404, message: '动态未找到', data: null });
    }

    return res.json({ code: 200, message: 'success', data: { post } });
  });

  // 更新动态 - 需要登录且只能更新自己的动态
  router.put('/:id', authMiddleware, (req, res) => {
    const currentUser = req.user;
    const postId = parseInt(req.params.id, 10);
    const { content, tags } = req.body;

    const post = posts.find(p => p.id === postId);
    if (!post) {
      return res.json({ code: 404, message: '动态未找到', data: null });
    }

    if (post.userId !== currentUser.userId) {
      return res.json({ code: 403, message: '无权修改他人动态', data: null });
    }

    // 只更新允许的字段
    if (typeof content === 'string') post.content = content;
    if (Array.isArray(tags)) post.tags = tags;
    post.updateTime = new Date().toISOString();

    return res.json({ code: 200, message: '更新成功', data: { postId: post.id, post } });
  });

  // 部分更新（PATCH） - 只更新提供的字段，响应格式与 PUT 保持一致
  router.patch('/:id', authMiddleware, (req, res) => {
    const currentUser = req.user;
    const postId = parseInt(req.params.id, 10);
    const { content, tags } = req.body;

    const post = posts.find(p => p.id === postId);
    if (!post) {
      return res.json({ code: 404, message: '动态未找到', data: null });
    }

    if (post.userId !== currentUser.userId) {
      return res.json({ code: 403, message: '无权修改他人动态', data: null });
    }

    let changed = false;
    if (typeof content === 'string') { post.content = content; changed = true; }
    if (Array.isArray(tags)) { post.tags = tags; changed = true; }
    if (changed) post.updateTime = new Date().toISOString();

    return res.json({ code: 200, message: '更新成功', data: { postId: post.id, post } });
  });

  // 删除动态 - 需要登录且只能删除自己的动态
  router.delete('/:id', authMiddleware, (req, res) => {
    const currentUser = req.user;
    const postId = parseInt(req.params.id, 10);

    const index = posts.findIndex(p => p.id === postId);
    if (index === -1) {
      return res.json({ code: 404, message: '动态未找到', data: null });
    }

    const post = posts[index];
    if (post.userId !== currentUser.userId) {
      return res.json({ code: 403, message: '无权删除他人动态', data: null });
    }

    posts.splice(index, 1);

    return res.json({ code: 200, message: '删除成功', data: { postId: postId } });
  });

  return router;
};