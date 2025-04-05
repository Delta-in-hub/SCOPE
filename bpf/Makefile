# 目标应用程序名称 (可以被命令行覆盖: make APP=...)
APP ?= openat2ring


# 输出目录，存放所有编译产物
OUTPUT := ./build/libbpf
# BPF 源文件
APP_BPF_SRCS := $(APP).bpf.c
# 用户空间 C 源文件
APP_C_SRCS := $(APP).c

# ---- 工具链定义 ----

# Clang 编译器，用于编译 BPF 代码，如果环境中定义了 CLANG，则使用环境的
CLANG ?= clang
# C 编译器，用于编译用户空间代码，如果环境中定义了 CC，则使用环境的
# $(call allow-override...) 用于确保命令行或环境变量可以覆盖默认值 'cc'
# 这里我们直接使用 ?= gcc，更简洁常用
CC ?= clang
# 链接器，通常由 CC 驱动，无需单独指定 LD 除非有特殊需求
# LD ?= ld
# BPFTOOL 工具路径，指向我们将从源码编译的版本
BPFTOOL ?= $(abspath $(OUTPUT)/bpftool/bootstrap/bpftool)
# AR 工具，用于创建静态库 (libbpf.a)
AR ?= ar
# RM 工具，用于清理
RM ?= rm -f

# ---- 源码路径 ----

# libbpf 源码路径 (相对于当前 Makefile)
LIBBPF_SRC := $(abspath ./libbpf-bootstrap/libbpf/src)
# bpftool 源码路径 (相对于当前 Makefile)
BPFTOOL_SRC := $(abspath ./libbpf-bootstrap/bpftool/src)
# 预生成的 vmlinux.h 头文件路径
# 检测系统架构，并指向 libbpf-bootstrap/vmlinux.h/ 下对应的头文件
ARCH ?= $(shell uname -m | sed 's/x86_64/x86/' \
             | sed 's/arm.*/arm/' \
             | sed 's/aarch64/arm64/' \
             | sed 's/ppc64le/powerpc/' \
             | sed 's/mips.*/mips/' \
             | sed 's/riscv64/riscv/' \
             | sed 's/loongarch64/loongarch/')
VMLINUX_H_DIR := $(abspath ./libbpf-bootstrap/vmlinux.h/include/$(ARCH))
VMLINUX := $(VMLINUX_H_DIR)/vmlinux.h

# ---- 编译产物路径 ----

# libbpf 静态库路径
LIBBPF_OBJ := $(abspath $(OUTPUT)/libbpf.a)
# bpftool 构建目录
BPFTOOL_OUTPUT := $(abspath $(OUTPUT)/bpftool)



# BPF 目标文件路径
APP_BPF_OBJ := ./build/$(APP).bpf.o
# BPF 骨架头文件路径
APP_SKEL := ./build/$(APP).skel.h
# 用户空间目标文件路径
APP_USER_OBJ := ./build/$(APP).o

# ---- 编译和链接标志 ----

# 用户空间 C 编译器标志
# -g: 生成调试信息
# -Wall: 开启所有警告
CFLAGS := -g -Wall -O2
# 用户空间链接器标志 (基础)
LDFLAGS :=
# libbpf 依赖的库
LIBBPF_LIBS := -lelf -lz # libelf 和 zlib 通常是 libbpf 的依赖

# 头文件搜索路径
# -I$(OUTPUT): 包含生成的 .skel.h 和可能的 vmlinux.h 副本
# -I.: 包含项目本地的 .h 文件 (如 opensnoop.h)
# -I$(LIBBPF_SRC)/.. : 指向 libbpf 源码的根目录，以便能找到 include/uapi 等
# -I$(VMLINUX_H_DIR): 包含 vmlinux.h
INCLUDES := -I./build -I$(OUTPUT) -I. -I$(abspath ./libbpf-bootstrap/libbpf/include/uapi) -I$(VMLINUX_H_DIR)
# 将 INCLUDES 添加到 CFLAGS
CFLAGS += $(INCLUDES)

# BPF 编译器标志
# -g: 生成调试信息 (对 CO-RE 很重要)
# -O2: 优化级别
# -target bpf: 指定目标架构为 BPF
# -D__TARGET_ARCH_$(ARCH): 定义目标架构宏 (vmlinux.h 可能需要)
BPF_CFLAGS := -g -O2 -target bpf -D__TARGET_ARCH_$(ARCH)
# 添加头文件路径到 BPF 编译标志
BPF_CFLAGS += $(INCLUDES)

# 获取 Clang 的默认系统包含路径，并使用 -idirafter 添加
# 这有助于在某些系统上找到 BPF 编译所需的架构特定头文件 (如 asm/types.h)
# 仅在系统头文件缺失时才生效，避免干扰正常包含
CLANG_BPF_SYS_INCLUDES ?= $(shell $(CLANG) -v -E - </dev/null 2>&1 \
	| sed -n '/<...> search starts here:/,/End of search list./{ s| \(/.*\)|-idirafter \1|p }')
BPF_CFLAGS += $(CLANG_BPF_SYS_INCLUDES)

# 最终链接时使用的所有库标志
ALL_LDFLAGS := $(LDFLAGS) $(LIBBPF_LIBS)

# ---- Verbosity Control (输出控制) ----
# 如果执行 make V=1，则显示详细命令，否则只显示简短消息
ifeq ($(V),1)
	Q =
	msg =
else
	Q = @
	msg = @printf '  %-8s %s%s\n' \
		      "$(1)" \
		      "$(patsubst $(abspath $(OUTPUT))/%,%,$(2))" \
		      "$(if $(3), $(3))";
	# 不打印目录切换信息
	MAKEFLAGS += --no-print-directory
endif

# ---- 编译目标 ----

# 默认目标：编译最终的应用程序
.PHONY: all
all: $(APP)

# 清理目标：删除所有生成的文件
.PHONY: clean
clean:
	$(call msg,CLEAN,$(OUTPUT) $(APP))
	$(Q)$(RM) -r $(OUTPUT) $(APP)

# 确保输出目录存在
$(OUTPUT):
	$(call msg,MKDIR,$@)
	$(Q)mkdir -p $@

# 确保 bpftool 的输出目录存在
$(BPFTOOL_OUTPUT):
	$(call msg,MKDIR,$@)
	$(Q)mkdir -p $@

# 编译 libbpf 静态库
# 依赖：libbpf 源码文件 和 Makefile
# 依赖：输出目录 $(OUTPUT)
$(LIBBPF_OBJ): $(wildcard $(LIBBPF_SRC)/*.[ch] $(LIBBPF_SRC)/Makefile) | $(OUTPUT)
	$(call msg,LIB,$@)
	$(Q)$(MAKE) -C $(LIBBPF_SRC) BUILD_STATIC_ONLY=1 \
		    CFLAGS="" LDFLAGS="" \
		    OBJDIR=$(abspath $(OUTPUT)/libbpf) \
		    DESTDIR=$(abspath $(OUTPUT)) \
		    INCLUDEDIR= LIBDIR= UAPIDIR= \
		    prefix= \
		    install
	$(Q)cp $(abspath $(OUTPUT)/libbpf)/libbpf.a $@


# 编译 bpftool 工具
# 依赖：bpftool 源码
# 依赖：bpftool 输出目录 $(BPFTOOL_OUTPUT)
$(BPFTOOL): $(wildcard $(BPFTOOL_SRC)/*.[ch] $(BPFTOOL_SRC)/Makefile) | $(BPFTOOL_OUTPUT)
	$(call msg,BPFTOOL,$@)
	$(Q)$(MAKE) ARCH=$(ARCH) CROSS_COMPILE= \
		    OUTPUT=$(BPFTOOL_OUTPUT)/ \
		    -C $(BPFTOOL_SRC) bootstrap


# 编译 BPF C 源文件 -> BPF 目标文件 (.bpf.o)
# 依赖：BPF C 源文件, 依赖的头文件 (.h), vmlinux.h, libbpf (间接), bpftool (生成最终对象)
# 依赖：输出目录 $(OUTPUT), bpftool 工具
$(APP_BPF_OBJ): $(APP_BPF_SRCS) $(VMLINUX) $(BPFTOOL) | $(OUTPUT)
	$(call msg,BPF,$@)
	# 编译 BPF C 代码为临时的 .o 文件
	$(Q)$(CLANG) $(BPF_CFLAGS) -c $(APP_BPF_SRCS) -o $(@).tmp
	$(Q)$(BPFTOOL) gen object $@ $(@).tmp
	$(Q)$(RM) $(@).tmp

# 生成 BPF 骨架头文件 (.skel.h)
# 依赖：BPF 目标文件 (.bpf.o)
# 依赖：输出目录 $(OUTPUT), bpftool 工具
$(APP_SKEL): $(APP_BPF_OBJ) | $(OUTPUT) $(BPFTOOL)
	$(call msg,GEN-SKEL,$@)
	$(Q)$(BPFTOOL) gen skeleton $< > $@

# 编译用户空间 C 源文件 -> 用户空间目标文件 (.o)
# 依赖：用户空间 C 源文件, 依赖的头文件 (.h), 生成的骨架头文件 (.skel.h)
# 依赖：输出目录 $(OUTPUT)
$(APP_USER_OBJ): $(APP_C_SRCS) $(APP_SKEL) | $(OUTPUT)
	$(call msg,CC,$@)
	$(Q)$(CC) $(CFLAGS) -c $(APP_C_SRCS) -o $@

# 链接最终的应用程序
# 依赖：用户空间目标文件 (.o), libbpf 静态库 (.a)
$(APP): $(APP_USER_OBJ) $(LIBBPF_OBJ)
	$(call msg,BINARY,$@)
	$(Q)$(CC) $(CFLAGS) $(APP_USER_OBJ) $(LIBBPF_OBJ) $(ALL_LDFLAGS) -o ./build/$@

# ---- 其他设置 ----

# 如果编译过程中发生错误，删除生成了一半的目标文件
.DELETE_ON_ERROR:

# 保留中间文件 (.skel.h, .bpf.o, .o)，即使它们不是最终目标
.SECONDARY: $(APP_BPF_OBJ) $(APP_SKEL) $(APP_USER_OBJ)

