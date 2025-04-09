package platform

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SymbolInfo 存储解析出的符号信息
type SymbolInfo struct {
	SymbolName  string  // 函数名或符号名
	FilePath    string  // 包含该符号的可执行文件或库路径
	Offset      uintptr // 传入地址相对于文件加载基址的偏移量
	BaseAddress uintptr // 文件在内存中的加载基址
	SourceFile  string  // 源代码文件名 (如果可用)
	SourceLine  int     // 源代码行号 (如果可用)
}

// mapEntry 表示 /proc/<pid>/maps 中的一行内部结构
type mapEntry struct {
	StartAddr  uintptr
	EndAddr    uintptr
	Perms      string
	FileOffset uintptr
	Dev        string
	Inode      uint64
	Path       string
}

// 正则表达式用于解析 /proc/<pid>/maps 的行 (增强以处理空格等情况)
// 示例: 7f0c8f8f8000-7f0c8f8fa000 r-xp 00000000 103:02 12345   /usr/lib/libc.so.6
//
//	address           perms offset   dev    inode      pathname
var mapLineRegex = regexp.MustCompile(`^([0-9a-f]+)-([0-9a-f]+)\s+([rwxp\-s]+)\s+([0-9a-f]+)\s+([0-9a-f]+:[0-9a-f]+)\s+([0-9]+)\s*(.*)$`)

// Cache for addr2line path to avoid repeated lookups
var (
	addr2linePath     string
	addr2linePathErr  error
	addr2linePathOnce sync.Once
)

func findAddr2line() (string, error) {
	addr2linePathOnce.Do(func() {
		addr2linePath, addr2linePathErr = exec.LookPath("addr2line")
		if addr2linePathErr != nil {
			addr2linePathErr = errors.New("addr2line command not found in PATH. Please install binutils (or equivalent)")
		}
	})
	return addr2linePath, addr2linePathErr
}

// parseMapLine 解析 maps 文件中的单行
func parseMapLine(line string) (*mapEntry, error) {
	matches := mapLineRegex.FindStringSubmatch(line)
	if len(matches) != 8 {
		return nil, fmt.Errorf("failed to parse map line: %q", line)
	}

	startAddr, err1 := strconv.ParseUint(matches[1], 16, 64)
	endAddr, err2 := strconv.ParseUint(matches[2], 16, 64)
	fileOffset, err3 := strconv.ParseUint(matches[4], 16, 64)
	inode, err4 := strconv.ParseUint(matches[6], 10, 64) // inode is decimal

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		// Combine parsing errors if any
		var errs []string
		if err1 != nil {
			errs = append(errs, fmt.Sprintf("start addr: %v", err1))
		}
		if err2 != nil {
			errs = append(errs, fmt.Sprintf("end addr: %v", err2))
		}
		if err3 != nil {
			errs = append(errs, fmt.Sprintf("offset: %v", err3))
		}
		if err4 != nil {
			errs = append(errs, fmt.Sprintf("inode: %v", err4))
		}
		return nil, fmt.Errorf("error parsing numeric fields in map line %q: %s", line, strings.Join(errs, "; "))
	}

	// Trim path carefully, as it might be empty or contain spaces
	path := strings.TrimSpace(matches[7])

	return &mapEntry{
		StartAddr:  uintptr(startAddr),
		EndAddr:    uintptr(endAddr),
		Perms:      matches[3],
		FileOffset: uintptr(fileOffset),
		Dev:        matches[5],
		Inode:      inode,
		Path:       path,
	}, nil
}

var (
	symbolCache        sync.Map
	lastClearCacheTime time.Time
)

// FindSymbolFromPidPtr 根据 PID 和内存地址查找符号信息
// pid: 目标进程的 ID
// ptr: 目标进程内的内存地址
// 返回: 符号信息指针和错误 (如果发生)
func FindSymbolFromPidPtr(pid int, ptr uintptr) (*SymbolInfo, error) {
	if pid <= 0 {
		return nil, errors.New("invalid PID provided (must be > 0)")
	}
	if ptr == 0 {
		// Technically a valid address, but unlikely to hold a meaningful user symbol.
		// Can be adjusted if resolving symbols at address 0 is required.
		return nil, errors.New("invalid pointer (0x0) provided, usually not a user symbol location")
	}

	if time.Since(lastClearCacheTime) > 10*time.Minute {
		symbolCache = sync.Map{}
		lastClearCacheTime = time.Now()
	}

	// Check cache first
	if cached, ok := symbolCache.Load(fmt.Sprintf("%d_%x", pid, ptr)); ok {
		return cached.(*SymbolInfo), nil
	}

	addr2line, err := findAddr2line()
	if err != nil {
		return nil, err // Return error if addr2line is not found
	}

	mapsPath := fmt.Sprintf("/proc/%d/maps", pid)
	mapsFile, err := os.Open(mapsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("process with PID %d not found or /proc filesystem not mounted", pid)
		}
		if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied to read %s. Try running as root or the process owner", mapsPath)
		}
		return nil, fmt.Errorf("failed to open maps file %s: %w", mapsPath, err)
	}
	defer mapsFile.Close()

	var targetEntry *mapEntry = nil
	// 存储每个映射文件实例（路径+inode）对应的最低加载地址（文件偏移为0的段的起始地址）
	baseAddresses := make(map[string]uintptr)
	// 存储路径对应的第一个遇到的inode，用于备选查找
	filePathToInode := make(map[string]uint64)

	scanner := bufio.NewScanner(mapsFile)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseMapLine(line)
		if err != nil {
			// Log or ignore parsing errors for individual lines? For now, ignore.
			// fmt.Fprintf(os.Stderr, "Warning: skipping unparsable map line in PID %d: %v\n", pid, err)
			continue
		}

		// 检查 ptr 是否在此映射范围内
		if targetEntry == nil && ptr >= entry.StartAddr && ptr < entry.EndAddr {
			targetEntry = entry
			// Continue scanning to collect all base addresses
		}

		// 记录文件的基地址（第一个映射段，文件偏移为0）
		// 使用 Inode 来区分不同的文件，即使路径相同（例如被删除后重新创建的文件）
		if entry.Path != "" && entry.Inode != 0 && entry.FileOffset == 0 {
			uniqueFileID := fmt.Sprintf("%s:%d", entry.Path, entry.Inode)

			// Record the first inode seen for a given path
			if _, exists := filePathToInode[entry.Path]; !exists {
				// Check if the file path actually points to a regular file,
				// to avoid recording base addresses for special fs entries if needed.
				// This check adds overhead, maybe optional? For now, keep it simple.
				filePathToInode[entry.Path] = entry.Inode
			}

			// Store the base address, taking the minimum if multiple offset 0 segments exist (unlikely but possible)
			if currentBase, exists := baseAddresses[uniqueFileID]; !exists || entry.StartAddr < currentBase {
				baseAddresses[uniqueFileID] = entry.StartAddr
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading maps file %s: %w", mapsPath, err)
	}

	if targetEntry == nil {
		return nil, fmt.Errorf("pointer 0x%x not found in any mapped region for PID %d", ptr, pid)
	}

	// --- Found the target mapping entry, now determine base address and offset ---

	// Handle anonymous mappings or special mappings (e.g., [heap], [stack], [vdso])
	if targetEntry.Path == "" || strings.HasPrefix(targetEntry.Path, "[") {
		// For non-file mappings, offset is relative to the segment start. Base address is segment start.
		return &SymbolInfo{
			SymbolName:  fmt.Sprintf("in %s", targetEntry.Path), // More concise name
			FilePath:    targetEntry.Path,
			Offset:      ptr - targetEntry.StartAddr, // Offset within this specific segment
			BaseAddress: targetEntry.StartAddr,       // Base address is the start of this segment
			SourceFile:  "N/A",
			SourceLine:  0,
		}, nil
	}

	// --- It's a file-backed mapping, find its base address and use addr2line ---

	// Clean up path (remove "(deleted)")
	targetPath := targetEntry.Path
	isDeleted := false
	if strings.HasSuffix(targetPath, " (deleted)") {
		targetPath = strings.TrimSuffix(targetPath, " (deleted)")
		isDeleted = true
	}

	// Check if path is absolute (it usually is in /proc/pid/maps)
	if !filepath.IsAbs(targetPath) {
		// This is unlikely but good to handle. We might need PWD of the process for accuracy.
		fmt.Fprintf(os.Stderr, "Warning: Path '%s' from maps for PID %d is relative. addr2line might fail if not run from the correct directory.\n", targetPath, pid)
	}

	// Find the base address for this specific file instance (path + inode)
	uniqueFileID := fmt.Sprintf("%s:%d", targetPath, targetEntry.Inode) // Use cleaned path
	baseAddr, found := baseAddresses[uniqueFileID]

	if !found {
		// Fallback: Try finding the base address using the path and the *first* inode we saw for it.
		if firstInode, inodeFound := filePathToInode[targetPath]; inodeFound {
			fallbackUniqueID := fmt.Sprintf("%s:%d", targetPath, firstInode)
			baseAddr, found = baseAddresses[fallbackUniqueID]
			if found {
				fmt.Fprintf(os.Stderr, "Warning: Using base address associated with first encountered inode (%d) for path '%s' as specific inode (%d) wasn't found in offset 0 maps.\n", firstInode, targetPath, targetEntry.Inode)
			}
		}
	}

	if !found {
		// Even fallback failed. Cannot reliably determine base address.
		// Option 1: Error out.
		return nil, fmt.Errorf("could not determine base address (offset 0 mapping) for file '%s' (inode %d) used by pointer 0x%x in PID %d", targetPath, targetEntry.Inode, ptr, pid)
		// Option 2: Guess base address (less reliable)
		// baseAddr = targetEntry.StartAddr - targetEntry.FileOffset
		// fmt.Fprintf(os.Stderr, "Warning: Could not find offset 0 mapping for '%s', guessing base address 0x%x\n", targetPath, baseAddr)
	}

	// Calculate ptr offset relative to the file's determined load base address
	offset := ptr - baseAddr

	// Prepare result struct even before calling addr2line
	result := &SymbolInfo{
		SymbolName:  fmt.Sprintf("symbol at offset 0x%x", offset), // Default name
		FilePath:    targetPath,                                   // Use potentially cleaned path
		Offset:      offset,
		BaseAddress: baseAddr,
		SourceFile:  "??", // Default value
		SourceLine:  0,    // Default value
	}
	if isDeleted {
		result.FilePath += " (deleted)" // Add back for clarity in output
	}

	// Execute addr2line
	// -e: Specify executable/library
	// -f: Output function names
	// -C: Demangle C++ symbols
	// -i: Resolve inlined functions (provides more detail)
	// -p: Pretty-print address (makes parsing easier)
	cmd := exec.Command(addr2line, "-e", targetPath, "-fCi", "-p", fmt.Sprintf("0x%x", offset))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Env = append(os.Environ(), "LC_ALL=C")

	err = cmd.Run()
	if err != nil {
		// addr2line failed (e.g., file not found, no debug info, invalid offset)
		errMsg := fmt.Sprintf("addr2line execution failed for path '%s' offset 0x%x: %v", targetPath, offset, err)
		stderrMsg := strings.TrimSpace(stderr.String())
		if stderrMsg != "" {
			errMsg += fmt.Sprintf(" | addr2line stderr: %s", stderrMsg)
		}
		// Return the partially filled result struct along with the error
		return result, errors.New(errMsg)
	}

	// Parse addr2line output (last line is usually most specific with -i)
	outputLines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(outputLines) == 0 {
		// Empty output from addr2line? Unexpected.
		return result, errors.New("addr2line produced empty output")
	}

	// fmt.Fprintf(os.Stderr, "Warning: addr2line output: %s\n", stdout.String())

	lastLine := outputLines[len(outputLines)-1]

	// Expected format with -p: "0xOFFSET: FUNCTION at FILE:LINE"
	addrPrefix := fmt.Sprintf("0x%x: ", offset) // Note the space after colon
	if strings.HasPrefix(lastLine, addrPrefix) {
		relevantPart := strings.TrimPrefix(lastLine, addrPrefix)
		parts := strings.SplitN(relevantPart, " at ", 2)
		symbolName := strings.TrimSpace(parts[0])
		if symbolName != "??" && symbolName != "" { // Update symbol name if found
			result.SymbolName = symbolName
		} else {
			result.SymbolName = fmt.Sprintf("Symbol at offset 0x%x (no name)", offset) // More specific default
		}

		if len(parts) == 2 {
			fileLinePart := strings.TrimSpace(parts[1])
			// Find the last colon for line number, robust against colons in filenames
			if locIdx := strings.LastIndex(fileLinePart, ":"); locIdx != -1 {
				filePart := fileLinePart[:locIdx]
				linePart := fileLinePart[locIdx+1:]

				// Clean up line number part from extras like (discriminator N) or (inlined by...)
				if parenIdx := strings.Index(linePart, " ("); parenIdx != -1 {
					linePart = linePart[:parenIdx]
				}
				if discIdx := strings.Index(linePart, " discriminator "); discIdx != -1 {
					linePart = linePart[:discIdx]
				}
				linePart = strings.TrimSpace(linePart) // Trim potential spaces

				if lineInt, err := strconv.Atoi(linePart); err == nil {
					result.SourceFile = filePart
					result.SourceLine = lineInt
				} else {
					// Could not parse line number part, treat whole thing as file path
					result.SourceFile = fileLinePart
					result.SourceLine = 0 // Reset line number
				}
			} else {
				// No colon found, treat whole part as file path
				result.SourceFile = fileLinePart
				result.SourceLine = 0
			}
		}
	} else {
		// Fallback or unexpected format
		// vectorAdd(float const*, float const*, float*, int) at ??:?

		result.SymbolName = strings.Split(lastLine, " at ")[0]
	}

	symbolCache.Store(fmt.Sprintf("%d_%x", pid, ptr), result)
	return result, nil // Success
}
