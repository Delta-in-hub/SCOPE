#!/bin/bash

# 启用严格错误检查
set -euo pipefail

# 检查是否在项目根目录
if [ ! -f "go.mod" ] || [ ! -d "web" ]; then
    echo "错误：必须在项目根目录下执行此脚本"
    exit 1
fi

# 检查 web/web 目录是否存在
if [ ! -d "web/web" ]; then
    echo "错误：web/web 目录不存在"
    exit 1
fi

# 进入 web/web 目录
echo "正在进入 web/web 目录..."
cd web/web || {
    echo "错误：无法进入 web/web 目录"
    exit 1
}

# 检查 pnpm 是否安装
if ! command -v pnpm &> /dev/null; then
    echo "错误：pnpm 未安装，请先安装 pnpm"
    exit 1
fi

# 检查 package.json 是否存在
if [ ! -f "package.json" ]; then
    echo "错误：package.json 文件不存在"
    exit 1
fi

# 安装依赖（可选，可根据需要取消注释）
echo "正在安装依赖..."
pnpm install
if [ $? -ne 0 ]; then
    echo "错误：依赖安装失败"
    exit 1
fi

# 执行构建
echo "正在执行 pnpm build..."
pnpm build
if [ $? -ne 0 ]; then
    echo "错误：构建失败"
    exit 1
fi

echo "前端构建成功完成！"