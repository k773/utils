package exec

import (
	"os/exec"
	"syscall"
)

func interruptProcessPID(pid int) (e error) {
	// Sending CTRL_BREAK_EVENT on windows (src: https://github.com/golang/go/blob/master/src/os/signal/signal_windows_test.go)
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return
	}
	return
}

func setupAttributesForInterruption(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP
}
