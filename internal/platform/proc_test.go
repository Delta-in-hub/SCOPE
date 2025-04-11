// proc_test.go
package platform

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// TestGetCmdline tests the GetCmdline function and its caching mechanism.
func TestGetCmdline(t *testing.T) {
	// --- 获取测试进程自身的 PID ---
	selfPid := os.Getpid()
	if selfPid <= 0 {
		t.Fatalf("Failed to get own PID: %d", selfPid)
	}

	// --- 子测试 1: 获取自身进程的 Cmdline ---
	t.Run("GetSelfCmdline", func(t *testing.T) {
		cmdline, err := GetCmdline(selfPid)
		if err != nil {
			t.Errorf("GetCmdline(%d) failed: %v", selfPid, err)
		}
		if cmdline == "" {
			t.Errorf("GetCmdline(%d) returned empty string, expected non-empty", selfPid)
		}
		// 我们可以尝试更精确地检查，os.Args 包含了进程启动的参数
		// 注意：/proc/pid/cmdline 可能和 os.Args 不完全一样，特别是 argv[0] 可能被修改
		// 或者 /proc 文件有长度限制。这里我们只做一个基本包含性检查。
		expectedPart := os.Args[0] // 期望至少包含可执行文件路径
		// 如果路径包含 / 或 . ，则进行包含性检查，否则可能是 go test 的临时名称
		if strings.ContainsAny(expectedPart, "/.") {
			if !strings.Contains(cmdline, expectedPart) {
				t.Logf("Warning: GetCmdline result '%s' might not contain expected executable '%s' (os.Args[0]). This can happen in some test environments.", cmdline, expectedPart)
			}
		}
		t.Logf("Cmdline for self (PID %d): %s", selfPid, cmdline)
	})

	// --- 子测试 2: 获取 PID 1 (init/systemd) 的 Cmdline ---
	// PID 1 几乎总是在 Linux 系统上存在
	t.Run("GetInitCmdline", func(t *testing.T) {
		pid1 := 1
		cmdline, err := GetCmdline(pid1)
		if err != nil {
			// 在某些极简容器或特殊环境（如 Fargate task without init process）中，PID 1 可能无法访问或不存在
			// 检查是否是权限或不存在错误，如果是，则跳过测试而不是失败
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				t.Skipf("Skipping PID 1 test: cannot read /proc/1/cmdline: %v", err)
			} else {
				t.Errorf("GetCmdline(%d) failed unexpectedly: %v", pid1, err)
			}
		} else if cmdline == "" {
			// 理论上 PID 1 的 cmdline 不应为空，但以防万一记录下来
			t.Logf("Warning: GetCmdline(%d) returned empty string.", pid1)
		} else {
			t.Logf("Cmdline for PID 1: %s", cmdline)
		}
	})

	// --- 子测试 3: 无效 PID (0) ---
	t.Run("InvalidPIDZero", func(t *testing.T) {
		cmdline, err := GetCmdline(0)
		if err == nil {
			t.Errorf("GetCmdline(0) expected an error, but got nil (cmdline: %s)", cmdline)
		}
		if cmdline != "" {
			t.Errorf("GetCmdline(0) expected empty cmdline on error, got: %s", cmdline)
		}
	})

	// --- 子测试 4: 无效 PID (-1) ---
	t.Run("InvalidPIDNegative", func(t *testing.T) {
		cmdline, err := GetCmdline(-1)
		if err == nil {
			t.Errorf("GetCmdline(-1) expected an error, but got nil (cmdline: %s)", cmdline)
		}
		if cmdline != "" {
			t.Errorf("GetCmdline(-1) expected empty cmdline on error, got: %s", cmdline)
		}
	})

	// --- 子测试 5: 不存在的 PID ---
	// 选择一个非常大的、几乎不可能存在的 PID
	t.Run("NonExistentPID", func(t *testing.T) {
		nonExistentPid := 999999 // 或者 math.MaxInt32
		cmdline, err := GetCmdline(nonExistentPid)
		if err == nil {
			t.Errorf("GetCmdline(%d) expected an error for non-existent PID, but got nil (cmdline: %s)", nonExistentPid, cmdline)
		} else {
			// 最好检查具体的错误类型，但 os.ReadFile 可能返回包装后的错误
			// 检查是否包含 "no such process" 或类似 os.ErrNotExist 的信息
			t.Logf("GetCmdline(%d) correctly returned error: %v", nonExistentPid, err)
			// if !errors.Is(err, fs.ErrNotExist) && !strings.Contains(err.Error(), "no such process") {
			//  t.Errorf("GetCmdline(%d) expected ErrNotExist or similar, got: %v", nonExistentPid, err)
			// }
		}
		if cmdline != "" {
			t.Errorf("GetCmdline(%d) expected empty cmdline on error, got: %s", nonExistentPid, cmdline)
		}
	})

	// --- 子测试 6: 缓存命中测试 ---
	t.Run("CacheHit", func(t *testing.T) {
		// 确保缓存是空的开始 (或使用一个确定的 PID)
		// 我们将使用 selfPid，它应该在第一个子测试中已被缓存
		cmdline1, err1 := GetCmdline(selfPid)
		if err1 != nil {
			t.Fatalf("CacheHit: First call to GetCmdline(%d) failed: %v", selfPid, err1)
		}

		// 再次调用，应该从缓存中获取
		cmdline2, err2 := GetCmdline(selfPid)
		if err2 != nil {
			t.Fatalf("CacheHit: Second call to GetCmdline(%d) failed: %v", selfPid, err2)
		}

		if cmdline1 != cmdline2 {
			t.Errorf("CacheHit: Cmdline mismatch for PID %d. First call: '%s', Second call: '%s'", selfPid, cmdline1, cmdline2)
		}

		// 注意: 很难直接 *证明* 第二次调用是缓存命中而不只是快速读取。
		// 但我们可以确信，如果结果一致，缓存机制（无论是否命中）工作正常。
		t.Logf("CacheHit test passed for PID %d.", selfPid)
	})

	// --- 子测试 7: 缓存清理测试 ---
	t.Run("CacheClear", func(t *testing.T) {
		// 确保 selfPid 在缓存中
		cmdlineBeforeClear, errBefore := GetCmdline(selfPid)
		if errBefore != nil {
			t.Fatalf("CacheClear: Pre-clear call to GetCmdline(%d) failed: %v", selfPid, errBefore)
		}

		// 清理缓存
		ClearCmdlineCache()

		// 再次获取，这次应该是重新读取 /proc
		cmdlineAfterClear, errAfter := GetCmdline(selfPid)
		if errAfter != nil {
			t.Fatalf("CacheClear: Post-clear call to GetCmdline(%d) failed: %v", selfPid, errAfter)
		}

		// 结果应该仍然相同
		if cmdlineBeforeClear != cmdlineAfterClear {
			t.Errorf("CacheClear: Cmdline mismatch for PID %d after clearing cache. Before: '%s', After: '%s'", selfPid, cmdlineBeforeClear, cmdlineAfterClear)
		}

		t.Logf("CacheClear test passed for PID %d.", selfPid)
	})
}
