## 一、 Go Backend RESTful API (面向前端和用户)

**设计原则:**

*   **协议:** HTTPS
*   **基础路径:** `/api`
*   **版本:** `/v1` (例如 `/api/v1/...`)
*   **认证:** JWT Bearer Token (在 `Authorization` Header 中)
*   **数据格式:** JSON
*   **资源命名:** 复数名词 (e.g., `machines`, `monitors`)
*   **错误处理:** 使用标准的 HTTP 状态码，并提供统一的 JSON 错误体，例如：
    ```json
    {
      "error": {
        "code": 404,
        "message": "Monitor with id 'xyz' not found on machine 'abc'",
        "status": "NOT_FOUND" // 类似 gRPC 状态码
      }
    }
    ```
*   **查询参数:**
    *   `filter`: 用于过滤结果 (语法待定，可参考 Google AIP-160，例如 `machine_id=='abc' AND status=='RUNNING'`)
    *   `page_size`: 每页数量 (e.g., `20`)
    *   `page_token`: 用于分页的令牌 (基于光标的分页)
    *   `order_by`: 排序字段 (e.g., `create_time desc`)

---

**API 端点定义:**

**1. 认证 (Auth)**

*   `POST /api/v1/auth:login`
    *   **描述:** 用户登录
    *   **请求体:** `{ "email": "user@example.com", "password": "your_password" }`
    *   **响应 (200 OK):** `{ "access_token": "...", "refresh_token": "...", "expires_in": 3600 }`
    *   **认证:** 无需

*   `POST /api/v1/auth:register`
    *   **描述:** 用户注册
    *   **请求体:** `{ "email": "user@example.com", "password": "your_password", "display_name": "User Name" }`
    *   **响应 (201 Created):** `{ "user_id": "...", "email": "...", "display_name": "..." }`
    *   **认证:** 无需

*   `POST /api/v1/auth:refreshToken`
    *   **描述:** 使用 Refresh Token 刷新 Access Token
    *   **请求体:** `{ "refresh_token": "..." }`
    *   **响应 (200 OK):** `{ "access_token": "...", "expires_in": 3600 }`
    *   **认证:** 无需 (依赖有效的 Refresh Token)

*   `POST /api/v1/auth:logout` (可选)
    *   **描述:** 用户登出 (如果需要服务端使 Refresh Token 失效)
    *   **请求体:** `{ "refresh_token": "..." }` (或者只依赖 Access Token)
    *   **响应 (204 No Content):**
    *   **认证:** 需要有效的 Access Token

**2. 机器 (Machines)** - 代表被监控的主机

*   `GET /api/v1/machines`
    *   **描述:** 列出已知的被监控机器 (例如，最近上报过数据或有活动监控的机器)。
    *   **查询参数:** `filter`, `page_size`, `page_token`, `order_by` (e.g., `last_seen desc`)
    *   **响应 (200 OK):** `{ "machines": [ { "id": "machine_uuid_1", "hostname": "...", "last_seen": "...", "status": "ACTIVE" / "INACTIVE" }, ... ], "next_page_token": "..." }`
    *   **认证:** 需要

*   `GET /api/v1/machines/{machine_id}`
    *   **描述:** 获取特定机器的详细信息。
    *   **路径参数:** `machine_id` (e.g., 由 `github.com/denisbrodbeck/machineid` 生成的 ID)
    *   **响应 (200 OK):** `{ "id": "...", "hostname": "...", "os_info": "...", "agent_manager_status": "RUNNING"/"STOPPED"/"UNKNOWN", "last_seen": "...", "monitors": [ /* monitor summaries */ ] }`
    *   **认证:** 需要

**3. 监控任务 (Monitors)** - 代表一个在特定机器上运行的 eBPF 监控实例

*   `POST /api/v1/machines/{machine_id}/monitors`
    *   **描述:** 在指定机器上启动一个新的监控任务。后端会通过控制通道将命令下发给对应的 Agent Manager。
    *   **路径参数:** `machine_id`
    *   **请求体:**
        ```json
        {
          "monitor_type": "CUDA" | "PYTHON" | "SYSCALL" | "OLLAMA", // 监控类型
          "target": { // 监控目标详情
            "pid": 12345, // 目标进程 PID (可选, 某些类型需要)
            "comm": "python", // 目标进程名 (可选, 辅助或用于过滤)
            "library_path": "/usr/local/cuda/lib64/libcudart.so", // uprobe 需要的库路径 (可选)
            "binary_path": "/usr/bin/python3.9" // uprobe 可能需要的二进制路径 (可选)
            // ... 其他特定类型的目标参数, e.g., syscall 过滤规则
          },
          "config": { // 监控配置 (可选)
            "sampling_rate": 1.0, // 采样率 (如果支持)
            "enable_profiling": true // 是否开启某些深度分析功能
            // ... 其他特定类型的配置
          },
          "display_name": "Monitor Ollama GPU usage" // 用户定义的显示名称 (可选)
        }
        ```
    *   **响应 (202 Accepted or 201 Created):** 返回表示任务已接受或已创建的监控资源。可能包含一个 `task_id` 用于后续查询状态，或直接返回 monitor 详情。
        ```json
        // 201 Created example
        {
          "id": "monitor_uuid_123", // 后端生成的唯一 ID
          "display_name": "Monitor Ollama GPU usage",
          "machine_id": "machine_uuid_1",
          "monitor_type": "CUDA",
          "target": { ... },
          "config": { ... },
          "status": "STARTING" | "RUNNING" | "FAILED_TO_START", // 初始状态
          "create_time": "..."
        }
        ```
    *   **认证:** 需要

*   `GET /api/v1/machines/{machine_id}/monitors`
    *   **描述:** 列出指定机器上当前活动或历史的监控任务。
    *   **路径参数:** `machine_id`
    *   **查询参数:** `filter` (e.g., `status=='RUNNING'`), `page_size`, `page_token`, `order_by`
    *   **响应 (200 OK):** `{ "monitors": [ { "id": "...", "display_name": "...", "monitor_type": "...", "status": "RUNNING", "start_time": "..." }, ... ], "next_page_token": "..." }`
    *   **认证:** 需要

*   `GET /api/v1/machines/{machine_id}/monitors/{monitor_id}`
    *   **描述:** 获取特定监控任务的详细信息。
    *   **路径参数:** `machine_id`, `monitor_id`
    *   **响应 (200 OK):** `{ "id": "...", "display_name": "...", "machine_id": "...", "monitor_type": "...", "target": {...}, "config": {...}, "status": "RUNNING" | "STOPPED" | "ERROR", "start_time": "...", "stop_time": "...", "error_message": "..." }`
    *   **认证:** 需要

*   `DELETE /api/v1/machines/{machine_id}/monitors/{monitor_id}`
    *   **描述:** 停止指定的监控任务。后端会通过控制通道将命令下发给对应的 Agent Manager。
    *   **路径参数:** `machine_id`, `monitor_id`
    *   **响应 (202 Accepted or 204 No Content):** 表示停止请求已被接受或已完成。
    *   **认证:** 需要

**4. 事件数据 (Events)** - 存储在 PostgreSQL 中的详细记录

*   `GET /api/v1/events`
    *   **描述:** 查询存储的事件数据。
    *   **查询参数:**
        *   `filter` (非常重要, e.g., `machine_id=='...' AND timestamp > '...' AND timestamp < '...' AND event_type=='CUDA_MEMCPY' AND pid==12345`)
        *   `page_size`, `page_token`
        *   `order_by` (e.g., `timestamp desc`)
        *   `fields` (可选, 选择返回哪些字段, e.g., `pid,comm,event_data.duration_ns`)
    *   **响应 (200 OK):** `{ "events": [ { "timestamp": "...", "machine_id": "...", "monitor_id": "...", "pid": 12345, "comm": "ollama", "event_type": "CUDA_MEMCPY", "event_data": { /* 特定事件的数据 */ } }, ... ], "next_page_token": "..." }`
    *   **认证:** 需要

**5. 时序指标 (Metrics)** - 存储在 InfluxDB 中的聚合数据

*   `GET /api/v1/metrics/summary` (示例 - 具体 API 取决于需求)
    *   **描述:** 获取关键指标的摘要或最新值 (供 React UI 使用, Grafana 可直连 InfluxDB)。
    *   **查询参数:** `filter` (e.g., `machine_id=='...' AND time_range='1h'`), `metrics` (e.g., `gpu_memory_usage,cpu_utilization`)
    *   **响应 (200 OK):** `{ "summary": { "machine_id": "...", "time_range": "...", "metrics": { "gpu_memory_usage": { "current": 8192, "avg_1h": 7500 }, ... } } }`
    *   **认证:** 需要

**6. 追踪数据 (Traces)** - Perfetto 格式

*   `POST /api/v1/traces:generate`
    *   **描述:** 请求后端基于存储的事件生成 Perfetto 追踪文件。这是一个异步操作。
    *   **请求体:** `{ "filter": "machine_id=='...' AND timestamp > '...' AND timestamp < '...' AND (pid==12345 OR monitor_id=='...')", "trace_name": "ollama_gpu_trace_1" }`
    *   **响应 (202 Accepted):** `{ "trace_job_id": "job_uuid_abc" }`

*   `GET /api/v1/traces/jobs/{trace_job_id}`
    *   **描述:** 查询追踪文件生成任务的状态。
    *   **路径参数:** `trace_job_id`
    *   **响应 (200 OK):**
        *   任务进行中: `{ "status": "PROCESSING", "progress": 0.5 }`
        *   任务完成: `{ "status": "COMPLETED", "trace_id": "trace_uuid_xyz", "download_url": "/api/v1/traces/files/trace_uuid_xyz" }` (或者直接包含 trace_id)
        *   任务失败: `{ "status": "FAILED", "error_message": "..." }`
    *   **认证:** 需要

*   `GET /api/v1/traces/{trace_id}` (或者 `GET /api/v1/traces/files/{trace_id}`)
    *   **描述:** 获取（下载）已生成的 Perfetto 追踪文件 (.json 或 .perfetto-trace)。
    *   **路径参数:** `trace_id`
    *   **响应 (200 OK):** 文件内容 (Content-Type: `application/json` 或 `application/octet-stream`, `Content-Disposition: attachment; filename="trace_name.json"`)
    *   **认证:** 需要

---

## 二、 Go Agent Manager 控制 API (内部接口)

**设计原则:**

*   **通信方式:** ZeroMQ (推荐使用 REQ/REP 或 ROUTER/DEALER 模式) 或 gRPC。
*   **数据格式:** JSON (如果用 ZMQ) 或 Protobuf (如果用 gRPC)。
*   **安全性:** **必须**启用加密和认证（例如 ZMQ CURVE 或 gRPC mTLS）。Agent Manager 以 root 运行，此通道必须严格保护。
*   **交互模式:** Backend (Client) 发送命令，Agent Manager (Server) 执行并响应。
*   **幂等性:** 某些操作（如 STOP）应设计为幂等的。

---

**控制命令与响应 (以 JSON over ZMQ 为例):**

**1. START_MONITOR 命令**

*   **Backend 发送 (Request):**
    ```json
    {
      "command": "START_MONITOR",
      "request_id": "backend_req_uuid_1", // 用于追踪请求
      "monitor_id": "monitor_uuid_123", // 由 Backend 生成，需传递给 C Agent
      "monitor_type": "CUDA",
      "target": {
        "pid": 12345,
        "library_path": "/usr/local/cuda/lib64/libcudart.so"
      },
      "config": {
        "sampling_rate": 1.0
      },
      "backend_zmq_data_endpoint": "tcp://backend-server:5556", // C Agent 推送数据的地址
      "machine_id": "machine_uuid_1" // Agent Manager 可用于自检
    }
    ```
*   **Agent Manager 响应 (Reply):**
    *   **成功:**
        ```json
        {
          "request_id": "backend_req_uuid_1",
          "status": "SUCCESS",
          "monitor_id": "monitor_uuid_123",
          "message": "Monitor started successfully.",
          "agent_pid": 54321 // C Agent 进程的 PID (可选, 用于管理)
        }
        ```
    *   **失败:**
        ```json
        {
          "request_id": "backend_req_uuid_1",
          "status": "FAILURE",
          "monitor_id": "monitor_uuid_123",
          "error_code": "TARGET_NOT_FOUND" | "PERMISSION_DENIED" | "AGENT_START_FAILED" | "INVALID_CONFIG",
          "message": "Failed to find target PID 12345."
        }
        ```

**2. STOP_MONITOR 命令**

*   **Backend 发送 (Request):**
    ```json
    {
      "command": "STOP_MONITOR",
      "request_id": "backend_req_uuid_2",
      "monitor_id": "monitor_uuid_123" // 要停止的监控任务 ID
    }
    ```
*   **Agent Manager 响应 (Reply):**
    *   **成功:**
        ```json
        {
          "request_id": "backend_req_uuid_2",
          "status": "SUCCESS",
          "monitor_id": "monitor_uuid_123",
          "message": "Monitor stopped successfully."
        }
        ```
    *   **失败 (例如，找不到对应的 Monitor 或 Agent 进程):**
        ```json
        {
          "request_id": "backend_req_uuid_2",
          "status": "FAILURE",
          "monitor_id": "monitor_uuid_123",
          "error_code": "MONITOR_NOT_FOUND" | "AGENT_STOP_FAILED",
          "message": "Monitor with id 'monitor_uuid_123' not found or already stopped."
        }
        ```

**3. GET_STATUS 命令 (可选)**

*   **Backend 发送 (Request):**
    ```json
    {
      "command": "GET_STATUS",
      "request_id": "backend_req_uuid_3",
      "monitor_id": "monitor_uuid_123" // 可选，获取特定 monitor 状态，否则获取 Agent Manager 状态
    }
    ```
*   **Agent Manager 响应 (Reply):**
    *   **获取 Agent Manager 状态:**
        ```json
        {
          "request_id": "backend_req_uuid_3",
          "status": "SUCCESS",
          "agent_manager_status": "RUNNING",
          "machine_id": "machine_uuid_1",
          "active_monitors": [
            { "monitor_id": "monitor_uuid_123", "agent_pid": 54321, "start_time": "...", "monitor_type": "CUDA" },
            { "monitor_id": "monitor_uuid_456", "agent_pid": 54322, "start_time": "...", "monitor_type": "PYTHON" }
          ],
          "version": "..." // Agent Manager 版本
        }
        ```
    *   **获取特定 Monitor 状态:**
        ```json
        {
          "request_id": "backend_req_uuid_3",
          "status": "SUCCESS",
          "monitor_id": "monitor_uuid_123",
          "monitor_status": "RUNNING" | "STOPPED" | "ERROR",
          "agent_pid": 54321, // 如果在运行
          "start_time": "...",
          "error_message": "..." // 如果是 ERROR 状态
        }
        ```
    *   **失败 (e.g., Monitor not found):**
        ```json
        {
          "request_id": "backend_req_uuid_3",
          "status": "FAILURE",
          "monitor_id": "monitor_uuid_123",
          "error_code": "MONITOR_NOT_FOUND",
          "message": "Monitor not found."
        }
        ```

**Agent Manager 的职责:**

1.  监听来自 Backend 的安全连接 (ZMQ/gRPC)。
2.  验证收到的命令。
3.  解析 `START_MONITOR` 命令：
    *   定位 C Agent 的可执行文件。
    *   准备 C Agent 的命令行参数或环境变量，必须包含：
        *   `--monitor-id <monitor_uuid_123>`
        *   `--monitor-type <CUDA|PYTHON|...>`
        *   `--target-pid <pid>` (如果需要)
        *   `--target-library <path>` (如果需要)
        *   `--backend-zmq-endpoint <tcp://backend:5556>`
        *   `--machine-id <machine_uuid_1>`
        *   其他配置参数。
    *   使用 `os/exec` 以 **root** 权限启动 C Agent 进程。
    *   记录下 `monitor_id` 与启动的 C Agent 进程 PID 的映射关系。
    *   向 Backend 回复成功或失败。
4.  解析 `STOP_MONITOR` 命令：
    *   根据 `monitor_id` 查找对应的 C Agent 进程 PID。
    *   如果找到，向该 PID 发送 SIGTERM 信号（或 SIGKILL 作为后备）。
    *   等待进程退出（可选，带超时）。
    *   清理 `monitor_id` -> PID 的映射。
    *   向 Backend 回复成功或失败。
5.  (可选) 解析 `GET_STATUS` 命令并根据内部状态回复。

这个设计划分了用户/前端与 Backend 交互的 RESTful API，以及 Backend 与 Agent Manager 之间更底层的、面向命令的控制 API，并强调了控制通道的安全性。