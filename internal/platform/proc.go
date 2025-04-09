// proc.go
package platform

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	lastClearCmdlineCacheTime time.Time
	cmdlineCache              sync.Map
	lastClearCommCacheTime    time.Time
	commCache                 sync.Map
)

func ClearCmdlineCache() {
	cmdlineCache = sync.Map{}
	lastClearCmdlineCacheTime = time.Now()
}

func ClearCommCache() {
	commCache = sync.Map{}
	lastClearCommCacheTime = time.Now()
}

// GetCmdline retrieves the command line string for a given process ID (PID)
// from /proc/[pid]/cmdline on Linux systems using sync.Map for caching.
// It replaces the null byte separators within the command line arguments with spaces.
// Results are cached in memory to reduce redundant reads from the /proc filesystem.
//
// Parameters:
//
//	pid (int): The process ID to look up.
//
// Returns:
//
//	string: The command line string with arguments separated by spaces.
//	        Returns an empty string if the cmdline file is empty or cannot be read.
//	error:  An error if the PID is invalid, the /proc entry cannot be read
//	        (e.g., process doesn't exist, permissions error), or other issues occur.
func GetCmdline(pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid pid: %d", pid)
	}

	// --- 0. 清理过期缓存 ---
	if time.Since(lastClearCmdlineCacheTime) > 5*time.Minute {
		ClearCmdlineCache()
	}

	// --- 1. 检查缓存 (使用 sync.Map.Load) ---
	// Load 原子地加载 key 对应的 value。
	// 返回值是 interface{} 类型，需要进行类型断言。
	cachedValue, found := cmdlineCache.Load(pid)
	if found {
		// 缓存命中，进行类型断言
		// 假设我们只存储 string 类型
		cmdlineStr, ok := cachedValue.(string)
		if ok {
			return cmdlineStr, nil
		}
	}

	// --- 2. 缓存未命中，从 /proc 读取 ---
	path := fmt.Sprintf("/proc/%d/cmdline", pid)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	if len(content) < 1 {
		return "", nil
	}

	cmdlineStr := string(bytes.ReplaceAll(content, []byte{'\x00'}, []byte{' '}))
	cmdlineStr = strings.TrimSuffix(cmdlineStr, " ")

	// --- 4. 更新缓存 (使用 sync.Map.Store) ---
	cmdlineCache.Store(pid, cmdlineStr)

	return cmdlineStr, nil
}

func GetComm(pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid pid: %d", pid)
	}

	if time.Since(lastClearCommCacheTime) > 5*time.Minute {
		ClearCommCache()
	}

	// --- 1. 检查缓存 (使用 sync.Map.Load) ---
	// Load 原子地加载 key 对应的 value。
	// 返回值是 interface{} 类型，需要进行类型断言。
	cachedValue, found := commCache.Load(pid)
	if found {
		// 缓存命中，进行类型断言
		// 假设我们只存储 string 类型
		commStr, ok := cachedValue.(string)
		if ok {
			return commStr, nil
		}
	}

	// --- 2. 缓存未命中，从 /proc 读取 ---
	path := fmt.Sprintf("/proc/%d/comm", pid)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}

	if len(content) < 1 {
		return "", nil
	}

	commStr := string(content)
	// --- 4. 更新缓存 (使用 sync.Map.Store) ---
	commCache.Store(pid, commStr)

	return commStr, nil
}
