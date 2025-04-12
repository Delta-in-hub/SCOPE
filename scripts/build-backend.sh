#!/bin/bash

# 启用错误检测和报错
set -e

# 检查是否在根目录下
if [ ! -f "go.mod" ] || [ ! -d "bpf" ] || [ ! -d "cmd" ]; then
    echo "错误：必须在项目根目录下执行此脚本"
    exit 1
fi

# 初始化 Go 模块
echo "正在初始化 Go 模块..."
go mod download
if [ $? -ne 0 ]; then
    echo "错误：go mod download 失败"
    exit 1
fi

# 构建 scope-backend
echo "正在构建 scope-backend..."
go build ./cmd/scope-backend
if [ $? -ne 0 ]; then
    echo "错误：构建 scope-backend 失败"
    exit 1
fi

echo "所有构建步骤完成！"