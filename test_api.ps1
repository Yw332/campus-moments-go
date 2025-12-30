# Campus Moments Go API 测试脚本
# 使用 PowerShell 运行: .\test_api.ps1

$API_BASE = "http://106.52.165.122:8080"

# 全局变量
$TOKEN = ""

function Test-API {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Data,
        [string]$Description
    )

    Write-Host "测试 $Description: " -NoNewline

    $headers = @{
        "Content-Type" = "application/json"
    }

    if ($TOKEN -ne "") {
        $headers["Authorization"] = "Bearer $TOKEN"
    }

    $fullUrl = "$API_BASE$Url"

    try {
        if ($Method -eq "GET") {
            $response = Invoke-RestMethod -Uri $fullUrl -Method GET -Headers $headers
        } elseif ($Method -eq "POST") {
            $response = Invoke-RestMethod -Uri $fullUrl -Method POST -Headers $headers -Body $Data
        } elseif ($Method -eq "PUT") {
            $response = Invoke-RestMethod -Uri $fullUrl -Method PUT -Headers $headers -Body $Data
        } elseif ($Method -eq "DELETE") {
            $response = Invoke-RestMethod -Uri $fullUrl -Method DELETE -Headers $headers
        }

        Write-Host "✅ 成功" -ForegroundColor Green
        Write-Host "   响应: $($response | ConvertTo-Json -Depth 3 -Compress)" -ForegroundColor Cyan
    } catch {
        Write-Host "❌ 失败" -ForegroundColor Red
        Write-Host "   错误: $($_.Exception.Message)" -ForegroundColor Red
    }

    Write-Host ""
}

# ==================== 测试开始 ====================
Write-Host "`n🚀 开始测试 Campus Moments Go API 接口` -ForegroundColor Yellow
Write-Host "==================================`n" -ForegroundColor Yellow

# 1. 系统接口
Write-Host "📋 1. 系统基础接口`n" -ForegroundColor Yellow
Test-API "GET" "/" "" "首页"
Test-API "GET" "/health" "" "健康检查"

# 2. 认证接口
Write-Host "`n📋 2. 认证接口`n" -ForegroundColor Yellow
Test-API "POST" "/auth/register" '{"username":"testuser009","phone":"13800138009","password":"Test123456"}' "用户注册"
Test-API "POST" "/auth/login" '{"account":"Yw166332","password":"JiangCan030"}' "用户登录"

# 保存 token（需要手动修改下面的 token）
# $TOKEN = "你的实际token"
Write-Host "`n提示: 请手动设置 `$TOKEN 变量，格式: `$TOKEN = `"eyJhbG...`"`n" -ForegroundColor Yellow

if ($TOKEN -eq "") {
    Write-Host "`n⚠️ 未设置Token，跳过需要认证的接口测试" -ForegroundColor Yellow
    Read-Host "`n按 Enter 键退出"
    exit
}

# 3. 用户接口
Write-Host "`n📋 3. 用户接口`n" -ForegroundColor Yellow
Test-API "GET" "/api/users/profile" "" "获取用户资料"

# 4. 帖子接口
Write-Host "`n📋 4. 帖子接口`n" -ForegroundColor Yellow
Test-API "GET" "/public/posts?page=1&pageSize=5" "" "公开帖子列表"
Test-API "GET" "/api/posts/my?page=1" "" "我的帖子"
Test-API "POST" "/api/posts" '{"title":"测试标题","content":"这是测试内容","visibility":0,"tags":["测试"]}' "创建帖子"

# 5. 评论接口
Write-Host "`n📋 5. 评论接口`n" -ForegroundColor Yellow
Test-API "GET" "/public/posts/1/comments" "" "评论列表"
Test-API "POST" "/api/comments/post/1" '{"content":"测试评论"}' "创建评论"

# 6. 点赞接口
Write-Host "`n📋 6. 点赞接口`n" -ForegroundColor Yellow
Test-API "POST" "/api/likes/post/1" "" "点赞帖子"
Test-API "GET" "/api/likes/posts/1" "" "帖子点赞列表"
Test-API "GET" "/api/likes/users?page=1" "" "我的点赞列表"

# 7. 好友接口
Write-Host "`n📋 7. 好友接口`n" -ForegroundColor Yellow
Test-API "GET" "/api/friends" "" "好友列表"
Test-API "GET" "/api/friends/requests" "" "好友请求列表"

# 8. 消息接口
Write-Host "`n📋 8. 消息接口`n" -ForegroundColor Yellow
Test-API "GET" "/api/conversations" "" "会话列表"
Test-API "GET" "/api/conversations/unread" "" "未读消息数"

# 9. 搜索接口
Write-Host "`n📋 9. 搜索接口`n" -ForegroundColor Yellow
Test-API "GET" "/search?keyword=测试" "" "搜索内容"
Test-API "GET" "/search/hot-words" "" "热门关键词"

# 10. 标签接口
Write-Host "`n📋 10. 标签接口`n" -ForegroundColor Yellow
Test-API "GET" "/public/tags" "" "标签列表"
Test-API "GET" "/public/tags/hot" "" "热门标签"

Write-Host "`n==================================" -ForegroundColor Yellow
Write-Host "🎯 测试完成！" -ForegroundColor Yellow
Read-Host "`n按 Enter 键退出"
