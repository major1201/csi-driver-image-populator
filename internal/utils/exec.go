package utils

import (
	"os/exec"
	"syscall"
	"time"
)

// SafeExec exec command safely with timeout
//
// If the command ended within timeout duration, it returns with (false, run error)
//
// If the command not been ended within timeout duration,
// it would be killed with SIGKILL recursively (both parent and children)
func SafeExec(cmd *exec.Cmd, timeout time.Duration, fn func(*exec.Cmd) error) (killed bool, err error) {
	if cmd == nil {
		return
	}

	// kill recursively (both parent and children)
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	} else {
		cmd.SysProcAttr.Setpgid = true
	}

	// do your stuff
	chDone := make(chan struct{}, 1)
	go func() {
		err = fn(cmd)
		chDone <- struct{}{}
	}()

	// wait done or timeout
	select {
	case <-time.After(timeout):
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			killed = true
		}
	case <-chDone:
	}
	return
}

// SafeExecWithCombinedOutput safely run the command with cmd.CombinedOutput() and returns the output
//
// It is the wrapper of SafeExec()
func SafeExecWithCombinedOutput(cmd *exec.Cmd, timeout time.Duration) (output []byte, killed bool, err error) {
	killed, err = SafeExec(cmd, timeout, func(cmd *exec.Cmd) error {
		var outErr error
		output, outErr = cmd.CombinedOutput()
		return outErr
	})
	return
}
