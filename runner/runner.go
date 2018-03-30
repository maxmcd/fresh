package runner

import (
	"io"
	"os/exec"
	"syscall"
)

func run(args []string) bool {
	runnerLog("Running...")

	cmd := exec.Command(buildPath(), args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Terminating PID %d", pid)
		// cmd.Process.Kill()
		// cmd.Process.Signal(syscall.SIGINT)
		cmd.Process.Signal(syscall.SIGTERM)
	}()

	return true
}
