package agentmanager

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"scope/internal/utils"
	"syscall"
)

func RunEBPF(ebpf_name string, args []string) int {
	bpfdir := utils.GetEnvOrDefault("BPF_DIR", "/home/delta/workspace/ebpf-golang/bpf")
	binpath := bpfdir + "/build/" + ebpf_name

	// 验证文件是否存在
	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		log.Printf("Error: eBPF program not found at %s", binpath)
		return -1
	}

	// 准备命令
	cmd := exec.Command(binpath, args...)

	// 启动进程
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting eBPF program: %v", err)
		return -1
	}

	// 获取进程ID并返回
	pid := cmd.Process.Pid

	// 分离进程，避免僵尸进程
	go func() {
		cmd.Process.Release()
	}()

	return pid
}

func StopProcess(pid int) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("invalid process ID")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}

	if process.Signal(syscall.Signal(0)) != nil {
		return false, fmt.Errorf("process %d is not running", pid)
	}

	process.Signal(syscall.SIGTERM)
	process.Signal(syscall.SIGINT)

	return true, nil
}
