const express = require('express');
const cors = require('cors');
const path = require('path');

const app = express();
const PORT = 3000;

// 中间件
app.use(cors());
app.use(express.json());
app.use('/uploads', express.static(path.join(__dirname, 'uploads')));

// 模拟数据（临时使用，后续可连接数据库）
const users = [
  { id: 1, username: 'admin', password: '123456' }
];
const posts = [
  { 
    id: 1, 
    content: '欢迎使用校园时刻！', 
    userId: 1, 
    username: 'admin',
    createTime: new Date().toISOString()
  }
];

// 路由
app.use('/api/auth', require('./routes/auth')(users));
app.use('/api/posts', require('./routes/posts')(posts));

// 健康检查接口
app.get('/api/hello', (req, res) => {
  res.json({ 
    code: 200, 
    message: '后端服务正常运行！',
    data: { service: '校园时刻后端', version: '1.0' }
  });
});

// 启动服务器
app.listen(PORT, () => {
  console.log('🎉 后端服务器启动成功！');
  console.log(`📍 访问地址: http://localhost:${PORT}`);
  console.log('✅ 测试接口: http://localhost:3000/api/hello');
});