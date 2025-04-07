#!/bin/bash

# 要处理的目录列表
# DIRS=("syscalls" "execv" "subprocess")
DIRS=("sched" "execv")


# 输出文件名
OUTPUT_FILE="combined_files.txt"

# 清空或创建输出文件
> "$OUTPUT_FILE"

# 遍历目录列表
for dir in "${DIRS[@]}"; do
  # 检查目录是否存在并且是一个目录
  if [[ -d "$dir" ]]; then
    echo "Processing directory: $dir"
    # 使用 find 查找当前目录下的所有普通文件 (-type f)
    # -maxdepth 1 确保只查找当前目录, 不递归子目录
    # -print0 和 read -d '' -r 处理包含特殊字符的文件名
    find "$dir" -maxdepth 1 -type f -print0 | while IFS= read -r -d '' file; do
        echo "Adding file: $file"
        # 写入文件名作为分隔符/标记
        echo "########## Start of File: $file ##########" >> "$OUTPUT_FILE"
        # 追加文件内容
        cat "$file" >> "$OUTPUT_FILE"
        # 写入文件结束标记和换行, 增加可读性
        echo "" >> "$OUTPUT_FILE" # 添加一个空行
        echo "########## End of File: $file ##########" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE" # 再添加一个空行
    done
  else
    echo "Warning: Directory '$dir' not found or is not a directory. Skipping." >&2
  fi
done

echo "Done. All files combined into '$OUTPUT_FILE'"