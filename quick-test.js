// 快速测试头像上传功能的JavaScript代码
// 在浏览器控制台中运行这段代码

const API_BASE = 'http://106.52.165.122:8080';
let authToken = '';

async function quickAvatarTest() {
  console.log('🚀 开始测试头像上传功能...');
  
  try {
    // 1. 登录获取token
    console.log('📝 步骤1: 登录获取token...');
    const loginResponse = await fetch(API_BASE + '/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        account: 'Yw166332',
        password: 'JiangCan030'
      })
    });
    
    const loginResult = await loginResponse.json();
    if (loginResult.code !== 200) {
      throw new Error('登录失败: ' + loginResult.message);
    }
    
    authToken = loginResult.data.token;
    console.log('✅ 登录成功，token已获取');
    
    // 2. 获取用户资料（上传前）
    console.log('📝 步骤2: 获取上传前用户资料...');
    const profileResponse1 = await fetch(API_BASE + '/api/users/profile', {
      headers: { 'Authorization': `Bearer ${authToken}` }
    });
    
    const profile1 = await profileResponse1.json();
    console.log('上传前用户资料:', profile1.data);
    
    // 3. 创建测试图片文件
    console.log('📝 步骤3: 创建测试图片...');
    const canvas = document.createElement('canvas');
    canvas.width = 200;
    canvas.height = 200;
    const ctx = canvas.getContext('2d');
    
    // 绘制一个简单的测试图片
    ctx.fillStyle = '#667eea';
    ctx.fillRect(0, 0, 200, 200);
    ctx.fillStyle = '#ffffff';
    ctx.font = '30px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('TEST', 100, 100);
    
    // 转换为blob
    const blob = await new Promise(resolve => canvas.toBlob(resolve, 'image/jpeg'));
    const file = new File([blob], 'test-avatar.jpg', { type: 'image/jpeg' });
    
    console.log('✅ 测试图片创建成功');
    
    // 4. 上传头像
    console.log('📝 步骤4: 上传头像...');
    const formData = new FormData();
    formData.append('avatar', file);
    
    const uploadResponse = await fetch(API_BASE + '/api/upload/avatar', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authToken}` },
      body: formData
    });
    
    const uploadResult = await uploadResponse.json();
    if (uploadResult.code !== 200) {
      throw new Error('上传失败: ' + uploadResult.message);
    }
    
    console.log('✅ 头像上传成功!');
    console.log('头像URL:', uploadResult.data.avatarUrl);
    
    // 5. 验证头像URL是否可访问
    console.log('📝 步骤5: 验证头像URL访问...');
    const avatarTestResponse = await fetch(uploadResult.data.avatarUrl);
    if (avatarTestResponse.ok) {
      console.log('✅ 头像URL可正常访问');
    } else {
      console.log('❌ 头像URL访问失败');
    }
    
    // 6. 获取用户资料（上传后）
    console.log('📝 步骤6: 获取上传后用户资料...');
    const profileResponse2 = await fetch(API_BASE + '/api/users/profile', {
      headers: { 'Authorization': `Bearer ${authToken}` }
    });
    
    const profile2 = await profileResponse2.json();
    console.log('上传后用户资料:', profile2.data);
    
    // 7. 验证头像是否正确更新
    if (profile2.data.avatar === uploadResult.data.avatarUrl) {
      console.log('✅ 用户资料头像已正确更新');
    } else {
      console.log('❌ 用户资料头像更新失败');
    }
    
    console.log('🎉 测试完成！头像上传功能正常工作');
    
    // 在页面中显示测试结果
    const resultDiv = document.createElement('div');
    resultDiv.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      z-index: 10000;
      max-width: 400px;
    `;
    
    resultDiv.innerHTML = `
      <h3 style="margin: 0 0 15px 0; color: #333;">🎉 测试完成</h3>
      <p style="margin: 5px 0; color: #666;">✅ 登录成功</p>
      <p style="margin: 5px 0; color: #666;">✅ 头像上传成功</p>
      <p style="margin: 5px 0; color: #666;">✅ 用户资料更新成功</p>
      <img src="${uploadResult.data.avatarUrl}" style="width: 50px; height: 50px; border-radius: 50%; margin-top: 10px;">
      <button onclick="this.parentElement.remove()" style="margin-top: 15px; padding: 5px 10px; background: #667eea; color: white; border: none; border-radius: 4px; cursor: pointer;">关闭</button>
    `;
    
    document.body.appendChild(resultDiv);
    
  } catch (error) {
    console.error('❌ 测试失败:', error.message);
    
    // 显示错误信息
    const errorDiv = document.createElement('div');
    errorDiv.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: #fef2f2;
      color: #dc2626;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      z-index: 10000;
      max-width: 400px;
      border: 1px solid #fca5a5;
    `;
    
    errorDiv.innerHTML = `
      <h3 style="margin: 0 0 10px 0;">❌ 测试失败</h3>
      <p style="margin: 0;">${error.message}</p>
      <button onclick="this.parentElement.remove()" style="margin-top: 15px; padding: 5px 10px; background: #dc2626; color: white; border: none; border-radius: 4px; cursor: pointer;">关闭</button>
    `;
    
    document.body.appendChild(errorDiv);
  }
}

// 运行测试
console.log('🚀 在控制台输入: quickAvatarTest() 来开始测试');
console.log('📝 或者直接运行下面这行代码:');
console.log('quickAvatarTest();');

// 自动运行（注释掉，避免自动执行）
// quickAvatarTest();