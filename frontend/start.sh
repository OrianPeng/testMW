#!/bin/bash

echo "启动订单中心前端应用..."
echo

cd frontend

echo "检查Node.js环境..."
if ! command -v node &> /dev/null; then
    echo "错误: 未找到Node.js，请先安装Node.js 16+"
    exit 1
fi

node --version

echo
echo "安装依赖..."
npm install
if [ $? -ne 0 ]; then
    echo "错误: 依赖安装失败"
    exit 1
fi

echo
echo "启动开发服务器..."
echo "前端应用将在 http://localhost:3000 启动"
echo "按 Ctrl+C 停止服务器"
echo

npm run dev
