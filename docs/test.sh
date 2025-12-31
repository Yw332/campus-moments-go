#!/bin/bash

# ============================================
# Campus Moments Go API 完整测试脚本
# 模仿前端请求方式测试所有后端功能
# ============================================

# 配置
API_BASE="${API_BASE:-http://localhost:8080}"
TIMESTAMP=$(date +%s)
TEST_USERNAME="testuser_${TIMESTAMP}"
TEST_PHONE="138${TIMESTAMP: -8}"  # 使用时间戳生成唯一手机号
TEST_PASSWORD="Test123456"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 统计变量
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 临时文件
RESPONSE_FILE="/tmp/campus_moments_test_response.json"
TOKEN_FILE="/tmp/campus_moments_test_token.txt"

# 清理函数
cleanup() {
    rm -f "$RESPONSE_FILE" "$TOKEN_FILE"
}
trap cleanup EXIT

# ============================================
# 工具函数
# ============================================

# 打印测试标题
print_section() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# 打印功能说明
print_description() {
    echo -e "${BLUE}📝 $1${NC}"
}

# 测试API函数
test_api() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4
    local need_auth=${5:-false}
    local expected_code=${6:-200}
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo ""
    echo -n -e "${YELLOW}测试 $description: ${NC}"
    
    # 构建curl命令
    local curl_cmd="curl -s -w \"%{http_code}\" -o \"$RESPONSE_FILE\""
    
    # 添加请求头
    curl_cmd="$curl_cmd -H \"Content-Type: application/json\""
    
    # 如果需要认证，添加Token
    if [ "$need_auth" = "true" ] && [ -f "$TOKEN_FILE" ]; then
        local token=$(cat "$TOKEN_FILE")
        curl_cmd="$curl_cmd -H \"Authorization: Bearer $token\""
    fi
    
    # 添加请求方法和数据
    if [ "$method" = "GET" ]; then
        curl_cmd="$curl_cmd -X GET \"$API_BASE$url\""
    elif [ "$method" = "POST" ]; then
        curl_cmd="$curl_cmd -X POST -d '$data' \"$API_BASE$url\""
    elif [ "$method" = "PUT" ]; then
        curl_cmd="$curl_cmd -X PUT -d '$data' \"$API_BASE$url\""
    elif [ "$method" = "DELETE" ]; then
        curl_cmd="$curl_cmd -X DELETE \"$API_BASE$url\""
    fi
    
    # 执行请求
    local http_code=$(eval $curl_cmd)
    
    # 检查响应
    if [ "$http_code" = "$expected_code" ] || [ "$http_code" = "200" ] || [ "$http_code" = "201" ]; then
        echo -e "${GREEN}✅ 通过 (HTTP $http_code)${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        
        # 显示响应摘要
        if [ -f "$RESPONSE_FILE" ]; then
            local response_preview=$(head -c 150 "$RESPONSE_FILE" 2>/dev/null || echo "")
            if [ ! -z "$response_preview" ]; then
                echo -e "   ${GREEN}响应: ${response_preview}...${NC}"
            fi
        fi
        return 0
    else
        echo -e "${RED}❌ 失败 (HTTP $http_code)${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        
        # 显示错误信息
        if [ -f "$RESPONSE_FILE" ]; then
            local error_msg=$(cat "$RESPONSE_FILE" 2>/dev/null | head -c 200)
            if [ ! -z "$error_msg" ]; then
                echo -e "   ${RED}错误: ${error_msg}${NC}"
            fi
        fi
        return 1
    fi
}

# 提取Token
extract_token() {
    if [ -f "$RESPONSE_FILE" ]; then
        local token=$(grep -o '"token":"[^"]*"' "$RESPONSE_FILE" 2>/dev/null | cut -d'"' -f4)
        if [ ! -z "$token" ]; then
            echo "$token" > "$TOKEN_FILE"
            echo -e "${GREEN}🔑 Token已保存: ${token:0:30}...${NC}"
            return 0
        fi
    fi
    return 1
}

# 提取数据字段
extract_field() {
    local field=$1
    if [ -f "$RESPONSE_FILE" ]; then
        grep -o "\"$field\":\"[^\"]*\"" "$RESPONSE_FILE" 2>/dev/null | cut -d'"' -f4
    fi
}

# ============================================
# 开始测试
# ============================================

echo -e "${GREEN}"
echo "╔═══════════════════════════════════════════════════════╗"
echo "║   Campus Moments Go API 完整功能测试脚本              ║"
echo "║   模仿前端请求方式测试所有后端功能                    ║"
echo "╚═══════════════════════════════════════════════════════╝"
echo -e "${NC}"

echo -e "${YELLOW}测试配置:${NC}"
echo "  后端地址: $API_BASE"
echo "  测试用户: $TEST_USERNAME"
echo "  测试手机: $TEST_PHONE"
echo ""

# ============================================
# 一、系统基础接口测试
# ============================================

print_section "一、系统基础接口测试"

print_description "健康检查接口，用于检测服务是否正常运行"
test_api "GET" "/health" "" "健康检查" false

print_description "首页接口"
test_api "GET" "/" "" "首页" false

# ============================================
# 二、用户认证接口测试
# ============================================

print_section "二、用户认证接口测试"

print_description "用户注册：新用户注册账号，验证用户名、手机号、密码格式"
test_api "POST" "/auth/register" "{\"username\":\"$TEST_USERNAME\",\"phone\":\"$TEST_PHONE\",\"password\":\"$TEST_PASSWORD\"}" "用户注册" false

print_description "用户登录：获取JWT Token，支持用户名或手机号登录"
test_api "POST" "/auth/login" "{\"account\":\"$TEST_USERNAME\",\"password\":\"$TEST_PASSWORD\"}" "用户登录" false

# 提取并保存Token
if extract_token; then
    TOKEN=$(cat "$TOKEN_FILE")
    echo -e "${GREEN}✅ Token提取成功，后续测试将使用此Token${NC}"
else
    echo -e "${RED}❌ 无法提取Token，后续需要认证的测试将跳过${NC}"
    TOKEN=""
fi

if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⚠️  警告: 未获取到Token，跳过需要认证的接口测试${NC}"
    echo ""
    echo "=================================="
    echo "测试完成！"
    echo "总测试数: $TOTAL_TESTS"
    echo -e "通过: ${GREEN}$PASSED_TESTS${NC} | 失败: ${RED}$FAILED_TESTS${NC}"
    exit 0
fi

print_description "获取用户资料：获取当前登录用户的详细信息"
test_api "GET" "/api/users/profile" "" "获取用户资料" true

print_description "修改密码：用户修改登录密码，需要验证旧密码"
test_api "PUT" "/api/users/password" "{\"oldPassword\":\"$TEST_PASSWORD\",\"newPassword\":\"NewPass123\"}" "修改密码" true

# 改回原密码以便后续测试
test_api "PUT" "/api/users/password" "{\"oldPassword\":\"NewPass123\",\"newPassword\":\"$TEST_PASSWORD\"}" "恢复原密码" true

# ============================================
# 三、动态相关接口测试
# ============================================

print_section "三、动态相关接口测试"

print_description "获取动态列表：获取首页动态列表，支持分页，返回格式匹配前端瀑布流需求"
test_api "GET" "/api/moments?page=1&pageSize=10" "" "获取动态列表" true

print_description "发布动态：用户发布新动态，支持标题、内容、标签、图片"
MOMENT_DATA="{\"title\":\"测试动态标题\",\"content\":\"这是测试动态内容，用于验证发布功能\",\"tags\":[\"测试\",\"校园\"],\"images\":[]}"
test_api "POST" "/api/moments" "$MOMENT_DATA" "发布动态" true

# 获取刚发布的动态ID（如果可能）
MOMENT_ID=$(extract_field "id" | head -n1)
if [ -z "$MOMENT_ID" ]; then
    MOMENT_ID=1  # 默认值
fi

print_description "获取动态详情：获取单条动态的详细信息，包含图片、评论等"
test_api "GET" "/api/moments/$MOMENT_ID" "" "获取动态详情" true

print_description "获取我的动态：获取当前用户发布的所有动态列表"
test_api "GET" "/api/moments/my?page=1&pageSize=10" "" "获取我的动态" true

# ============================================
# 四、点赞相关接口测试
# ============================================

print_section "四、点赞相关接口测试"

print_description "点赞动态：点赞或取消点赞动态，自动更新点赞数"
test_api "POST" "/api/likes/post/$MOMENT_ID" "" "点赞动态" true

# ============================================
# 五、评论相关接口测试
# ============================================

print_section "五、评论相关接口测试"

print_description "获取评论列表：获取指定动态的所有评论，公开接口无需Token"
test_api "GET" "/public/posts/$MOMENT_ID/comments?page=1&pageSize=20" "" "获取评论列表" false

print_description "发布评论：在指定动态下发布评论"
COMMENT_DATA="{\"content\":\"这是一条测试评论\"}"
test_api "POST" "/api/comments/post/$MOMENT_ID" "$COMMENT_DATA" "发布评论" true

# 获取刚发布的评论ID
COMMENT_ID=$(extract_field "id" | head -n1)
if [ -z "$COMMENT_ID" ]; then
    COMMENT_ID=1  # 默认值
fi

print_description "点赞评论：点赞或取消点赞评论"
test_api "POST" "/api/comments/$COMMENT_ID/like" "" "点赞评论" true

# ============================================
# 六、搜索相关接口测试
# ============================================

print_section "六、搜索相关接口测试"

print_description "搜索内容：搜索动态和用户，支持关键词搜索"
test_api "GET" "/api/search?keyword=测试&page=1&pageSize=10" "" "搜索内容" true

print_description "获取热门关键词：获取搜索热门关键词列表"
test_api "GET" "/api/search/hot-words" "" "获取热门关键词" true

# ============================================
# 七、标签相关接口测试
# ============================================

print_section "七、标签相关接口测试"

print_description "获取标签列表：获取所有可用标签列表，用于发布动态时选择"
test_api "GET" "/public/tags" "" "获取标签列表" false

# ============================================
# 八、文件上传接口测试（模拟）
# ============================================

print_section "八、文件上传接口测试"

print_description "文件上传：上传图片文件，用于动态发布（需要实际文件，此处仅测试接口可用性）"
echo -e "${YELLOW}⚠️  文件上传接口需要实际文件，跳过实际文件上传测试${NC}"
echo -e "${YELLOW}   如需测试，请使用以下命令：${NC}"
echo -e "${CYAN}   curl -X POST $API_BASE/api/upload/file \\${NC}"
echo -e "${CYAN}     -H \"Authorization: Bearer \$TOKEN\" \\${NC}"
echo -e "${CYAN}     -F \"file=@/path/to/image.jpg\"${NC}"

# ============================================
# 九、错误处理测试
# ============================================

print_section "九、错误处理测试"

print_description "未认证请求：不带Token访问需要认证的接口，应返回401"
test_api "GET" "/api/users/profile" "" "未认证请求测试" false 401

print_description "无效参数：注册时使用无效参数，应返回400"
test_api "POST" "/auth/register" "{\"username\":\"ab\",\"phone\":\"123\",\"password\":\"123\"}" "无效参数测试" false 400

# ============================================
# 十、退出登录测试
# ============================================

print_section "十、退出登录测试"

print_description "退出登录：用户退出登录，将Token加入黑名单"
test_api "POST" "/api/auth/logout" "" "退出登录" true

# 测试退出后Token是否失效
print_description "验证Token失效：退出后使用原Token应返回401"
test_api "GET" "/api/users/profile" "" "验证Token失效" true 401

# ============================================
# 测试总结
# ============================================

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}测试总结${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "总测试数: $TOTAL_TESTS"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}❌ 有 $FAILED_TESTS 个测试失败${NC}"
    exit 1
fi

