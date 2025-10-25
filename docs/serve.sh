#!/bin/bash

# GMC 文档本地预览脚本

echo "======================================"
echo "   GMC 框架文档本地预览"
echo "======================================"
echo ""

# 检查是否安装了 docsify-cli
if command -v docsify &> /dev/null; then
    echo "✓ 检测到 docsify-cli"
    echo ""
    echo "启动文档服务器..."
    echo "访问地址："
    echo "  - 英文: http://localhost:3000"
    echo "  - 中文: http://localhost:3000/zh/"
    echo ""
    docsify serve .
elif command -v python3 &> /dev/null; then
    echo "✓ 检测到 Python 3"
    echo ""
    echo "启动文档服务器..."
    echo "访问地址："
    echo "  - 英文: http://localhost:3000"
    echo "  - 中文: http://localhost:3000/zh/"
    echo ""
    python3 -m http.server 3000
elif command -v python &> /dev/null; then
    echo "✓ 检测到 Python"
    echo ""
    echo "启动文档服务器..."
    echo "访问地址："
    echo "  - 英文: http://localhost:3000"
    echo "  - 中文: http://localhost:3000/zh/"
    echo ""
    python -m SimpleHTTPServer 3000
else
    echo "✗ 未找到可用的服务器工具"
    echo ""
    echo "请安装以下工具之一："
    echo "  1. docsify-cli: npm install -g docsify-cli"
    echo "  2. Python 3: https://www.python.org/downloads/"
    echo ""
    exit 1
fi
