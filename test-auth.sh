#!/bin/bash

# UpFTP 认证功能测试脚本

echo "🔐 UpFTP 认证功能测试"
echo "====================="
echo ""

# 设置颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数器
PASS=0
FAIL=0

# 测试函数
test_case() {
    local test_name="$1"
    local command="$2"
    local expected="$3"
    
    echo -n "测试: $test_name ... "
    
    result=$(eval "$command" 2>&1)
    
    if [[ "$result" == *"$expected"* ]]; then
        echo -e "${GREEN}✓ 通过${NC}"
        ((PASS++))
    else
        echo -e "${RED}✗ 失败${NC}"
        echo "  预期: $expected"
        echo "  实际: $result"
        ((FAIL++))
    fi
}

# 1. 测试帮助信息
echo ""
echo "📋 测试 1: 验证命令行参数"
echo "-------------------------"
test_case "HTTP认证参数存在" \
    "/tmp/upftp-auth/upftp-test -h" \
    "-http-auth"

test_case "HTTP用户名参数存在" \
    "/tmp/upftp-auth/upftp-test -h" \
    "-http-user"

test_case "HTTP密码参数存在" \
    "/tmp/upftp-auth/upftp-test -h" \
    "-http-pass"

# 2. 测试配置文件解析
echo ""
echo "📋 测试 2: 验证配置文件"
echo "-------------------------"
test_case "测试配置文件有效" \
    "cat /tmp/upftp-auth/test-auth.yaml" \
    "http_auth:"

# 3. 测试认证中间件
echo ""
echo "📋 测试 3: 验证认证中间件逻辑"
echo "-------------------------"

# 创建临时测试目录
mkdir -p /tmp/upftp-test-share
echo "test file" > /tmp/upftp-test-share/test.txt

# 启动服务器（后台运行）
echo ""
echo -e "${YELLOW}启动测试服务器（启用认证）...${NC}"
/tmp/upftp-auth/upftp-test \
    -p 19999 \
    -d /tmp/upftp-test-share \
    -http-auth \
    -http-user testuser \
    -http-pass testpass \
    -auto &
SERVER_PID=$!

# 等待服务器启动
sleep 2

# 4. 测试未认证访问
echo ""
echo "📋 测试 4: 测试未认证访问（应该被拒绝）"
echo "-------------------------"

test_case "主页未认证访问" \
    "curl -s -o /dev/null -w '%{http_code}' http://localhost:19999/" \
    "401"

test_case "API未认证访问" \
    "curl -s -o /dev/null -w '%{http_code}' http://localhost:19999/api/info" \
    "401"

test_case "文件未认证访问" \
    "curl -s -o /dev/null -w '%{http_code}' http://localhost:19999/files/test.txt" \
    "401"

# 5. 测试认证访问
echo ""
echo "📋 测试 5: 测试认证访问（应该成功）"
echo "-------------------------"

test_case "主页认证访问" \
    "curl -s -o /dev/null -w '%{http_code}' -u testuser:testpass http://localhost:19999/" \
    "200"

test_case "API认证访问" \
    "curl -s -u testuser:testpass http://localhost:19999/api/info" \
    "version"

test_case "文件认证访问" \
    "curl -s -u testuser:testpass http://localhost:19999/files/test.txt" \
    "test file"

# 6. 测试错误凭证
echo ""
echo "📋 测试 6: 测试错误凭证（应该被拒绝）"
echo "-------------------------"

test_case "错误用户名" \
    "curl -s -o /dev/null -w '%{http_code}' -u wronguser:testpass http://localhost:19999/" \
    "401"

test_case "错误密码" \
    "curl -s -o /dev/null -w '%{http_code}' -u testuser:wrongpass http://localhost:19999/" \
    "401"

# 停止服务器
echo ""
echo -e "${YELLOW}停止测试服务器...${NC}"
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

# 清理
rm -rf /tmp/upftp-test-share

# 输出结果
echo ""
echo "📊 测试结果"
echo "==========="
echo -e "通过: ${GREEN}$PASS${NC}"
echo -e "失败: ${RED}$FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✅ 所有测试通过！${NC}"
    echo ""
    echo "🎉 Phase 1 认证修复已完成并验证成功！"
    echo ""
    echo "📝 下一步:"
    echo "1. 提交代码到 Git"
    echo "2. 创建 Pull Request"
    echo "3. 更新文档"
    exit 0
else
    echo -e "${RED}❌ 部分测试失败${NC}"
    exit 1
fi
