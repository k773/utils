package exec

import (
	"errors"
	"os/exec"
	"syscall"
)

func interruptProcessPID(pid int) (e error) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return
	}
	r, _, e := p.Call(syscall.CTRL_C_EVENT, uintptr(pid))
	if r == 0 {
		if e == nil {
			return errors.New("sending CTRL_C_EVENT: " + e.Error())
		} else {
			return errors.New("sending CTRL_C_EVENT: unknown error")
		}
	}
	return nil
}

// make sure that spawned process makes a call to SetConsoleCtrlHandler(NULL, FALSE)
func setupAttributesForInterruption(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP
}
