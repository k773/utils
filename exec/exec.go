package exec

import (
	"bytes"
	"context"
	"github.com/k773/utils"
	"io"
	"os/exec"
	"sync"
)

// CombinedOutputWithContext interrupts child by sending platform-specific signal to it (see InterruptProcess) on context expiration.
// All required preparations over cmd are performed by this function.
func CombinedOutputWithContext(ctx context.Context, cmd *exec.Cmd) (output []byte, e error) {
	var w sync.WaitGroup
	var goRun = func(f func()) { utils.ExecuteWitWGAsync(&w, f) }
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	SetupAttributesForInterruption(cmd)

	// joining out-pipes into a single out
	outPipe, e1 := cmd.StdoutPipe()
	errPipe, e2 := cmd.StderrPipe()
	var buf = bytes.NewBuffer(nil)
	if e = utils.JoinErrors(e1, e2); e == nil {
		goRun(func() { _, _ = io.Copy(buf, io.MultiReader(outPipe, errPipe)) })

		if e = cmd.Start(); e == nil {
			// Planning interruption
			var interruptionError error
			goRun(func() {
				select {
				case <-ctx.Done():
					interruptionError = ctx.Err()
					_ = InterruptProcess(cmd)
				}
			})
			e = cmd.Wait()
			// If the child was killed due to the context expiration, preserve the error
			if interruptionError != nil {
				e = interruptionError
			}
		}
	}

	cancel()
	w.Wait()
	output = buf.Bytes()

	return
}

// InterruptProcess is a convenience wrapper for InterruptProcessPid
func InterruptProcess(cmd *exec.Cmd) error {
	return InterruptProcessPID(cmd.Process.Pid)
}

// InterruptProcessPID sends different signals depending on the os:
// on Linux: sends SIGINT
// on Windows: sends CTRL_BREAK_EVENT, child process needs to be started with syscall.CREATE_NEW_PROCESS_GROUP, otherwise parent will receive this event too
func InterruptProcessPID(pid int) (e error) {
	return interruptProcessPID(pid)
}

// SetupAttributesForInterruption should be called before cmd.Start() is called; sets up correct cmd's interruption behaviour:
// on Windows: sets cmd.SysProcAttr.CreationFlags=syscall.CREATE_NEW_PROCESS_GROUP (so the parent won't receive CTRL_BREAK_EVENT when it is sent to the child)
func SetupAttributesForInterruption(cmd *exec.Cmd) {
	setupAttributesForInterruption(cmd)
}
