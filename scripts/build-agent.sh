#!/bin/bash

# 启用错误检测和报错
set -e

# 检查是否为 Debian/Ubuntu 系统
if [ -f "/etc/debian_version" ]; then
    echo "检测到 Debian/Ubuntu 系统，正在安装依赖..."
    sudo apt-get update
    sudo apt-get install -y \
        git \
        build-essential \
        clang \
        llvm \
        libelf-dev \
        libbpf-dev \
        libzmq3-dev \
        libmsgpack-dev
    if [ $? -ne 0 ]; then
        echo "错误：依赖安装失败"
        exit 1
    fi
else
    echo "注意：非 Debian/Ubuntu 系统，请确保已安装所有必要依赖"
fi

# 检查是否在根目录下
if [ ! -f "go.mod" ] || [ ! -d "bpf" ] || [ ! -d "cmd" ]; then
    echo "错误：必须在项目根目录下执行此脚本"
    exit 1
fi

# 拉取 git submodules
echo "正在更新 git submodules..."
git submodule update --init --recursive
if [ $? -ne 0 ]; then
    echo "错误：git submodule 更新失败"
    exit 1
fi

# 进入 bpf 目录
cd bpf || {
    echo "错误：无法进入 bpf 目录"
    exit 1
}

# 定义要构建的
apps=(
    "cuda"
    "execv"
    "ggml_base"
    "ggml_cpu"
    "ggml_cuda"
    "Ollamabin"
    "sched"
    "syscalls"
    "vfs_open"
)

# 遍历并构建每个应用
for app in "${apps[@]}"; do
    echo "正在构建 $app..."
    make APP="$app" -j$(nproc)
    if [ $? -ne 0 ]; then
        echo "错误：构建 $app 失败"
        exit 1
    fi
done

# 返回项目根目录
cd .. || {
    echo "错误：无法返回项目根目录"
    exit 1
}

# 初始化 Go 模块
echo "正在初始化 Go 模块..."
go mod download
if [ $? -ne 0 ]; then
    echo "错误：go mod download 失败"
    exit 1
fi

# 构建 scope-agent-manager
echo "正在构建 scope-agent-manager..."
go build ./cmd/scope-agent-manager
if [ $? -ne 0 ]; then
    echo "错误：构建 scope-agent-manager 失败"
    exit 1
fi

echo "所有构建步骤完成！"