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
	cmdlineCache              map[int]string
	cmdlineCacheLock          sync.Mutex

	lastClearCommCacheTime time.Time
	commCache              map[int]string
	commCacheLock          sync.Mutex
)

func ClearCmdlineCache() {
	if cmdlineCache == nil {
		cmdlineCache = make(map[int]string)
	}
	clear(cmdlineCache)
	lastClearCmdlineCacheTime = time.Now()
}

func ClearCommCache() {
	if commCache == nil {
		commCache = make(map[int]string)
	}
	clear(commCache)
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
		cmdlineCacheLock.Lock()
		ClearCmdlineCache()
		cmdlineCacheLock.Unlock()
	}

	cmdlineCacheLock.Lock()
	cachedValue, found := cmdlineCache[pid]
	cmdlineCacheLock.Unlock()
	if found {
		return cachedValue, nil
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

	cmdlineCacheLock.Lock()
	cmdlineCache[pid] = cmdlineStr
	cmdlineCacheLock.Unlock()

	return cmdlineStr, nil
}

func GetComm(pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid pid: %d", pid)
	}

	if time.Since(lastClearCommCacheTime) > 5*time.Minute {
		commCacheLock.Lock()
		ClearCommCache()
		commCacheLock.Unlock()
	}

	commCacheLock.Lock()
	cachedValue, found := commCache[pid]
	commCacheLock.Unlock()
	if found {
		return cachedValue, nil
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

	commCacheLock.Lock()
	commCache[pid] = commStr
	commCacheLock.Unlock()

	return commStr, nil
}
