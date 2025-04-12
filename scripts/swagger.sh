#!/bin/bash

# 安全设置
set -euo pipefail  # 启用严格模式：错误退出、未定义变量报错、管道错误检测

# 查找swag可执行文件
SWAG_PATH=""
if command -v swag &>/dev/null; then
    SWAG_PATH=$(command -v swag)
elif [ -f "$HOME/go/bin/swag" ]; then
    SWAG_PATH="$HOME/go/bin/swag"
else
    echo "错误: 未找到swag可执行文件" >&2
    exit 1
fi

# 检查main.go是否存在
MAIN_GO_PATH="./cmd/scope-backend/main.go"
if [ ! -f "$MAIN_GO_PATH" ]; then
    echo "错误: 未找到 $MAIN_GO_PATH" >&2
    exit 1
fi

# 执行命令
echo "使用swag路径: $SWAG_PATH"
echo "处理文件: $MAIN_GO_PATH"
"$SWAG_PATH" init --parseDependency --parseInternal -g "$MAIN_GO_PATH" -o docs/backend

# 检查执行结果
if [ $? -ne 0 ]; then
    echo "错误: swag命令执行失败" >&2
    exit 1
fi

echo "swag文档生成成功"