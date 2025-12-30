# Campus Moments Go API - Windows 测试指南

## 运行 PowerShell 测试脚本

### 1. 首先获取 Token
先运行一次登录测试，从响应中复制 token：

```powershell
# 运行测试脚本
.\test_api.ps1
```

从登录响应中复制 `token` 字段的值。

### 2. 设置 Token 并运行测试
打开 `test_api.ps1` 文件，找到这一行：

```powershell
# $TOKEN = "你的实际token"
```

修改为：

```powershell
$TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

然后再次运行脚本：

```powershell
.\test_api.ps1
```

---

## 单独测试各个接口

### 1. 用户登录

```powershell
$headers = @{"Content-Type" = "application/json"}
$body = @{
    account = "Yw166332"
    password = "JiangCan030"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://106.52.165.122:8080/auth/login" -Method POST -Headers $headers -Body $body
$response | ConvertTo-Json -Depth 10
```

### 2. 获取用户资料

```powershell
$TOKEN = "你的token"
$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

$response = Invoke-RestMethod -Uri "http://106.52.165.122:8080/api/users/profile" -Method GET -Headers $headers
$response | ConvertTo-Json -Depth 10
```

### 3. 获取公开帖子列表

```powershell
$response = Invoke-RestMethod -Uri "http://106.52.165.122:8080/public/posts?page=1&pageSize=10" -Method GET
$response | ConvertTo-Json -Depth 10
```

### 4. 创建帖子

```powershell
$TOKEN = "你的token"
$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

$body = @{
    title = "测试标题"
    content = "这是一条测试帖子"
    visibility = 0
    tags = @("测试", "API")
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://106.52.165.122:8080/api/posts" -Method POST -Headers $headers -Body $body
$response | ConvertTo-Json -Depth 10
```

### 5. 上传头像

```powershell
$TOKEN = "你的token"
$headers = @{
    "Authorization" = "Bearer $TOKEN"
}

$filePath = "C:\path\to\your\avatar.jpg"
$response = Invoke-RestMethod -Uri "http://106.52.165.122:8080/api/upload/avatar" -Method POST -Headers $headers -Form @{
    avatar = Get-Item -Path $filePath
}

$response | ConvertTo-Json -Depth 10
```

### 6. 获取好友列表

```powershell
$TOKEN = "你的token"
$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

$response = Invoke-RestMethod -Uri "http://106.52.165.122:8080/api/friends?page=1" -Method GET -Headers $headers
$response | ConvertTo-Json -Depth 10
```

---

## 常见问题

### PowerShell 执行策略错误
如果提示脚本执行策略错误，运行：

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 查看响应详细信息

```powershell
$response = Invoke-RestMethod -Uri "http://106.52.165.122:8080/health" -Method GET

# 查看详细响应
Write-Host "Code: $($response.code)"
Write-Host "Message: $($response.message)"
Write-Host "Data: $($response.data | ConvertTo-Json)"
```

---

## 推荐工具

除了 PowerShell，也可以使用以下工具：

1. **Postman** - 图形化 API 测试工具
2. **Apifox** - 支持导入 `apifox-collection.json`
3. **Bruno** - 轻量级 API 测试工具
4. **curl (WSL)** - 在 Windows Subsystem for Linux 中使用
