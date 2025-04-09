package postgres

/*
好的，我们现在基于之前的分类，并结合你提供的 Go `Processor` 代码生成的 `eventData` 结构，来设计**按大类分组**的 TimescaleDB 表结构。

这种方案的核心是为每个事件大类创建一个 Hypertable，表内包含该类别下所有具体事件类型的字段超集，并用一个 `event_subtype` 字段来区分它们。

**1. 操作系统/内核事件表 (`events_os`)**

*   **包含事件 Topics:** `vfs_open`, `syscalls`, `sched`, `execv`
*   **对应 `eventData` 关键字段:** `topic` (映射到 `event_subtype`), `timestamp`, `pid`, `comm`, `cmdline`, `machineid`, `filename`, `syscall`, `cpu`, `type` (sched), `ppid`, `ppid_comm`, `ppid_cmdline`, `args`
*   **SQL 定义:**

```sql
CREATE TABLE events_os (
    -- 公共字段 (来自 eventData)
    ts TIMESTAMPTZ NOT NULL,         -- 事件时间戳 (从纳秒转换)
    machine_id TEXT NOT NULL,         -- 机器 ID
    event_subtype TEXT NOT NULL,      -- 子事件类型: 'vfs_open', 'syscalls', 'sched', 'execv'
    pid INT,                          -- 进程 ID
    comm TEXT,                        -- 进程名
    cmdline TEXT,                     -- 完整命令行

    -- vfs_open 特定字段 (来自 eventData["filename"])
    vfs_filename TEXT,

    -- syscalls 特定字段 (来自 eventData["syscall"])
    syscall_name TEXT,                -- 系统调用名 (重命名以清晰)

    -- sched 特定字段 (来自 eventData["cpu"], eventData["type"])
    cpu INT,                          -- 调度 CPU ID
    sched_type TEXT,                  -- 调度类型 ('switch_in', 'switch_out', 'unknown') - 使用 TEXT

    -- execv 特定字段 (来自 eventData["ppid*"], eventData["filename"], eventData["args"])
    ppid INT,
    ppid_comm TEXT,
    ppid_cmdline TEXT,
    exec_filename TEXT,               -- (重命名自 eventData["filename"] 以区分 vfs_open)
    exec_args TEXT                    -- (重命名自 eventData["args"])
);

-- 转换为 Hypertable
SELECT create_hypertable('events_os', 'ts', chunk_time_interval => INTERVAL '1 day');

-- 索引策略
CREATE INDEX ix_events_os_machine_id_ts ON events_os (machine_id, ts DESC);
CREATE INDEX ix_events_os_subtype_ts ON events_os (event_subtype, ts DESC); -- 重要：用于快速过滤子类型
CREATE INDEX ix_events_os_pid_ts ON events_os (pid, ts DESC);
```

**2. CUDA/GPU 通用事件表 (`events_cuda`)**

*   **包含事件 Topics:** `cudaMalloc`, `cudaFree`, `cudaLaunchKernel`, `cudaMemcpy`, `cudaDeviceSynchronize`
*   **对应 `eventData` 关键字段:** `topic` (->`event_subtype`), `timestamp`, `pid`, `comm`, `cmdline`, `machineid`, `operation` (cudaOperation), `ptr` (cudaMalloc/Free), `size` (Malloc/Memcpy), `retval`, `func_ptr`, `symbol_*`, `src`, `dst`, `kind`, `type` (Memcpy), `duration_ns` (Sync)
*   **SQL 定义:**

```sql
CREATE TABLE events_cuda (
    -- 公共字段 (来自 eventData)
    ts TIMESTAMPTZ NOT NULL,
    machine_id TEXT NOT NULL,
    event_subtype TEXT NOT NULL,      -- 'cudaMalloc', 'cudaFree', 'cudaLaunchKernel', 'cudaMemcpy', 'cudaDeviceSynchronize'
    pid INT,
    comm TEXT,
    cmdline TEXT,

    -- 操作标识 (来自 eventData["operation"])
    operation TEXT,              -- 'cudaMalloc', 'cudaFree', 'cudaLaunchKernel', 'cudaMemcpy', 'cudaDeviceSynchronize'

    -- cudaMalloc 特定/共享字段
    cuda_ptr BIGINT,                  -- 设备指针 (合并 Malloc的ptr, Free的ptr) - Use BIGINT for uint64
    cuda_size BIGINT,                 -- 大小 (合并 Malloc的size, Memcpy的size) - Use BIGINT for uint64
    cuda_retval INT,                  -- 返回值 (仅 Malloc)

    -- cudaLaunchKernel 特定字段
    cuda_func_ptr BIGINT,             -- 核函数入口指针 - Use BIGINT for uint64
    cuda_symbol_name TEXT,            -- 解析的符号名
    cuda_symbol_file TEXT,            -- 符号所在文件
    cuda_symbol_offset BIGINT,        -- 符号偏移 - Use BIGINT for uint64/offset
    cuda_symbol_sourcefile TEXT,      -- 源码位置 (e.g., "file.cu:123")

    -- cudaMemcpy 特定字段
    cuda_memcpy_src BIGINT,           -- 源地址 (重命名自 eventData["src"]) - Use BIGINT for uint64
    cuda_memcpy_dst BIGINT,           -- 目标地址 (重命名自 eventData["dst"]) - Use BIGINT for uint64
    cuda_memcpy_kind INT,             -- 原始拷贝类型 (来自 eventData["kind"])
    cuda_memcpy_type TEXT,            -- 人类可读拷贝类型 (来自 eventData["type"])

    -- cudaDeviceSynchronize 特定字段
    cuda_sync_duration_ns BIGINT      -- 同步耗时 (重命名自 eventData["duration_ns"])
);

-- 转换为 Hypertable
SELECT create_hypertable('events_cuda', 'ts', chunk_time_interval => INTERVAL '1 day');

-- 索引策略
CREATE INDEX ix_events_cuda_machine_id_ts ON events_cuda (machine_id, ts DESC);
CREATE INDEX ix_events_cuda_subtype_ts ON events_cuda (event_subtype, ts DESC); -- 重要
CREATE INDEX ix_events_cuda_pid_ts ON events_cuda (pid, ts DESC);
CREATE INDEX ix_events_operation_ts ON events_cuda (operation, ts DESC); -- 可替代 subtype 或一起使用
```

**3. GGML 特定库事件表 (`events_ggml`)**

*   **包含事件 Topics:** `ggml_cuda`, `ggml_graph_compute`, `ggml_base`
*   **对应 `eventData` 关键字段:** `topic` (->`event_subtype`), `timestamp`, `pid`, `comm`, `cmdline`, `machineid`, `func_name`, `duration_ns` (ggml_cuda), `operation` (ggml_graph_compute, ggml_base), `graph_*`, `cost_ns`, `size`, `ptr` (ggml_base)
*   **SQL 定义:**

```sql
CREATE TABLE events_ggml (
    -- 公共字段 (来自 eventData)
    ts TIMESTAMPTZ NOT NULL,
    machine_id TEXT NOT NULL,
    event_subtype TEXT NOT NULL,      -- 'ggml_cuda', 'ggml_graph_compute', 'ggml_base'
    pid INT,
    comm TEXT,
    cmdline TEXT,

    -- 操作标识 (来自 eventData["operation"] 或 subtype)
    operation TEXT,
-- 'ggml_aligned_malloc', 'ggml_aligned_free' for ggml_base
-- 'ggml_graph_compute'   for   ggml_cpu
-- 'ggml_cuda_op_mul_mat_vec_q' , 'ggml_cuda_op_mul_mat_q'  for ggml_cuda

    -- ggml_cuda 特定字段 (来自 eventData["func_name"], eventData["duration_ns"])
    ggml_cuda_func_name TEXT,         -- (重命名 func_name)
    ggml_cuda_duration_ns BIGINT,     -- (重命名 duration_ns)

    -- ggml_graph_compute 特定字段 (来自 eventData["graph_*"], eventData["cost_ns"], eventData["graph_order"])
    ggml_graph_size INT,
    ggml_graph_nodes INT,
    ggml_graph_leafs INT,
    ggml_graph_order TEXT,            -- 图计算顺序 ('LEFT_TO_RIGHT', etc.) - Use TEXT
    ggml_cost_ns BIGINT,

    -- ggml_base 特定字段 (来自 eventData["size"], eventData["ptr"])
    ggml_mem_size BIGINT,             -- (重命名 size) - Use BIGINT for uint64
    ggml_mem_ptr BIGINT               -- (重命名 ptr) - Use BIGINT for uint64
);

-- 转换为 Hypertable
SELECT create_hypertable('events_ggml', 'ts', chunk_time_interval => INTERVAL '1 day');

-- 索引策略
CREATE INDEX ix_events_ggml_machine_id_ts ON events_ggml (machine_id, ts DESC);
CREATE INDEX ix_events_ggml_subtype_ts ON events_ggml (event_subtype, ts DESC); -- 重要
CREATE INDEX ix_events_ggml_pid_ts ON events_ggml (pid, ts DESC);
CREATE INDEX ix_events_operation_ts ON events_ggml (operation, ts DESC);
```

**4. 应用日志事件表 (`events_app_log`)**

*   **包含事件 Topics:** `llamaLog` (以及未来可能的其他应用日志)
*   **对应 `eventData` 关键字段:** `topic` (->`event_subtype`), `timestamp`, `pid`, `comm`, `cmdline`, `machineid`, `text`
*   **SQL 定义:**

```sql
CREATE TABLE events_app_log (
    -- 公共字段 (来自 eventData)
    ts TIMESTAMPTZ NOT NULL,
    machine_id TEXT NOT NULL,
    event_subtype TEXT NOT NULL,      -- 'llamaLog', 'otherAppXLog', etc.
    pid INT,
    comm TEXT,
    cmdline TEXT,

    -- llamaLog 特定字段 (来自 eventData["text"])
    log_text TEXT

    -- 如果未来有其他应用日志，在此处添加对应字段
    -- other_app_field_1 TEXT,
    -- other_app_field_2 INT,
);

-- 转换为 Hypertable
SELECT create_hypertable('events_app_log', 'ts', chunk_time_interval => INTERVAL '1 day'); -- 应用日志分区可不同

-- 索引策略
CREATE INDEX ix_events_app_log_machine_id_ts ON events_app_log (machine_id, ts DESC);
CREATE INDEX ix_events_app_log_subtype_ts ON events_app_log (event_subtype, ts DESC); -- 重要
CREATE INDEX ix_events_app_log_pid_ts ON events_app_log (pid, ts DESC);
```

**设计说明:**

1.  **数据源:** 设计严格基于 Go `Processor` 生成并放入 Redis 的 `eventData` map 结构。
2.  **分组依据:** 沿用之前的四大类别：OS、CUDA、GGML、AppLog。
3.  **核心列:** `ts` (转换后), `machine_id`, `event_subtype` (原 `topic`), `pid`, `comm`, `cmdline` 是所有表的基础。
4.  **字段超集:** 每个表包含其类别下所有事件的字段，不存在于某个具体事件的字段将为 `NULL`。
5.  **命名:** 尽量保持与 `eventData` key 一致，但在必要时添加前缀 (`cuda_`, `ggml_`) 或重命名 (`vfs_filename` vs `exec_filename`) 以消除歧义或提高可读性。
6.  **数据类型:** `timestamp` -> `TIMESTAMPTZ`, `string` -> `TEXT`, `int32`/`int` -> `INT`, `int64`/`uint64` -> `BIGINT`。Go 代码中添加的字符串类型（如 `sched_type`, `cuda_memcpy_type`）直接用 `TEXT`。
7.  **`event_subtype` 索引:** 这是分组表设计的关键索引，用于高效过滤特定类型的事件。
8.  **部分索引 (Partial Indexes):** 对于只在特定 `event_subtype` 下有意义且经常查询的字段，建议使用部分索引以提高效率和减少存储。

**插入流程 (概念):**

1.  从 Redis Stream 读取 `eventData` map。
2.  根据 map 中的 `topic` 字段确定目标表（`events_os`, `events_cuda`, `events_ggml`, `events_app_log`）。
3.  将 `topic` 字段的值赋给目标表的 `event_subtype` 列。
4.  将 `timestamp` (nanoseconds) 转换为 `TIMESTAMPTZ` (e.g., `time.Unix(0, ns).UTC()`) 并赋给 `ts` 列。
5.  将 `machineid`, `pid`, `comm`, `cmdline` 赋给对应列。
6.  根据 `topic` (即 `event_subtype`)，将 map 中其他相关的 key-value 对赋给目标表中对应的列。不存在的 key 对应的列将插入 `NULL`。
7.  使用批量插入 (`CopyFrom`) 将准备好的数据写入对应的 TimescaleDB 表。
*/

import (
	"context"
	"database/sql"

	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// --- SQL Statements for Schema Initialization ---

const (
	// Enable TimescaleDB Extension
	enableTimescaleDBExtensionSQL = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`

	// --- events_os ---
	createEventsOsTableSQL = `
CREATE TABLE events_os (
    ts TIMESTAMPTZ NOT NULL,
    machine_id TEXT NOT NULL,
    event_subtype TEXT NOT NULL,
    pid INT NOT NULL,
    comm TEXT,
    cmdline TEXT,
    vfs_filename TEXT,
    syscall_name TEXT,
    cpu INT,
    sched_type TEXT,
    ppid INT,
    ppid_comm TEXT,
    ppid_cmdline TEXT,
    exec_filename TEXT,
    exec_args TEXT
);`
	createEventsOsHypertableSQL   = `SELECT create_hypertable('events_os', 'ts', chunk_time_interval => INTERVAL '1 day');`
	createEventsOsIdxMachineIDSQL = `CREATE INDEX IF NOT EXISTS ix_events_os_machine_id_ts ON events_os (machine_id, ts DESC);`
	createEventsOsIdxSubtypeSQL   = `CREATE INDEX IF NOT EXISTS ix_events_os_subtype_ts ON events_os (event_subtype, ts DESC);`
	createEventsOsIdxPidSQL       = `CREATE INDEX IF NOT EXISTS ix_events_os_pid_ts ON events_os (pid, ts DESC);`

	// --- events_cuda ---
	createEventsCudaTableSQL = `
CREATE TABLE events_cuda (
    ts TIMESTAMPTZ NOT NULL,
    machine_id TEXT NOT NULL,
    event_subtype TEXT NOT NULL,
    pid INT NOT NULL,
    comm TEXT,
    cmdline TEXT,
    operation TEXT,
    cuda_ptr BIGINT,
    cuda_size BIGINT,
    cuda_retval INT,
    cuda_func_ptr BIGINT,
    cuda_symbol_name TEXT,
    cuda_symbol_file TEXT,
    cuda_symbol_offset BIGINT,
    cuda_symbol_sourcefile TEXT,
    cuda_memcpy_src BIGINT,
    cuda_memcpy_dst BIGINT,
    cuda_memcpy_kind INT,
    cuda_memcpy_type TEXT,
    cuda_sync_duration_ns BIGINT
);`
	createEventsCudaHypertableSQL   = `SELECT create_hypertable('events_cuda', 'ts', chunk_time_interval => INTERVAL '1 day');`
	createEventsCudaIdxMachineIDSQL = `CREATE INDEX IF NOT EXISTS ix_events_cuda_machine_id_ts ON events_cuda (machine_id, ts DESC);`
	createEventsCudaIdxSubtypeSQL   = `CREATE INDEX IF NOT EXISTS ix_events_cuda_subtype_ts ON events_cuda (event_subtype, ts DESC);`
	createEventsCudaIdxPidSQL       = `CREATE INDEX IF NOT EXISTS ix_events_cuda_pid_ts ON events_cuda (pid, ts DESC);`
	createEventsCudaIdxOperationSQL = `CREATE INDEX IF NOT EXISTS ix_events_cuda_operation_ts ON events_cuda (operation, ts DESC);` // Corrected index name

	// --- events_ggml ---
	createEventsGgmlTableSQL = `
CREATE TABLE events_ggml (
    ts TIMESTAMPTZ NOT NULL,
    machine_id TEXT NOT NULL,
    event_subtype TEXT NOT NULL,
    pid INT NOT NULL,
    comm TEXT,
    cmdline TEXT,
    operation TEXT,
    ggml_cuda_func_name TEXT,
    ggml_cuda_duration_ns BIGINT,
    ggml_graph_size INT,
    ggml_graph_nodes INT,
    ggml_graph_leafs INT,
    ggml_graph_order TEXT,
    ggml_cost_ns BIGINT,
    ggml_mem_size BIGINT,
    ggml_mem_ptr BIGINT
);`
	createEventsGgmlHypertableSQL   = `SELECT create_hypertable('events_ggml', 'ts', chunk_time_interval => INTERVAL '1 day');`
	createEventsGgmlIdxMachineIDSQL = `CREATE INDEX IF NOT EXISTS ix_events_ggml_machine_id_ts ON events_ggml (machine_id, ts DESC);`
	createEventsGgmlIdxSubtypeSQL   = `CREATE INDEX IF NOT EXISTS ix_events_ggml_subtype_ts ON events_ggml (event_subtype, ts DESC);`
	createEventsGgmlIdxPidSQL       = `CREATE INDEX IF NOT EXISTS ix_events_ggml_pid_ts ON events_ggml (pid, ts DESC);`
	createEventsGgmlIdxOperationSQL = `CREATE INDEX IF NOT EXISTS ix_events_ggml_operation_ts ON events_ggml (operation, ts DESC);` // Corrected index name

	// --- events_app_log ---
	createEventsAppLogTableSQL = `
CREATE TABLE events_app_log (
    ts TIMESTAMPTZ NOT NULL,
    machine_id TEXT NOT NULL,
    event_subtype TEXT NOT NULL,
    pid INT NOT NULL,
    comm TEXT,
    cmdline TEXT,
    log_text TEXT
);`
	createEventsAppLogHypertableSQL   = `SELECT create_hypertable('events_app_log', 'ts', chunk_time_interval => INTERVAL '1 day');` // Adjusted interval for consistency, change if needed
	createEventsAppLogIdxMachineIDSQL = `CREATE INDEX IF NOT EXISTS ix_events_app_log_machine_id_ts ON events_app_log (machine_id, ts DESC);`
	createEventsAppLogIdxSubtypeSQL   = `CREATE INDEX IF NOT EXISTS ix_events_app_log_subtype_ts ON events_app_log (event_subtype, ts DESC);`
	createEventsAppLogIdxPidSQL       = `CREATE INDEX IF NOT EXISTS ix_events_app_log_pid_ts ON events_app_log (pid, ts DESC);`
)

// InitializeTSDBSchema ensures the required TimescaleDB extension and tables exist.
// It creates them idempotently if they are missing.
func InitializeTSDBSchema(ctx context.Context, db *sqlx.DB) error {
	log.Println("开始初始化 TimescaleDB schema...")

	// 1. Enable TimescaleDB Extension
	log.Println("确保 TimescaleDB extension 已启用...")
	if _, err := db.ExecContext(ctx, enableTimescaleDBExtensionSQL); err != nil {
		return fmt.Errorf("启用 TimescaleDB extension 失败: %w", err)
	}
	log.Println("TimescaleDB extension 已启用.")

	// 2. Initialize each table group
	if err := initializeTableGroup(ctx, db, "events_os", createEventsOsTableSQL, createEventsOsHypertableSQL, []string{
		createEventsOsIdxMachineIDSQL,
		createEventsOsIdxSubtypeSQL,
		createEventsOsIdxPidSQL,
	}); err != nil {
		return err
	}

	if err := initializeTableGroup(ctx, db, "events_cuda", createEventsCudaTableSQL, createEventsCudaHypertableSQL, []string{
		createEventsCudaIdxMachineIDSQL,
		createEventsCudaIdxSubtypeSQL,
		createEventsCudaIdxPidSQL,
		createEventsCudaIdxOperationSQL,
	}); err != nil {
		return err
	}

	if err := initializeTableGroup(ctx, db, "events_ggml", createEventsGgmlTableSQL, createEventsGgmlHypertableSQL, []string{
		createEventsGgmlIdxMachineIDSQL,
		createEventsGgmlIdxSubtypeSQL,
		createEventsGgmlIdxPidSQL,
		createEventsGgmlIdxOperationSQL,
	}); err != nil {
		return err
	}

	if err := initializeTableGroup(ctx, db, "events_app_log", createEventsAppLogTableSQL, createEventsAppLogHypertableSQL, []string{
		createEventsAppLogIdxMachineIDSQL,
		createEventsAppLogIdxSubtypeSQL,
		createEventsAppLogIdxPidSQL,
	}); err != nil {
		return err
	}

	log.Println("数据库 schema 初始化完成.")
	return nil
}

// initializeTableGroup checks and creates a specific table, its hypertable conversion, and indexes.
func initializeTableGroup(ctx context.Context, db *sqlx.DB, tableName, createTableSQL, createHypertableSQL string, createIndexSQLs []string) error {
	log.Printf("检查表 '%s'...", tableName)
	exists, err := tableExists(ctx, db, "public", tableName) // Assuming 'public' schema
	if err != nil {
		return fmt.Errorf("检查表 '%s' 是否存在时出错: %w", tableName, err)
	}

	if !exists {
		log.Printf("表 '%s' 不存在，正在创建...", tableName)

		// Begin transaction for atomic creation
		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("为表 '%s' 启动事务失败: %w", tableName, err)
		}
		// Defer rollback in case of error
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p) // re-throw panic after Rollback
			} else if err != nil {
				log.Printf("事务因错误回滚: %v", err)
				tx.Rollback() // err is non-nil; don't change it
			} else {
				err = tx.Commit() // err is nil; if Commit fails, update err
				if err != nil {
					log.Printf("事务提交失败: %v", err)
				}
			}
		}()

		// Create the regular table
		log.Printf(" -> 创建表 '%s'...", tableName)
		if _, err = tx.ExecContext(ctx, createTableSQL); err != nil {
			return fmt.Errorf("创建表 '%s' 失败: %w", tableName, err)
		}

		// Convert to Hypertable
		// Note: create_hypertable might return notices, not errors, if already a hypertable,
		// but since we checked table existence, this should run only once.
		// It also doesn't work well inside transactions in some older Timescale versions,
		// but generally recommended now. If issues arise, move it outside the tx.
		log.Printf(" -> 转换 '%s' 为 Hypertable...", tableName)
		// We expect create_hypertable to potentially return rows (e.g., the hypertable name)
		// Using QueryRowContext and Scan is safer than ExecContext here.
		// We don't actually need the result, just check for errors. ExecContext is fine.
		if _, err = tx.ExecContext(ctx, createHypertableSQL); err != nil {
			// Check if error is "already a hypertable" - PQTSError potentially needed for specific code
			pqErr, ok := err.(*pq.Error)
			// Specific error code for "already a hypertable" might vary or not exist,
			// often it just returns success or a notice. If it's a real error, report it.
			// This simple check might suffice if the extension handles it gracefully.
			// If create_hypertable reliably errors *outside* the IF NOT EXISTS block,
			// this check needs refinement. Given we check table existence first, it *shouldn't* error here
			// unless something else creates it concurrently.
			if !ok || pqErr.Code != "SOME_SPECIFIC_CODE_FOR_ALREADY_HYPERTABLE" { // Replace with actual code if known
				return fmt.Errorf("转换 '%s' 为 Hypertable 失败: %w", tableName, err)
			}
			log.Printf(" -> 注意：表 '%s' 似乎已经是 Hypertable (错误: %v)", tableName, err)
			// Reset err to nil as this specific "error" is acceptable here
			err = nil
		}

		// Create Indexes (using CREATE INDEX IF NOT EXISTS for idempotency)
		log.Printf(" -> 为 '%s' 创建索引...", tableName)
		for _, indexSQL := range createIndexSQLs {
			if _, err = tx.ExecContext(ctx, indexSQL); err != nil {
				return fmt.Errorf("为表 '%s' 创建索引失败 (%s): %w", tableName, indexSQL, err)
			}
		}

		log.Printf("表 '%s' 创建成功.", tableName)

	} else {
		log.Printf("表 '%s' 已存在.", tableName)
	}
	return nil
}

// tableExists checks if a table exists in a given schema.
func tableExists(ctx context.Context, db *sqlx.DB, schemaName, tableName string) (bool, error) {
	query := `SELECT EXISTS (
        SELECT FROM information_schema.tables
        WHERE table_schema = $1 AND table_name = $2
    );`
	var exists bool
	err := db.GetContext(ctx, &exists, query, schemaName, tableName)
	if err != nil && err != sql.ErrNoRows { // ErrNoRows should not happen with EXISTS
		return false, err
	}
	return exists, nil
}
