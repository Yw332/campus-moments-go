# 📸 头像上传 API 示例

## 🚀 快速开始

### 1. 先登录获取Token
```bash
curl -X POST "http://106.52.165.122:8080/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "account": "Yw166332",
    "password": "JiangCan030"
  }'
```

**复制返回的 token，后续请求需要使用**

### 2. 上传头像（使用curl）
```bash
curl -X POST "http://106.52.165.122:8080/api/upload/avatar" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "avatar=@/path/to/your/avatar.jpg"
```

### 3. 获取用户资料查看头像
```bash
curl -X GET "http://106.52.165.122:8080/api/users/profile" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 📱 前端JavaScript示例

### 完整的上传头像流程
```javascript
// 1. 登录函数
async function login() {
  const response = await fetch('http://106.52.165.122:8080/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      account: 'Yw166332',
      password: 'JiangCan030'
    })
  });
  
  const result = await response.json();
  if (result.code === 200) {
    localStorage.setItem('token', result.data.token);
    console.log('登录成功，token:', result.data.token);
    return result.data.token;
  }
  throw new Error('登录失败');
}

// 2. 上传头像函数
async function uploadAvatar(file) {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('请先登录');
  }

  const formData = new FormData();
  formData.append('avatar', file);

  const response = await fetch('http://106.52.165.122:8080/api/upload/avatar', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  });

  const result = await response.json();
  if (result.code === 200) {
    console.log('头像上传成功:', result.data.avatarUrl);
    return result.data.avatarUrl;
  } else {
    throw new Error(result.message || '上传失败');
  }
}

// 3. 获取用户资料函数
async function getUserProfile() {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('请先登录');
  }

  const response = await fetch('http://106.52.165.122:8080/api/users/profile', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  const result = await response.json();
  if (result.code === 200) {
    console.log('用户资料:', result.data);
    return result.data;
  } else {
    throw new Error(result.message || '获取用户资料失败');
  }
}

// 4. 完整测试流程
async function testAvatarUpload() {
  try {
    // 登录
    const token = await login();
    
    // 模拟文件选择（实际使用中用户通过input选择）
    const testFile = new File(['test image content'], 'test.jpg', { type: 'image/jpeg' });
    
    // 上传头像
    const avatarUrl = await uploadAvatar(testFile);
    
    // 获取用户资料验证头像
    const userProfile = await getUserProfile();
    
    console.log('✅ 测试完成');
    console.log('头像URL:', avatarUrl);
    console.log('用户资料头像:', userProfile.avatar);
    
    // 验证头像是否一致
    if (userProfile.avatar === avatarUrl) {
      console.log('✅ 头像上传和显示验证成功');
    } else {
      console.log('❌ 头像显示不一致');
    }
    
  } catch (error) {
    console.error('❌ 测试失败:', error.message);
  }
}

// 在浏览器控制台运行：testAvatarUpload()
```

### HTML 表单示例
```html
<!DOCTYPE html>
<html>
<head>
    <title>头像上传示例</title>
</head>
<body>
    <h3>头像上传测试</h3>
    
    <!-- 登录部分 -->
    <div>
        <h4>1. 登录</h4>
        <button onclick="quickLogin()">快速登录测试账号</button>
    </div>
    
    <!-- 头像上传部分 -->
    <div>
        <h4>2. 选择头像文件</h4>
        <input type="file" id="avatarInput" accept="image/*" onchange="previewAvatar(event)">
        <div id="preview" style="width: 100px; height: 100px; border: 1px solid #ccc; margin: 10px 0;">
            <!-- 预览区域 -->
        </div>
        <button onclick="uploadSelectedAvatar()">上传头像</button>
    </div>
    
    <!-- 查看用户资料 -->
    <div>
        <h4>3. 查看用户资料</h4>
        <button onclick="showUserProfile()">获取用户资料</button>
        <div id="userProfile"></div>
    </div>

    <script>
        const API_BASE = 'http://106.52.165.122:8080';
        let currentToken = '';

        async function quickLogin() {
            try {
                const response = await fetch(API_BASE + '/auth/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        account: 'Yw166332',
                        password: 'JiangCan030'
                    })
                });
                
                const result = await response.json();
                if (result.code === 200) {
                    currentToken = result.data.token;
                    console.log('登录成功，token已保存');
                    alert('登录成功！');
                } else {
                    alert('登录失败：' + result.message);
                }
            } catch (error) {
                alert('登录失败：' + error.message);
            }
        }

        function previewAvatar(event) {
            const file = event.target.files[0];
            if (file) {
                const reader = new FileReader();
                reader.onload = function(e) {
                    document.getElementById('preview').innerHTML = 
                        `<img src="${e.target.result}" style="width: 100%; height: 100%; object-fit: cover;">`;
                };
                reader.readAsDataURL(file);
            }
        }

        async function uploadSelectedAvatar() {
            const fileInput = document.getElementById('avatarInput');
            const file = fileInput.files[0];
            
            if (!file) {
                alert('请选择头像文件');
                return;
            }
            
            if (!currentToken) {
                alert('请先登录');
                return;
            }

            try {
                const formData = new FormData();
                formData.append('avatar', file);

                const response = await fetch(API_BASE + '/api/upload/avatar', {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${currentToken}`
                    },
                    body: formData
                });

                const result = await response.json();
                if (result.code === 200) {
                    alert('头像上传成功！');
                    console.log('头像URL:', result.data.avatarUrl);
                } else {
                    alert('上传失败：' + result.message);
                }
            } catch (error) {
                alert('上传失败：' + error.message);
            }
        }

        async function showUserProfile() {
            if (!currentToken) {
                document.getElementById('userProfile').innerHTML = '请先登录';
                return;
            }

            try {
                const response = await fetch(API_BASE + '/api/users/profile', {
                    headers: {
                        'Authorization': `Bearer ${currentToken}`
                    }
                });

                const result = await response.json();
                if (result.code === 200) {
                    const user = result.data;
                    const avatarHtml = user.avatar ? 
                        `<img src="${user.avatar}" style="width: 50px; height: 50px; border-radius: 50%;">` : 
                        '无头像';
                    
                    document.getElementById('userProfile').innerHTML = `
                        <div style="padding: 10px; background: #f5f5f5;">
                            <p><strong>用户名:</strong> ${user.username}</p>
                            <p><strong>头像:</strong> ${avatarHtml}</p>
                            ${user.avatar ? `<p><a href="${user.avatar}" target="_blank">查看头像</a></p>` : ''}
                        </div>
                    `;
                } else {
                    document.getElementById('userProfile').innerHTML = '获取用户资料失败：' + result.message;
                }
            } catch (error) {
                document.getElementById('userProfile').innerHTML = '获取用户资料失败：' + error.message;
            }
        }
    </script>
</body>
</html>
```

## 🧪 Postman 测试示例

### 1. 环境变量设置
```
base_url = http://106.52.165.122:8080
token = {{登录响应中复制的token}}
```

### 2. 登录请求
- **Method**: POST
- **URL**: `{{base_url}}/auth/login`
- **Headers**: `Content-Type: application/json`
- **Body**: 
```json
{
  "account": "Yw166332",
  "password": "JiangCan030"
}
```

### 3. 上传头像请求
- **Method**: POST
- **URL**: `{{base_url}}/api/upload/avatar`
- **Headers**: `Authorization: Bearer {{token}}`
- **Body**: 
  - Type: `form-data`
  - Key: `avatar` (Type: File)
  - Value: 选择你的头像文件

### 4. 获取用户资料
- **Method**: GET
- **URL**: `{{base_url}}/api/users/profile`
- **Headers**: `Authorization: Bearer {{token}}`

## 🔍 响应示例

### 成功上传头像
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatarUrl": "http://106.52.165.122:8080/uploads/avatars/1_20231214150000_a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6.jpg",
    "userId": 1
  }
}
```

### 获取用户资料（含头像）
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "userId": 1,
    "username": "Yw166332",
    "phone": "17875242005",
    "avatar": "http://106.52.165.122:8080/uploads/avatars/1_20231214150000_a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6.jpg",
    "status": 1,
    "createdAt": "2023-12-14T15:00:00Z",
    "updatedAt": "2023-12-14T15:00:00Z"
  }
}
```

---

**📞 快速测试推荐**: 使用 `frontend-test.html` 文件进行可视化测试！