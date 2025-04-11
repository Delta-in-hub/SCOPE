#!/bin/bash

# 检查是否以 root 权限运行
if [ "$(id -u)" -ne 0 ]; then
    echo "请使用 sudo 运行此脚本"
    exit 1
fi

# 检查是否提供了APP参数
if [ $# -eq 0 ]; then
    echo "使用方法: $0 <APP>"
    echo "示例: $0 myapp"
    exit 1
fi

APP=$1
PREFIX="./bpf/build"

# 定义可执行文件列表（仅文件名，不含路径）
executables=(
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

# 存储所有子进程的PID
pids=()

# 捕获SIGINT信号（Ctrl+C），并杀死所有子进程
trap 'echo "正在终止所有子进程..."; kill ${pids[@]} 2>/dev/null; exit 1' SIGINT

# 遍历并运行每个可执行文件
for exe in "${executables[@]}"; do
    full_path="${PREFIX}/${exe}"
    
    if [ -x "$full_path" ]; then
        echo "正在运行: $full_path -c $APP"
        "$full_path" -c "$APP" &  # 后台运行并传递APP参数
        pids+=($!) # 存储PID
    else
        echo "警告: $full_path 不存在或不可执行"
    fi
done

echo "所有BPF程序已在后台运行（使用APP参数: $APP）"
echo "按 Ctrl+C 停止所有程序"

# 等待所有子进程结束（或等待Ctrl+C）
wait

echo "所有BPF程序已停止"