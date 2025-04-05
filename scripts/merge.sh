# 定义输出文件名
output_file="combined_code.txt"

# 清空或创建输出文件
> "$output_file"

# 查找当前目录下的 .c 和 .h 文件，并处理
# -maxdepth 1: 只查找当前目录，不进入子目录
# -type f: 只查找普通文件
# \( ... \): 分组条件
# -name "*.c" -o -name "*.h": 查找名字以 .c 或 .h 结尾的文件
# -exec sh -c '...': 对找到的每个文件执行一个 shell 片段
#   echo "--- File: $1 ---": 打印文件名作为分隔符
#   cat "$1": 打印文件内容
#   echo "": 打印一个空行增加可读性
# sh {}: 将找到的文件名 ($1) 传递给 sh -c
# \;: 表示 -exec 命令的结束，为每个文件单独执行 sh -c
# >> "$output_file": 将 sh -c 的输出追加到目标文件
find . -maxdepth 1 -type f \( -name "*.c" -o -name "*.h" \) \
    -exec sh -c '
        echo "--- File: $1 ---"
        cat "$1"
        echo "" # 在文件内容后添加一个空行，可选
    ' sh {} \; >> "$output_file"

echo "所有 .c 和 .h 文件已合并到 $output_file"
