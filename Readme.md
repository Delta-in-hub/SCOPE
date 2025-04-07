# SCOPE: **S**calable **C**omprehensive **O**bservability **P**latform with **e**BPF



技术栈：

- libbpf
  - uprobe
  - probe points:
    - cudaMalloc
    - cudaMemcpy
    - cudaFree
    - PyCallFuncEntry
    - PyCallFuncExit
    - 进程 PID, COMM：
      - 
- golang
  - chi
  - pgx
  - sqlx
  - github.com/denisbrodbeck/machineid
- zeromq
- TimescaleDB
- Grafana
- docker
- perfetto tracing json  （Tracing 图）
- Postgres
- redis
- Nginx
- React



我想要实现一个系统，使用 eBPF（libbpf）观测系统, 比如：

- 类似 GPUprobe 是一个零代码修改、低开销、应用级的 GPU 监控工具，通过 eBPF uprobes 检查 CUDA 运行时 API 调用。
- 监控 ollama
  - tcpdump 实现 Ollama 流量分析
  - llama.cpp
  -  CUDA 运行时库 (libcudart.so, libcuda.so 等)。
    - 显存分配/释放: 跟踪 cudaMalloc, cudaFree 等调用，了解 Ollama（通过 llama.cpp）的显存使用模式，检测潜在泄漏。
    - 数据传输: 监控 cudaMemcpy 等函数，了解主机与设备之间的数据拷贝频率和大小。
    - Kernel 执行: 监控 cudaLaunchKernel 等，了解哪些计算核心（kernel）正在被调用，调用频率如何。可以像 GPUprobe 那样尝试解析内核符号名称。
  - 监控 llama.cpp 内部函数 （要思考符号问题）
- 进程的系统调用监控
  - 通过 PID ,COMM 过滤
- 资源使用追踪



使用 HASH MAP 在 probe 之间进行通信， 使用 ringbuf 将数据传输到用户态C程序。 用户态C程序再把数据通过zeromq等ipc手段，传给 go 语言程序。 go语言程序将数据收集，保存到数据库（可以考虑时序数据库），同时可以做数据转换等。

此外， golang 使用 chi ，pgx， sqlx 等实现服务器后端，包括用户管理的相关功能，双token等。

后端通过 restfulapi， 前后端分类，前端使用 Grafana ，perfetto 等将数据合理美观展示。

用 Docker 部署开发。

Nginx 反代实现 https 卸载。

设计上要考虑多机器监控，分布式。





## GOAL



**项目名称:**

- **全称 (Long Name):** **S**calable **C**omprehensive **O**bservability **P**latform with **e**BPF (SCOPE)
- **中文:** 基于 eBPF 的可扩展综合可观测性平台
- **核心理念:** 构建一个功能全面、可灵活扩展、低侵入性的分布式系统观测解决方案，利用 eBPF 技术提供深度系统洞察，并通过特权辅助服务安全地管理 eBPF 探针的部署与生命周期。

**项目目标:**

构建一个低开销、高可观测性的分布式系统监控平台。该平台利用 eBPF 技术在内核层面无侵入地采集应用（特别是 GPU 应用如 Ollama、Python 应用）和系统的性能指标、函数调用、系统调用、资源使用等数据。通过在每台受监控机器上部署一个由 systemd 管理的特权 Go 服务（Agent 管理器），负责按需启动和停止具体的 eBPF 数据采集 Agent (C 语言实现)。采集到的数据通过分布式架构高效汇聚到中心化的 Go 后端进行处理、存储和分析，最终通过 Grafana、Perfetto 和自定义 React 界面进行多维度可视化展示，为开发者和运维人员提供深入、实时、可定制、可远程管理的系统洞察力。

**核心功能:**

1. **无侵入应用级监控:**
   - **GPU 应用监控 (类 GPUprobe):**
     - 使用 libbpf 创建 uprobes 动态挂载到目标进程的 CUDA 运行时库 (如 libcudart.so, libcuda.so)。
     - **探测点:** cudaMalloc, cudaFree (追踪显存分配/释放，分析使用模式，检测泄漏), cudaMemcpy (监控主机<->设备数据传输频率、大小、方向), cudaLaunchKernel (识别调用的计算核心、频率，尝试解析 Kernel 名称)。
   - **Python 应用监控:**
     - 使用 uprobes 监控 Python 解释器 C API 关键函数。
     - **探测点:** PyCallFuncEntry, PyCallFuncExit (追踪 Python 函数调用栈、执行时间，分析性能瓶颈)。
   - **特定应用监控 (Ollama / llama.cpp):**
     - 结合上述 CUDA 监控能力。
     - **网络流量分析:** 利用 eBPF 的 socket filter 或 kprobe 挂载到网络相关系统调用 (如 sendmsg, recvmsg) 过滤 Ollama 的网络流量，或作为补充使用 tcpdump。
     - **内部函数监控 (挑战):** 尝试使用 uprobes 监控 llama.cpp 编译后的二进制文件中的特定内部函数。需要重点解决 C++ 符号名称修饰 (Name Mangling) 和符号查找的稳定性问题（可能依赖调试信息或特定符号解析库）。
2. **进程/系统级监控:**
   - **系统调用追踪:** 使用 tracepoints 或 kprobes 监控关键系统调用 (如 openat, read, write, execve, clone 等)。
   - **过滤:** eBPF 程序内通过检查当前进程的 PID 和 COMM (进程名) 实现只监控目标进程。
   - **资源使用追踪:** 结合系统调用和 eBPF 辅助函数获取进程的 CPU 时间、内存使用（可能通过监控 mmap/munmap 等）、I/O 统计等信息。
   - **进程生命周期:** 监控 execve/exit 等系统调用追踪进程创建和退出事件。
3. **高效数据采集与传输:**
   - **内核态数据共享:** 使用 eBPF HASH MAP 在不同探针之间共享状态。例如，在函数入口探针记录时间戳和参数，在函数出口探针读取并计算执行时间，然后将完整事件发送。
   - **内核态到用户态:** 使用 eBPF Ring Buffer (ringbuf) 作为主要的、高性能的、内存安全的数据传输通道，将采集到的结构化事件数据批量发送到用户态。
   - **用户态数据采集 Agent (C 语言):**
     - **职责:** 专注于使用 libbpf 框架与内核 eBPF 程序交互（加载、挂载、Map 操作、读取 Ring Buffer），执行实际的 eBPF 监控任务。此 Agent 由特权的 Agent 管理服务启动。
     - 从 Ring Buffer 中高效读取批量事件数据。
     - **数据序列化:** 将 C 结构体数据序列化成适合网络传输的格式（例如 MessagePack, 或自定义二进制格式），需要仔细选择以平衡性能和易用性。
     - **跨进程/机器传输:** 使用 ZeroMQ (C/C++ 库) 将序列化后的数据通过 PUSH/PULL 或 PUB/SUB 模式发送给中心化的 Go 后端服务。需要考虑连接管理、重连机制。
   - **Agent 管理服务 (Go 语言):**
     - **部署:** 在每台目标机器上，使用 systemd 配置并运行一个 Go 编写的服务，该服务以 root 权限运行。
     - **职责:**
       - **接收控制命令:** 通过一个安全的通信通道（例如 ZeroMQ REQ/REP 或 gRPC）从中心化的 Go 后端服务接收指令（如启动对特定进程/应用的监控，停止某个监控任务）。
       - **管理 C Agent 生命周期:** 根据收到的指令，负责启动或终止对应的 C 语言数据采集 Agent 进程。每个 C Agent 实例可能负责一种或一组特定的 eBPF 探针。
       - **状态上报:** （可选）向后端上报其管理下的 C Agent 的运行状态。
     - **通信:** 建立 Go 后端到各 Go Agent 管理服务的通信机制，用于下发控制命令。
4. **分布式后端数据处理与存储:**
   - **Go 后端服务:**
     - **数据接收:** 使用 ZeroMQ Go Binding (如 pebbe/zmq4) 接收来自多个 C Agent 的数据流。
     - **控制命令下发:** 实现与各 Agent 管理服务 (Go) 的通信（例如，通过 ZeroMQ 或 gRPC），发送启动/停止监控任务的指令。
     - **数据处理:** 采用 Go 的并发原语 (goroutines, channels) 高效地进行数据反序列化、解析、校验。
     - **数据丰富化:** 使用 github.com/denisbrodbeck/machineid 获取 Agent 所在机器的唯一 ID，并结合时间戳、元数据（如应用名、环境标签）丰富事件信息。
     - **数据转换:** 将原始事件数据转换为适合存储到时序数据库和关系数据库的格式。生成 Perfetto UI 可识别的 JSON 追踪格式。
   - **数据存储:**
     - **时序数据 (Metrics):** 使用 InfluxDB 存储聚合后的性能指标、资源使用率、事件计数等时间序列数据。使用 InfluxDB Go Client 进行交互。
     - **事件/追踪数据 (Logs/Traces):** 使用 PostgreSQL (通过 pgx/sqlx 驱动) 存储详细的、结构化的事件日志、函数调用记录、系统调用记录等。Perfetto JSON 格式的追踪数据可以直接存储为文本或 JSONB 类型，或解析后存入专用表结构。
     - **元数据/用户数据:** 使用 PostgreSQL 存储系统配置（如监控目标列表）、Agent 信息、用户信息、认证凭据、权限等关系型数据。
5. **用户认证与管理:**
   - Go 后端使用 Chi 框架构建 RESTful API，提供用户注册、登录接口。
   - 实现基于 JWT 的双 Token 机制 (短效 Access Token + 长效 Refresh Token) 进行 API 访问控制和会话管理。
6. **多维数据可视化与 API:**
   - **后端 API:** Go (Chi) 提供安全的 RESTful API，供前端查询处理后的指标、事件、追踪数据，以及用于触发 Agent 管理操作（如启动/停止特定监控）的接口。
   - **数据展示:**
     - Grafana: 连接 InfluxDB 和 PostgreSQL 数据源，创建仪表盘，展示系统概览、性能趋势、资源消耗等指标。
     - Perfetto UI: 加载由 Go 后端生成的 Perfetto JSON 追踪文件，提供强大的、交互式的、细粒度的性能分析能力（火焰图、调用栈、事件时间线等）。
     - React 前端: 开发自定义的管理界面，增加用于管理（新增、查看、停止）不同机器上运行的 eBPF 监控任务的功能，作为 Grafana 和 Perfetto 的补充和整合入口。
7. **分布式架构设计:**
   - **多 Agent 支持:** 架构设计支持在多台目标机器上部署独立的 **Agent 管理服务 (Go)** 和由其管理的 **数据采集 Agent (C)**。
   - **中心化后端:** Go 后端服务设计为可水平扩展的集群（虽然初期可能单实例），统一处理和存储所有 Agent 发来的数据。
   - **机器识别:** 每个 C Agent 启动时获取本机唯一 ID (machineid) 并随数据上报，以便后端区分和聚合来自不同机器的数据。Agent 管理服务也需识别自身所在机器。
   - **控制流:** 用户通过 React 前端发起监控操作 -> 请求发送到 Go 后端 API -> Go 后端通过控制通道 (如 ZeroMQ/gRPC) 将命令发送给目标机器上的 Go Agent 管理服务 -> Go Agent 管理服务执行命令，启动或停止相应的 C 数据采集 Agent 进程。
8. **部署与运维:**
   - **容器化:**
     - 使用 Docker 和 Docker Compose 封装 Go Backend、PostgreSQL、InfluxDB、Redis、Nginx、React App。
     - **Agent 部署:**
       - **Go Agent 管理服务:** 需要在目标宿主机上部署。提供 systemd service 文件模板，用于将其安装为系统服务并以 root 权限运行。其运行可能需要访问宿主机的某些资源（如启动 C Agent 进程）。直接在宿主机上通过 systemd 运行可能是较简单的方式。
       - **C 数据采集 Agent:** 其二进制文件需要与 Go Agent 管理服务一起部署在目标宿主机上，以便管理服务可以按需启动它。
   - **反向代理与安全:** 使用 Nginx 作为前端入口，处理静态资源请求（React App），反向代理 API 请求到 Go 后端集群，实现负载均衡，并进行 SSL/TLS 卸载，强制 HTTPS 通信。
   - **缓存:** 使用 Redis (通过 Redis Go Client) 缓存热点数据、用户会话信息、或用作简单的消息队列/分布式锁，提升后端性能和响应速度。

**技术栈梳理:**

- **eBPF & 内核监控:** libbpf (C), uprobes, tracepoints/kprobes , eBPF Maps (HASH MAP), eBPF Ring Buffer
- **数据转发 Agent:** C, libbpf, ZeroMQ (C/C++ lib), Serialization library (e.g., msgpack-c)
- **Agent 管理服务:** Go, ZeroMQ/gRPC (for control), systemd (for service management)
- **后端服务:** Go, Chi (Web Framework), pgx/sqlx (PostgreSQL Driver), pebbe/zmq4 (ZeroMQ Binding for data), ZeroMQ/gRPC Client (for control), influxdb-client-go (InfluxDB Client), go-redis (Redis Client), github.com/denisbrodbeck/machineid, Go Concurrency Primitives
- **数据库 & 缓存:** PostgreSQL, InfluxDB, Redis
- **前端 & 可视化:** React, Grafana, Perfetto (UI & JSON format)
- **基础设施 & 部署:** Docker, Docker Compose, Nginx, systemd

**数据流:**

1. **触发:** 目标应用/内核执行到 eBPF 探针挂载点。
2. **采集:** eBPF 程序被激活，收集上下文信息 (参数, 返回值, PID, COMM, 时间戳等)。
3. **内核处理:** 可选地使用 eBPF Map 聚合或关联信息 (如计算耗时)。
4. **传输 (内核->用户态):** 结构化事件数据推入 eBPF Ring Buffer。
5. **读取 & 序列化:** 被 Go Agent 管理服务启动的用户态 C Agent 通过 libbpf 从 Ring Buffer 读取批量事件，并将其序列化， 添加机器 ID。
6. **转发 (Agent->后端):** C Agent 使用 ZeroMQ 将序列化数据发送到 Go 后端。
7. **接收 & 处理:** Go 后端监听 ZeroMQ，接收数据，反序列化，解析，校验等元数据。
8. **转换 & 存储:** Go 后端将数据转换为目标格式：
   - 时序指标写入 InfluxDB。
   - 详细事件/追踪数据写入 PostgreSQL。
   - 按需生成 Perfetto JSON 追踪文件并存储 (可能在 Postgres 或文件系统)。
9. **查询:** 前端 (Grafana, Perfetto, React) 通过 Go 后端的 RESTful API 请求数据。
10. **可视化:** Grafana 查询 DB 展示仪表盘；Perfetto UI 加载 JSON 文件进行追踪分析；React App 展示自定义视图。

**控制流:**

1. **触发:** 用户在 React 前端点击按钮，请求在某台机器上启动对特定应用（如 Ollama）的 GPU 监控。
2. **API 请求:** React 前端调用 Go 后端的 API，传递机器 ID 和监控任务类型。
3. **命令下发:** Go 后端服务通过其控制通道 (如 ZeroMQ/gRPC) 向指定机器 ID 的 Go Agent 管理服务发送“启动 GPU 监控”的命令，可能包含目标进程信息。
4. **执行命令:** 目标机器上的 Go Agent 管理服务收到命令。
5. **启动 C Agent:** 该服务以 root 权限（因为它本身是 root 运行的）启动一个配置好的 C 数据采集 Agent 进程，该 C Agent 专门负责 GPU 监控的 eBPF 程序加载和数据采集/发送。
6. **(反向操作)** 停止监控时，流程类似，最终由 Go Agent 管理服务终止对应的 C Agent 进程。

**设计考虑与最佳实践:**

- **模块化:** 清晰划分 eBPF 程序、C Agent、Go Agent 管理服务、Go 后端、数据库、前端的职责。
- **配置化:** 监控目标、采样率、数据保留期、Agent 连接信息、Agent 管理器的监听地址、与后端的通信方式等应外部化配置。
- **错误处理与韧性:** 实现 Agent 与后端的连接重试、数据处理失败的记录与告警、组件间的健康检查。增加对 Agent 管理服务与后端通信、管理 C Agent 进程失败的处理。
- **性能优化:**
  - 选择合适的 Ring Buffer 大小和用户态拉取策略。
  - ZeroMQ 传输模式 (PUSH/PULL vs PUB/SUB) 和序列化格式的选择对性能影响显着，需测试评估。
  - Go 后端充分利用并发能力处理数据流。
  - 数据库设计合理的索引，优化查询语句。
- **安全性:**
  - API 使用 Token 认证，强制 HTTPS 通信。
  - Agent 管理服务与后端之间的控制通道需要加密和认证，防止未授权的控制命令。
  - 严格限制 Agent 管理服务的功能，仅限于启动/停止预定义的、受信任的 C Agent 二进制文件，避免任意命令执行漏洞。
  - 将 eBPF 操作所需的 root 权限集中到 Go Agent 管理服务，确保该服务自身安全。
- **可扩展性:** Go 后端服务设计为无状态或状态易于外部化 (Redis/DB)，便于水平扩展；数据库根据需要考虑分片或集群方案。
- **测试:** 编写单元测试 (Go, C), 集成测试 (组件间交互), 端到端测试 (模拟用户场景)。eBPF 程序测试需要特定环境和策略。增加对 Agent 管理服务、控制流、systemd 部署的测试。
- **文档:** 提供清晰的架构图、数据模型说明、API 文档、部署指南，包括 Agent 管理服务的配置、API 和部署说明。

**潜在挑战:**

- **符号解析:** 获取 C++ (如 llama.cpp) 用户态函数准确名称和参数是主要难点，可能需要依赖调试符号 (-g) 或专门的库/技术。
- **数据量与成本:** 高频监控可能产生巨大数据量，对网络带宽、处理能力、存储成本带来压力，需要实施采样、聚合或过滤策略。
- **eBPF 兼容性:** 不同 Linux 内核版本对 eBPF 功能支持度不同，需确定最低内核要求或使用 CO-RE (Compile Once – Run Everywhere) 技术。
- **Agent 管理复杂性:** 管理分布在多台机器上的 Agent 管理服务及其启动的 C Agent 进程，包括版本更新、状态监控、错误排查等，增加了运维复杂性。
- **控制通道安全:** 确保 Go 后端与所有 Agent 管理服务之间的控制命令通道安全可靠。
- **部署一致性:** 确保 Go Agent 管理服务和 C Agent 二进制文件在所有目标机器上正确部署和配置。





## Dirs



```
scope/
├── .github/              # CI/CD 工作流 (e.g., GitHub Actions)
│   └── workflows/
│       ├── build.yml
│       └── test.yml
├── .gitignore            # Git 忽略文件配置
├── README.md             # 项目根 README，包含概述、快速开始、架构图链接等
├── LICENSE               # 项目许可证文件
├── go.mod                # Go 模块定义
├── go.sum                # Go 模块校验和
├── Makefile              # 顶层 Makefile，用于协调构建、测试、部署等任务
│
├── api/                  # API 定义与规范 (例如 OpenAPI, Protobuf)
│   ├── openapi/          # RESTful API 的 OpenAPI v3 规范
│   │   └── v1/
│   │       └── scope.yaml
│   └── proto/            # (如果使用 gRPC) Protobuf 定义
│       └── v1/
│           └── control.proto # Agent 控制命令的 proto 定义
│
├── bpf/                  # 所有 eBPF 相关代码
│   ├── src/              # 内核态 eBPF 程序 C 源码 (.c)
│   │   ├── probes/       # 按功能组织的 eBPF 探针代码
│   │   │   ├── cuda_monitor.c
│   │   │   ├── python_monitor.c
│   │   │   └── syscall_monitor.c
│   │   ├── headers/      # 通用头文件, vmlinux.h (如果使用 BTF/CO-RE)
│   │   │   └── common.h
│   │   └── bpf_defs.h    # eBPF 程序使用的共享结构体定义 (也会被 C Agent 包含)
│   ├── agent/            # 用户态 C 数据采集 Agent
│   │   ├── src/          # C Agent 源码
│   │   │   ├── main.c
│   │   │   ├── ringbuf_reader.c
│   │   │   ├── zmq_sender.c
│   │   │   └── event_serializer.c # (例如使用 msgpack-c)
│   │   ├── include/      # C Agent 头文件
│   │   │   └── agent.h
│   │   └── lib/          # 依赖库 (可选, 如 libbpf submodule 或 msgpack-c)
│   └── Makefile          # 用于编译 eBPF 程序和 C Agent 的 Makefile
│
├── cmd/                  # 项目的可执行应用程序入口
│   ├── scope-backend/    # 中心化 Go 后端服务
│   │   └── main.go
│   └── scope-agent-manager/ # Go Agent 管理服务 (部署在被监控机器)
│       └── main.go
│
├── config/               # 配置文件模板或默认配置
│   ├── backend/
│   │   └── config.yaml.example
│   └── agent-manager/
│       └── config.yaml.example
│
├── database/             # 数据库相关文件
│   ├── postgres/
│   │   ├── migrations/   # PostgreSQL 数据库迁移脚本 (e.g., using golang-migrate)
│   │   │   ├── 001_initial_schema.up.sql
│   │   │   └── 001_initial_schema.down.sql
│   │   └── schema.sql    # (可选) 完整的数据库 Schema 定义
│   └── influxdb/         # InfluxDB 相关配置说明或脚本 (Buckets, Tasks)
│       └── setup_notes.md
│
├── deploy/               # 部署相关配置和脚本
│   ├── docker/           # Docker & Docker Compose 配置
│   │   ├── docker-compose.yml # 编排后端服务 (Backend, DBs, Nginx, Grafana, Redis)
│   │   ├── backend/
│   │   │   └── Dockerfile
│   │   ├── frontend/
│   │   │   └── Dockerfile     # 用于构建和服务 React App
│   │   ├── nginx/
│   │   │   └── nginx.conf     # Nginx 配置文件
│   │   └── grafana/
│   │       └── provisioning/  # Grafana 数据源和仪表盘自动配置
│   │           ├── datasources/
│   │           │   └── default.yaml
│   │           └── dashboards/
│   │               └── default.yaml
│   └── systemd/          # systemd 服务单元文件
│       └── scope-agent-manager.service.template # Agent Manager 服务模板
│
├── docs/                 # 项目文档
│   ├── architecture.md   # 架构设计文档
│   ├── api.md            # API 使用说明 (或链接到生成的文档)
│   ├── deployment.md     # 部署指南 (后端堆栈和 Agent)
│   ├── development.md    # 开发环境设置和指南
│   ├── bpf_probes.md     # eBPF 探针详细说明
│   └── data_flow.md      # 数据流和控制流说明
│
├── internal/             # 私有的应用程序代码 (不应被外部项目导入)
│   ├── agentmanager/     # Go Agent Manager 的内部实现
│   │   ├── cmdctrl/      # 处理来自 Backend 的控制命令 (ZMQ/gRPC Server)
│   │   ├── config/       # Agent Manager 配置加载与管理
│   │   ├── identity/     # 获取机器 ID
│   │   └── proc/         # C Agent 进程管理逻辑
│   ├── backend/          # Go Backend 的内部实现
│   │   ├── api/          # API 层 (Chi 路由、处理器、中间件)
│   │   ├── app/          # 核心业务逻辑/服务层
│   │   ├── agentctrl/    # 与 Agent Manager 通信的客户端逻辑 (ZMQ/gRPC Client)
│   │   ├── auth/         # 用户认证与 JWT 处理
│   │   ├── cache/        # Redis 缓存交互
│   │   ├── config/       # Backend 配置加载与管理
│   │   ├── data/         # 数据接收 (ZMQ Pull)、解析、转换 (Perfetto)
│   │   └── storage/      # 数据存储层 (PostgreSQL, InfluxDB 交互)
│   ├── models/           # 跨内部包共享的数据模型 (事件结构体, API 请求/响应体)
│   ├── platform/         # 平台相关工具 (如 machineid 的封装)
│   └── shared/           # 项目内部共享的通用库/工具 (如自定义错误、日志封装)
│       ├── errors/
│       └── logger/
│
├── pkg/                  # 公开的库代码 (如果希望其他项目可以导入)
│   └── eventtypes/       # (示例) 如果事件类型定义需要被外部工具使用
│
├── scripts/              # 构建、测试、部署等辅助脚本
│   ├── build.sh          # 主构建脚本 (调用其他脚本)
│   ├── build_bpf.sh      # 编译 eBPF 和 C Agent
│   ├── build_go.sh       # 编译 Go 应用
│   ├── build_web.sh      # 构建 React 应用
│   ├── run_dev.sh        # 本地开发环境启动脚本
│   ├── test.sh           # 运行测试 (单元测试, 集成测试, lint)
│   └── db_migrate.sh     # 应用数据库迁移脚本
│
├── test/                 # 测试相关文件 (集成测试、端到端测试)
│   ├── integration/
│   └── e2e/
│
└── web/                  # 前端 React 应用
    └── app/              # React 项目根目录
        ├── public/
        ├── src/
        ├── package.json
        ├── tsconfig.json   # (如果是 TypeScript)
        └── ...           # 其他 React 项目文件
```











## Steps



**阶段一：核心数据管道原型 (单机手动运行)**

- 
- **目标:** 验证 eBPF -> C Agent -> Go Backend 的基本数据流。
- **步骤:**
  1. 
  2. **eBPF 程序 (C/libbpf):**
     - 
     - 创建一个简单的 eBPF 程序（例如 minimal_probe.c），使用 libbpf。
     - 选择一个简单的 Tracepoint (如 syscalls/sys_enter_execve) 或 Kprobe。
     - 定义一个简单的数据结构 (e.g., struct event { pid_t pid; char comm[16]; };)。
     - 在 eBPF 程序中获取 PID 和 COMM，填充结构体。
     - 设置一个 eBPF Ring Buffer (BPF_MAP_TYPE_RINGBUF)。
     - 将填充好的事件数据提交到 Ring Buffer。
     - 使用 bpftool 或 libbpf 的辅助工具编译 eBPF 程序为 BPF 对象文件 (.o)。
  3. **用户态 C Agent (数据采集):**
     - 
     - 创建一个 C 程序 (c_agent.c)。
     - 使用 libbpf 加载上面编译的 BPF 对象文件 (.o)。
     - 查找并附加 (attach) eBPF 程序到指定的 Tracepoint/Kprobe。
     - 查找 Ring Buffer Map 的文件描述符 (FD)。
     - 设置 Ring Buffer 的回调函数。
     - 在回调函数中，从 Ring Buffer 读取事件数据。
     - **[临时]** 将读取到的数据打印到标准输出，验证数据采集。
     - **[集成 ZMQ]** 引入 ZeroMQ C 库 (libzmq)。
     - 初始化 ZeroMQ Context 和一个 PUSH socket。
     - 连接到 Go Backend 将要监听的地址 (e.g., tcp://localhost:5555)。
     - 将从 Ring Buffer 读取到的事件数据（C 结构体）序列化（初期可直接发送原始字节，后续考虑 MessagePack 或 Protobuf）并通过 ZMQ PUSH socket 发送出去。
  4. **用户态 Go Backend (数据接收与存储):**
     - 
     - 创建一个 Go 程序 (go_backend/main.go)。
     - 使用 pebbe/zmq4 Go binding。
     - 初始化 ZeroMQ Context 和一个 PULL socket。
     - 绑定到 C Agent 连接的地址 (e.g., tcp://*:5555)。
     - 在一个 goroutine 中循环接收来自 C Agent 的消息。
     - 反序列化接收到的数据（与 C Agent 的序列化方式对应）。
     - **[临时]** 将接收到的数据打印到标准输出。
     - **[集成 Postgres]** 引入 pgx/sqlx。
     - 连接到本地运行的 PostgreSQL 数据库。
     - 设计一个简单的表来存储接收到的事件 (e.g., events table with timestamp, pid, comm columns)。
     - 将反序列化后的数据插入到 PostgreSQL 表中。
  5. **数据库准备:**
     - 
     - 启动一个 PostgreSQL 实例（可用 Docker）。
     - 创建数据库和前面设计的 events 表。
  6. **测试:**
     - 
     - 编译 eBPF 程序。
     - 编译 C Agent。
     - 运行 Go Backend。
     - 以 root 权限运行 C Agent。
     - 触发 eBPF 探针（例如，执行一个新命令来触发 execve）。
     - 检查 Go Backend 的日志输出和 PostgreSQL 中的数据。
- **产出:** 一个能在单机上手动运行的、从内核 eBPF 探针采集数据，通过 C Agent 和 ZeroMQ 发送到 Go Backend，并存入 PostgreSQL 的最小化数据管道。

------



**阶段二：实现核心 eBPF 监控功能**

- 
- **目标:** 实现对 CUDA、Python 和系统调用的具体监控逻辑。
- **步骤:**
  1. 
  2. **CUDA 监控 (eBPF & C Agent):**
     - 
     - 扩展 eBPF 程序，添加 Uprobes。
     - 探测点: cudaMalloc, cudaFree, cudaMemcpy, cudaLaunchKernel (需要目标进程的 libcudart.so 路径)。
     - 设计更复杂的事件结构体，包含函数名、参数（尽可能获取）、返回值、时间戳、PID、COMM。
     - 使用 eBPF Hash Map 在入口 (entry) 和出口 (return) 探针间传递状态（如记录开始时间戳）。
     - 在出口探针计算耗时，并将完整的事件（包含耗时）提交到 Ring Buffer。
     - 更新 C Agent 加载和附加这些 Uprobes 的逻辑。可能需要命令行参数指定目标进程 PID 或二进制路径。
     - 更新 C Agent 的事件处理和序列化逻辑以支持 CUDA 事件。
  3. **Python 监控 (eBPF & C Agent):**
     - 
     - 类似 CUDA，添加 Uprobes 到 Python 解释器的 C API 函数。
     - 探测点: PyCallFuncEntry, PyCallFuncExit (需要目标 Python 进程使用的 libpythonX.Y.so 路径)。
     - 设计 Python 函数调用事件结构体（函数名、文件名、行号 - 获取可能困难，时间戳、耗时、PID、COMM）。
     - 同样使用 Hash Map 关联 Entry 和 Exit。
     - 更新 C Agent 以支持 Python 监控的加载、事件处理和序列化。
  4. **系统调用与进程监控 (eBPF & C Agent):**
     - 
     - 扩展 eBPF 程序，使用 Tracepoints 或 Kprobes 监控关键系统调用 (e.g., openat, read, write, clone, exit)。
     - 在 eBPF 程序内部实现基于 PID 和 COMM 的过滤逻辑。
     - 设计相应的事件结构体。
     - 更新 C Agent 以支持这些探针和事件。
  5. **Go Backend 数据处理:**
     - 
     - 更新 Go Backend 的 ZMQ 接收逻辑，使其能识别和反序列化来自不同监控类型（CUDA, Python, Syscall）的事件。
     - 扩展 PostgreSQL 数据库模式，创建专门的表来存储不同类型的事件数据（或使用 JSONB 存储灵活结构）。
     - 实现将不同事件数据存入对应表的逻辑。
  6. **测试:**
     - 
     - 准备测试环境：安装 CUDA Toolkit、Python 环境、运行一个简单的 CUDA 程序、一个 Python 脚本。
     - 分别或组合运行 C Agent（可能需要指定目标 PID 或库路径）和 Go Backend。
     - 验证 PostgreSQL 中是否记录了预期的 CUDA API 调用、Python 函数调用和系统调用事件。
- **产出:** C Agent 具备了核心的监控能力，Go Backend 能够接收、区分并存储这些多样化的监控数据。

------



**阶段三：引入 Agent 管理服务与控制流**

- 
- **目标:** 实现对 C Agent 的远程、按需启动和停止，为分布式部署打下基础。
- **步骤:**
  1. 
  2. **Go Agent 管理服务 (新组件):**
     - 
     - 创建新的 Go 项目 (go_agent_manager)。
     - 该服务需要以 root 权限运行。
     - 使用 ZeroMQ (或其他如 gRPC) 实现一个 REQ/REP 或 ROUTER/DEALER 模式的监听 Socket，用于接收来自 Go Backend 的控制命令。定义清晰的命令格式（如 JSON: {"command": "start_cuda_monitor", "pid": 1234} 或 {"command": "stop_monitor", "monitor_id": "xyz"}）。
     - 实现命令解析逻辑。
     - 根据命令，使用 os/exec 启动或终止对应的 C Agent 进程。需要管理 C Agent 进程的生命周期（记录 PID，处理退出等）。
     - 考虑安全性：确保只执行预定义的、安全的 C Agent 启动命令。
     - **[部署]** 编写 systemd service 文件，将 Go Agent Manager 安装为系统服务，配置为开机启动并以 root 运行。
  3. **C Agent 调整:**
     - 
     - 确保 C Agent 可以通过命令行参数接收必要的配置（如目标 PID、监控类型、连接后端的 ZMQ 地址等）。
     - Go Agent Manager 在启动 C Agent 时传递这些参数。
  4. **Go Backend (控制命令发送):**
     - 
     - 在 Go Backend 中添加与 Go Agent Manager 通信的逻辑（使用 ZMQ REQ 或 DEALER socket）。
     - 需要知道目标机器上 Agent Manager 的地址（初期可配置，后续考虑服务发现）。
     - 实现发送启动/停止监控命令的功能，可能由后续的 API 触发。
  5. **机器 ID:**
     - 
     - 在 Go Agent Manager 启动时，使用 github.com/denisbrodbeck/machineid 获取本机唯一 ID。
     - Agent Manager 在与 Backend 通信时（例如，上报状态或响应命令时）可以带上此 ID。
     - C Agent 启动时，由 Agent Manager 将机器 ID 通过参数传递给它，或者 C Agent 自己获取。C Agent 在发送数据给 Backend 时附带此机器 ID。
     - Go Backend 在存储数据时，记录下数据来源的机器 ID。
  6. **测试 (单机):**
     - 
     - 编译 Go Agent Manager。
     - 将其配置为 systemd 服务并启动。
     - 手动（或通过简单的测试脚本）向 Go Agent Manager 发送启动 C Agent 的命令（通过 ZMQ）。
     - 验证 C Agent 是否被正确启动，并且数据是否开始流向 Go Backend。
     - 发送停止命令，验证 C Agent 是否被终止。
     - 检查 Go Backend 接收到的数据是否包含机器 ID。
- **产出:** 一个能在目标机器上运行的、可通过 Go Backend 控制的 Agent 管理服务，能够按需启动/停止 C 数据采集 Agent。数据流中加入了机器标识。

------



**阶段四：增强数据处理、存储与可视化准备**

- 
- **目标:** 引入时序数据库，实现 Perfetto 追踪格式转换，为可视化做准备。
- **步骤:**
  1. 
  2. **集成 InfluxDB:**
     - 
     - 启动一个 InfluxDB 实例（可用 Docker）。
     - 在 Go Backend 中引入 InfluxDB Go Client (influxdb-client-go/v2)。
     - 配置 InfluxDB 连接信息（URL, Token, Org, Bucket）。
     - 在 Go Backend 的数据处理逻辑中，识别适合存为时序指标的数据（如 GPU 显存使用量、函数调用频率/耗时统计、系统调用计数等）。
     - 将这些指标数据写入 InfluxDB。
  3. **Perfetto JSON 转换:**
     - 
     - 研究 Perfetto Trace Event Format (JSON)。
     - 在 Go Backend 中，实现将收集到的函数调用（CUDA, Python）、系统调用等事件转换为 Perfetto JSON 格式中对应的事件类型（如 Slice events for duration, Instant events, Counter events）。
     - 需要精确的时间戳（纳秒精度，eBPF 可提供 bpf_ktime_get_ns()）和线程/进程信息。
     - 决定如何存储和提供这些 JSON 数据：
       - 
       - 选项 A: 实时生成并存储在 PostgreSQL 的某个字段（如 TEXT 或 JSONB）。
       - 选项 B: 批量生成追踪文件，存储在文件系统，并通过 API 提供下载链接。
       - 选项 C: 组合，存储原始事件，API 按需查询并生成 JSON。
  4. **PostgreSQL Schema 优化:**
     - 
     - 根据实际采集到的数据和查询需求，优化 PostgreSQL 的表结构和索引。
     - 考虑使用分区（Partitioning）来管理大量的事件数据。
  5. **数据丰富化:**
     - 
     - 确保所有存储的数据（Postgres, InfluxDB）都关联了时间戳、机器 ID、进程 PID、进程 COMM 等关键元数据。
  6. **测试:**
     - 
     - 运行完整的流程（Agent Manager -> C Agent -> Go Backend）。
     - 验证指标数据是否成功写入 InfluxDB。
     - 验证 Go Backend 是否能生成（至少是打印或临时存储）符合 Perfetto 格式的 JSON 数据。
     - 检查 PostgreSQL 中的数据是否结构合理且包含所有元数据。
- **产出:** Go Backend 具备了将数据存储到 InfluxDB 和转换为 Perfetto 格式的能力。数据存储更加完善。

------



**阶段五：构建后端 API 和用户界面**

- 
- **目标:** 提供用户交互接口，包括数据查询、监控任务管理和用户认证。
- **步骤:**
  1. 
  2. **Go Backend RESTful API (Chi):**
     - 
     - 使用 chi 框架构建 API 服务。
     - 实现用户管理 API: 注册、登录（返回 JWT Access Token 和 Refresh Token）。
     - 实现中间件进行 JWT 验证。
     - 实现 API 端点用于：
       - 
       - 查询存储在 PostgreSQL 中的事件数据（支持过滤、分页）。
       - 查询存储在 InfluxDB 中的指标数据（可能需要代理或直接让前端连接 Grafana）。
       - 获取 Perfetto 追踪数据（根据上一步的存储方式，提供下载或直接返回 JSON）。
       - 触发监控任务：向指定的 Agent Manager 发送启动/停止命令（需要机器 ID 和监控参数）。
       - 查看当前活动的监控任务状态（可能需要 Agent Manager 上报状态）。
  3. **用户认证 (JWT):**
     - 
     - 实现 JWT 的生成（登录时）和校验（API 中间件）。
     - 处理 Token 过期和刷新逻辑。
     - 在 PostgreSQL 中存储用户信息和（可选的）Refresh Token。
  4. **React 前端:**
     - 
     - 初始化 React 项目。
     - 实现用户登录/注册页面。
     - 实现一个管理页面，用于：
       - 
       - 列出已知的被监控机器（可能需要一个注册机器的流程或从接收到的数据中动态发现）。
       - 在选定的机器上发起新的监控任务（调用后端 API，传递 PID、应用类型等）。
       - 查看/停止正在运行的监控任务。
     - **[初期]** 实现简单的数据展示页面，调用后端 API 显示来自 PostgreSQL 的原始事件表格。
  5. **测试:**
     - 
     - 测试用户注册、登录、Token 刷新。
     - 测试 API 的访问控制。
     - 通过 React UI 触发启动/停止监控任务，验证是否成功。
     - 通过 React UI 查询并显示基本的事件数据。
- **产出:** 一个带用户认证的 RESTful API 后端，以及一个可以管理监控任务和查看基本数据的 React 前端原型。

------



**阶段六：集成可视化工具 (Grafana & Perfetto)**

- 
- **目标:** 利用专业工具美观、高效地展示监控数据。
- **步骤:**
  1. 
  2. **Grafana 集成:**
     - 
     - 启动 Grafana 实例（可用 Docker）。
     - 配置 Grafana 数据源：连接到 InfluxDB 和 PostgreSQL。
     - 创建 Grafana 仪表盘 (Dashboards)：
       - 
       - 展示来自 InfluxDB 的时序指标（GPU 显存、利用率、函数调用速率、耗时 P99 等）。
       - 展示来自 PostgreSQL 的事件摘要或统计信息。
       - 使用变量 (Variables) 实现按机器 ID、应用名等维度过滤。
  3. **Perfetto UI 集成:**
     - 
     - 确保 Go Backend API 可以提供 Perfetto JSON 数据。
     - 在 React 前端：
       - 
       - 提供一个按钮或链接，调用后端 API 获取 Perfetto JSON 数据。
       - 或者，提供一个文件上传入口，让用户上传从后端下载的 .json / .perfetto-trace 文件到 Perfetto UI ([https://ui.perfetto.dev](https://www.google.com/url?sa=E&q=https%3A%2F%2Fui.perfetto.dev)) 或本地托管的 Perfetto UI 实例。
       - （高级）尝试在 React 应用中嵌入 Perfetto UI 组件（如果可行）。
  4. **React UI 增强:**
     - 
     - 在 React UI 中嵌入 Grafana 的图表（使用 iframe 或 Grafana 的嵌入功能）。
     - 提供更友好的 Perfetto 追踪数据访问方式。
     - 完善整体 UI/UX。
  5. **测试:**
     - 
     - 验证 Grafana 仪表盘是否能正确显示来自 InfluxDB 和 PostgreSQL 的数据。
     - 验证能否从 React UI 获取 Perfetto JSON 数据，并成功在 Perfetto UI 中加载和分析。
- **产出:** 集成了 Grafana 和 Perfetto 的可视化能力，用户可以通过专业的工具深入分析系统性能。

------



**阶段七：容器化、部署与基础设施完善**

- 
- **目标:** 实现服务的容器化部署，配置反向代理和缓存。
- **步骤:**
  1. 
  2. **Docker 化:**
     - 
     - 为 Go Backend、React (构建后的静态文件) 创建 Dockerfile。
     - 创建 docker-compose.yml 文件，编排 Go Backend, PostgreSQL, InfluxDB, Redis, Nginx, Grafana 服务。
     - 确保容器间的网络配置正确。
  3. **Nginx 配置:**
     - 
     - 配置 Nginx 作为反向代理：
       - 
       - 代理 API 请求到 Go Backend 服务。
       - 服务 React 前端的静态文件。
       - （可选）代理 Grafana 服务。
     - 配置 SSL/TLS 卸载，实现 HTTPS。
  4. **Redis 集成:**
     - 
     - 在 docker-compose.yml 中添加 Redis 服务。
     - 在 Go Backend 中引入 Redis Go Client (go-redis/redis)。
     - 使用 Redis 缓存热点数据、存储用户会话/JWT Refresh Token（如果选择这种方式）、或用于简单的分布式锁/消息队列（如果需要）。
  5. **Agent 部署 (非 Docker):**
     - 
     - 确定 Go Agent Manager 和 C Agent 二进制文件的部署方式（如使用配置管理工具 Ansible/SaltStack，或手动 SCP 和 systemctl）。
     - 提供清晰的 Agent 安装和配置指南。
     - 确保 Agent Manager 的 systemd 服务配置正确，并能在目标机器上以 root 权限运行。
  6. **测试:**
     - 
     - 使用 docker-compose up 启动整个后端和可视化堆栈。
     - 测试通过 Nginx 访问前端和 API (HTTPS)。
     - 在目标机器上部署并启动 Agent Manager 服务。
     - 通过前端 UI 触发监控，验证整个流程在容器化环境下是否正常工作。
     - 测试 Redis 缓存是否生效（如果已实现相关逻辑）。
- **产出:** 一套可通过 Docker Compose 部署的后端服务和基础设施。Agent 的部署流程和文档。

------



**阶段八：多机部署、测试与最终完善**

- 
- **目标:** 验证分布式架构，进行端到端测试，并进行最终的优化和文档编写。
- **步骤:**
  1. 
  2. **多机部署:**
     - 
     - 选择至少两台目标机器。
     - 在每台机器上部署 Go Agent Manager (systemd) 和 C Agent 二进制文件。
     - 确保每台 Agent Manager 都配置了正确的 Go Backend 地址（用于控制命令）和 C Agent 知道正确的 Go Backend 地址（用于数据发送）。
     - 确保 Go Backend 可以接收来自多个 Agent 的数据和处理来自多个 Agent Manager 的连接（如果使用 ROUTER socket）。
  3. **分布式测试:**
     - 
     - 通过 React UI 或 API，在不同的机器上启动监控任务。
     - 验证 Go Backend 能正确接收、处理、存储来自不同机器的数据，并能通过机器 ID 进行区分。
     - 验证 Grafana 和 React UI 能按机器 ID 过滤和展示数据。
     - 测试控制流在多机环境下的可靠性。
  4. **性能与压力测试:**
     - 
     - 模拟高负载场景，测试系统的性能瓶颈（eBPF 开销、网络传输、后端处理、数据库写入）。
     - 根据测试结果进行优化（如调整 Ring Buffer 大小、采样率、ZeroMQ 配置、Go Backend 并发数、数据库索引等）。
  5. **安全加固:**
     - 
     - 仔细审查 Agent Manager 的权限和执行逻辑，防止安全漏洞。
     - 确保 Go Backend 和 Agent Manager 之间的控制通道是加密和认证的（例如，使用 ZeroMQ 的 CURVE 安全机制或 gRPC 的 TLS）。
     - 对所有 API 输入进行严格验证。
  6. **错误处理与健壮性:**
     - 
     - 完善各组件的错误处理、日志记录和重试机制（如 Agent 连接后端失败）。
     - 考虑添加健康检查端点。
  7. **文档完善:**
     - 
     - 编写详细的架构文档、部署指南（包括 Agent 和后端）、API 文档、用户手册、开发贡献指南。
     - 提供 eBPF 程序、数据结构、配置选项的说明。
  8. **代码清理与审查:**
     - 
     - 进行代码格式化、移除调试代码、添加注释。
     - 进行 Code Review。
- **产出:** 一个经过多机测试、性能优化、安全加固、文档齐全的 SCOPE 可观测性平台。



