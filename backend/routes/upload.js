// routes/upload.js
const multer = require('multer');
const path = require('path');

// 配置存储
const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    cb(null, 'uploads/'); // 保存到uploads文件夹
  },
  filename: (req, file, cb) => {
    // 生成唯一文件名：时间戳+随机数+原扩展名
    const uniqueName = Date.now() + '-' + Math.round(Math.random() * 1E9) + path.extname(file.originalname);
    cb(null, uniqueName);
  }
});

const upload = multer({ 
  storage: storage,
  limits: {
    fileSize: 10 * 1024 * 1024 // 限制10MB
  },
  fileFilter: (req, file, cb) => {
    // 只允许图片和视频
    if (file.mimetype.startsWith('image/') || file.mimetype.startsWith('video/')) {
      cb(null, true);
    } else {
      cb(new Error('只支持图片和视频文件！'), false);
    }
  }
});

// 文件上传接口
router.post('/upload', authMiddleware, upload.single('file'), (req, res) => {
  if (!req.file) {
    return res.json({ code: 400, message: '请选择文件', data: null });
  }

  // 返回文件访问URL
  const fileUrl = `/uploads/${req.file.filename}`;
  
  res.json({
    code: 200,
    message: '上传成功',
    data: {
      url: fileUrl,
      filename: req.file.filename,
      size: req.file.size,
      mimetype: req.file.mimetype
    }
  });
});