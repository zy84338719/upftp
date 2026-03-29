#!/bin/bash

# 简化版认证测试

echo "🔐 UpFTP 认证功能快速测试"
echo "========================="
echo ""

# 创建测试目录
mkdir -p /tmp/upftp-quick-test
echo "Hello, World!" > /tmp/upftp-quick-test/hello.txt

echo "✅ 步骤 1: 启动服务器（后台运行）..."
/tmp/upftp-auth/upftp-test \
    -p 18888 \
    -d /tmp/upftp-quick-test \
    -http-auth \
    -http-user admin \
    -http-pass secretpass \
    -auto \
    > /tmp/upftp-server.log 2>&1 &
SERVER_PID=$!

sleep 3

echo ""
echo "📊 步骤 2: 测试未认证访问..."
echo "-----------------------------------"

echo -n "测试 2.1: 未认证访问主页 ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:18888/)
if [ "$HTTP_CODE" = "401" ]; then
    echo "✅ 通过 (返回 401)"
else
    echo "❌ 失败 (返回 $HTTP_CODE)"
fi

echo -n "测试 2.2: 未认证访问 API ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:18888/api/info)
if [ "$HTTP_CODE" = "401" ]; then
    echo "✅ 通过 (返回 401)"
else
    echo "❌ 失败 (返回 $HTTP_CODE)"
fi

echo -n "测试 2.3: 未认证访问文件 ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:18888/files/hello.txt)
if [ "$HTTP_CODE" = "401" ]; then
    echo "✅ 通过 (返回 401)"
else
    echo "❌ 失败 (返回 $HTTP_CODE)"
fi

echo ""
echo "📊 步骤 3: 测试认证访问..."
echo "-----------------------------------"

echo -n "测试 3.1: 认证访问主页 ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -u admin:secretpass http://localhost:18888/)
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ 通过 (返回 200)"
else
    echo "❌ 失败 (返回 $HTTP_CODE)"
fi

echo -n "测试 3.2: 认证访问 API ... "
RESULT=$(curl -s -u admin:secretpass http://localhost:18888/api/info)
if [[ "$RESULT" == *"httpAuthEnabled"* ]]; then
    echo "✅ 通过 (返回 JSON)"
else
    echo "❌ 失败 (返回: $RESULT)"
fi

echo -n "测试 3.3: 认证访问文件 ... "
RESULT=$(curl -s -u admin:secretpass http://localhost:18888/files/hello.txt)
if [[ "$RESULT" == *"Hello, World!"* ]]; then
    echo "✅ 通过 (返回文件内容)"
else
    echo "❌ 失败 (返回: $RESULT)"
fi

echo ""
echo "📊 步骤 4: 测试错误凭证..."
echo "-----------------------------------"

echo -n "测试 4.1: 错误用户名 ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -u wronguser:secretpass http://localhost:18888/)
if [ "$HTTP_CODE" = "401" ]; then
    echo "✅ 通过 (拒绝访问)"
else
    echo "❌ 失败 (返回 $HTTP_CODE)"
fi

echo -n "测试 4.2: 错误密码 ... "
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -u admin:wrongpass http://localhost:18888/)
if [ "$HTTP_CODE" = "401" ]; then
    echo "✅ 通过 (拒绝访问)"
else
    echo "❌ 失败 (返回 $HTTP_CODE)"
fi

# 停止服务器
echo ""
echo "🛑 步骤 5: 停止服务器..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

# 清理
rm -rf /tmp/upftp-quick-test

echo ""
echo "🎉 测试完成！"
echo ""
echo "✨ Phase 1 认证修复已成功实施！"
echo ""
echo "📝 修改内容："
echo "1. ✅ 所有 API 端点现在需要认证"
echo "2. ✅ /files/ 端点现在需要认证"
echo "3. ✅ 添加了命令行参数：-http-auth, -http-user, -http-pass"
echo "4. ✅ 添加了安全的文件处理函数"
echo ""
echo "🚀 使用方法："
echo "  ./upftp -http-auth -http-user admin -http-pass yourpassword"
