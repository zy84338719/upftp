#!/bin/bash

# 测试现代化页面

echo "🧪 测试现代化页面"
echo "===================="

# 创建测试目录
mkdir -p /tmp/upftp-modern-test
echo "Test file 1" > /tmp/upftp-modern-test/test1.txt
echo "Test file 2" > /tmp/upftp-modern-test/test2.pdf
echo "Test file 3" > /tmp/upftp-modern-test/test3.png

# 编译
echo ""
echo "📦 编译项目..."
cd /tmp/upftp-auth
go build -o upftp-modern

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译成功"

# 启动服务器
echo ""
echo "🚀 启动服务器..."
./upftp-modern \
    -p 16666 \
    -d /tmp/upftp-modern-test \
    -http-auth \
    -http-user admin \
    -http-pass test123 \
    -auto \
    > /tmp/upftp-modern.log 2>&1 &
SERVER_PID=$!

sleep 3

# 测试
echo ""
echo "🧪 测试现代化页面..."
echo "-------------------------"

echo "1. 测试访问根路径..."
HTTP_CODE=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:16666/)
if [ "$HTTP_CODE" = "401" ]; then
    echo "   ✅ 栜路径需要认证"
else
    echo "   ❌ 根路径不需要认证"
fi

echo ""
echo "2. 测试访问现代化页面..."
HTTP_CODE=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:16666/modern/)
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ 现代化页面可访问"
else
    echo "   ❌ 现代化页面不可访问 (HTTP $HTTP_CODE)"
fi

echo ""
echo "3. 测试文件列表..."
RESPONSE=$(curl -s http://localhost:16666/modern/)
if echo "$RESPONSE" | grep -q "test1.txt"; then
    echo "   ✅ 文件列表正常显示"
else
    echo "   ❌ 文件列表未正常显示"
fi

# 清理
echo ""
echo "🧹 清理..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null
rm -rf /tmp/upftp-modern-test

echo ""
echo "✅ 测试完成！"
echo ""
echo "📝 访问方式:"
echo "   - 旧版页面: http://localhost:16666/"
echo "   - 新版页面: http://localhost:16666/modern/"
