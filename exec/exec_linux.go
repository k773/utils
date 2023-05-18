package exec

import (
	"os/exec"
	"syscall"
)

func interruptProcessPID(pid int) (e error) {
	// Sending SIGINT on linux
	return syscall.Kill(pid, syscall.SIGINT)
}

func setupAttributesForInterruption(cmd *exec.Cmd) {}
