const http = require('http');

// 测试函数
function testAPI(method, path, data = null, token = null) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'localhost',
      port: 8080,
      path: path,
      method: method,
      headers: {
        'Content-Type': 'application/json'
      }
    };

    if (token) {
      options.headers['Authorization'] = `Bearer ${token}`;
    }

    const req = http.request(options, (res) => {
      let body = '';
      res.on('data', (chunk) => {
        body += chunk;
      });
      res.on('end', () => {
        try {
          const response = JSON.parse(body);
          resolve({
            statusCode: res.statusCode,
            response: response,
            rawBody: body
          });
        } catch (e) {
          resolve({
            statusCode: res.statusCode,
            response: body,
            rawBody: body
          });
        }
      });
    });

    req.on('error', (e) => {
      reject(e);
    });

    if (data) {
      req.write(JSON.stringify(data));
    }
    req.end();
  });
}

// 实际测试互动功能
async function testInteractionLive() {
  console.log('🔥 开始实际测试互动功能...\n');

  let testToken = null;
  let testUserId = null;

  try {
    // 1. 先注册一个测试用户
    console.log('📝 1. 注册测试用户...');
    const registerData = {
      username: `testuser${Date.now()}`,
      phone: `138${Math.floor(Math.random() * 100000000)}`,
      password: 'TestPass123'
    };

    const registerResult = await testAPI('POST', '/auth/register', registerData);
    console.log('注册结果:', JSON.stringify(registerResult.response, null, 2));

    // 2. 登录获取token
    console.log('\n🔑 2. 用户登录获取token...');
    const loginData = {
      account: registerData.username,
      password: registerData.password
    };

    const loginResult = await testAPI('POST', '/auth/login', loginData);
    console.log('登录结果:', JSON.stringify(loginResult.response, null, 2));

    if (loginResult.response.code === 200) {
      testToken = loginResult.response.data.token;
      testUserId = loginResult.response.data.userInfo.userId;
      console.log('✅ 获取token成功，用户ID:', testUserId);
    } else {
      throw new Error('登录失败');
    }

    // 3. 创建一条测试动态
    console.log('\n📝 3. 创建测试动态...');
    const momentData = {
      content: '这是一条用于测试互动功能的动态',
      tags: ['测试', '互动'],
      visibility: 0
    };

    const createMomentResult = await testAPI('POST', '/api/moments', momentData, testToken);
    console.log('创建动态结果:', JSON.stringify(createMomentResult.response, null, 2));

    let momentId = null;
    if (createMomentResult.response.code === 200) {
      momentId = createMomentResult.response.data.id;
      console.log('✅ 动态创建成功，ID:', momentId);
    }

    // 4. 测试发表评论
    console.log('\n💬 4. 测试发表评论...');
    const commentData = {
      momentId: momentId,
      content: '这是第一条测试评论'
    };

    const createCommentResult = await testAPI('POST', '/api/comments', commentData, testToken);
    console.log('发表评论结果:', JSON.stringify(createCommentResult.response, null, 2));

    let commentId = null;
    if (createCommentResult.response.code === 200) {
      commentId = createCommentResult.response.data.id;
      console.log('✅ 评论发表成功，ID:', commentId);
    }

    // 5. 测试获取评论列表
    console.log('\n📋 5. 测试获取评论列表...');
    const getCommentsResult = await testAPI('GET', `/api/moments/${momentId}/comments`, null, testToken);
    console.log('获取评论列表结果:', JSON.stringify(getCommentsResult.response, null, 2));

    // 6. 测试点赞动态
    console.log('\n👍 6. 测试点赞动态...');
    const likeMomentData = {
      targetId: momentId,
      targetType: 1 // 1:动态
    };

    const likeMomentResult = await testAPI('POST', '/api/likes', likeMomentData, testToken);
    console.log('点赞动态结果:', JSON.stringify(likeMomentResult.response, null, 2));

    // 7. 测试点赞评论
    console.log('\n❤️ 7. 测试点赞评论...');
    const likeCommentData = {
      targetId: commentId,
      targetType: 2 // 2:评论
    };

    const likeCommentResult = await testAPI('POST', '/api/likes', likeCommentData, testToken);
    console.log('点赞评论结果:', JSON.stringify(likeCommentResult.response, null, 2));

    // 8. 测试再次点赞评论（应该取消点赞）
    console.log('\n🔄 8. 测试取消点赞评论...');
    const unlikeCommentResult = await testAPI('POST', '/api/likes', likeCommentData, testToken);
    console.log('取消点赞评论结果:', JSON.stringify(unlikeCommentResult.response, null, 2));

    // 9. 测试回复评论
    console.log('\n🔄 9. 测试回复评论...');
    const replyCommentData = {
      momentId: momentId,
      content: '这是对第一条评论的回复',
      parentId: commentId,
      replyToUserId: testUserId.toString()
    };

    const replyCommentResult = await testAPI('POST', '/api/comments', replyCommentData, testToken);
    console.log('回复评论结果:', JSON.stringify(replyCommentResult.response, null, 2));

    // 10. 再次获取评论列表查看嵌套结构
    console.log('\n📋 10. 查看嵌套评论列表...');
    const updatedCommentsResult = await testAPI('GET', `/api/moments/${momentId}/comments`, null, testToken);
    console.log('更新后评论列表:', JSON.stringify(updatedCommentsResult.response, null, 2));

    // 11. 测试删除评论
    console.log('\n🗑️ 11. 测试删除评论...');
    const deleteCommentResult = await testAPI('DELETE', `/api/comments/${commentId}`, null, testToken);
    console.log('删除评论结果:', JSON.stringify(deleteCommentResult.response, null, 2));

    // 12. 最终检查动态列表，看计数是否正确
    console.log('\n📊 12. 最终检查动态数据...');
    const finalMomentsResult = await testAPI('GET', '/moments', null, testToken);
    console.log('最终动态数据:', JSON.stringify(finalMomentsResult.response, null, 2));

    console.log('\n🎉 互动功能测试完成！');

  } catch (error) {
    console.error('❌ 测试过程中发生错误:', error);
  }
}

// 运行测试
testInteractionLive().catch(console.error);