package tracer

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

type Tracer struct {
	handlers struct {
		syscallExit  func(*Syscall)
		syscallEnter func(*Syscall)
		signal       func(syscall.Signal)
		processExit  func(int)
	}
	pid        int
	cmd        *exec.Cmd
	isExit     bool
	lastCall   *Syscall
	lastSignal int
}

func New(pid int) *Tracer {
	return &Tracer{
		pid: pid,
	}
}

func FromCommand(suppressOutput bool, command string, args ...string) (*Tracer, error) {

	runtime.LockOSThread()

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	if !suppressOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Tracer{
		pid: cmd.Process.Pid,
		cmd: cmd,
	}, nil
}

func (t *Tracer) SetSyscallExitHandler(handler func(*Syscall)) {
	t.handlers.syscallExit = handler
}

func (t *Tracer) SetSyscallEnterHandler(handler func(*Syscall)) {
	t.handlers.syscallEnter = handler
}

func (t *Tracer) SetSignalHandler(handler func(syscall.Signal)) {
	t.handlers.signal = handler
}

func (t *Tracer) SetProcessExitHandler(handler func(int)) {
	t.handlers.processExit = handler
}

func (t *Tracer) Start() error {

	runtime.LockOSThread()

	if _, err := os.FindProcess(t.pid); err != nil {
		return fmt.Errorf("could not find process with pid %d: %w", t.pid, err)
	}

	if t.cmd == nil {
		if err := syscall.PtraceAttach(t.pid); err == syscall.EPERM {
			return fmt.Errorf("could not attach to process with pid %d: %w - check your permissions", t.pid, err)
		} else if err != nil {
			return err
		}
	}

	status := syscall.WaitStatus(0)
	if _, err := syscall.Wait4(t.pid, &status, 0, nil); err != nil {
		return err
	}

	if t.cmd == nil {
		defer func() {
			_ = syscall.PtraceDetach(t.pid)
			_, _ = syscall.Wait4(t.pid, &status, 0, nil)
		}()
	}

	// deliver SIGTRAP|0x80
	if err := syscall.PtraceSetOptions(t.pid, syscall.PTRACE_O_TRACESYSGOOD); err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGQUIT)
	go func() {
		for sig := range signalChan {
			interrupted := sig.(syscall.Signal)
			fmt.Printf("((SIGNAL: %s))", interrupted) // TODO
			_ = syscall.Kill(t.pid, syscall.SIGSTOP)
		}
	}()

	return t.loop()
}

var errExited = fmt.Errorf("process exited")

func (t *Tracer) loop() error {
	for {
		if err := t.waitForSyscall(); err != nil {
			if err == errExited {
				return nil
			}
			return err
		}
	}
}

func isStopSig(sig syscall.Signal) bool {
	return sig == syscall.SIGSTOP || sig == syscall.SIGTSTP || sig == syscall.SIGTTIN || sig == syscall.SIGTTOU
}

func (t *Tracer) waitForSyscall() error {

	// intercept syscall
	err := syscall.PtraceSyscall(t.pid, t.lastSignal)
	if err != nil {
		return fmt.Errorf("could not intercept syscall: %w", err)
	}

	// wait for a syscall
	status := syscall.WaitStatus(0)
	_, err = syscall.Wait4(t.pid, &status, 0, nil)
	if err != nil {
		return fmt.Errorf("wait failed: %w", err)
	}

	if status.TrapCause() != -1 {
		return nil
	}

	if status.Exited() {
		if t.handlers.processExit != nil {
			t.handlers.processExit(status.ExitStatus())
		}
		return errExited
	}

	if status.StopSignal() != syscall.SIGTRAP|0x80 {

		// NOTE: once https://groups.google.com/g/golang-codereviews/c/t2SwaIV-hFs is merged, we can use
		// syscall.PtraceGetSigInfo() to retrieve the siginfo_t struct and pass it to the signal handler instead
		if t.handlers.signal != nil {
			t.handlers.signal(status.StopSignal())
		}

		if isStopSig(status.StopSignal()) {
			// TODO: if we received an interrupt via notify above, return an error here to break the loop
			t.lastSignal = int(status.StopSignal())
		} else if t.lastSignal != 0 {
			if status.StopSignal() == syscall.SIGCONT {
				t.lastSignal = 0
				return nil
			}
			if err = syscall.PtraceSyscall(t.pid, t.lastSignal); err != nil && err != syscall.ESRCH {
				return err
			}
			return nil
		}
		return nil
	}

	// if interrupted, stop tracing
	if status.StopSignal().String() == "interrupt" {
		_ = syscall.PtraceSyscall(t.pid, int(status.StopSignal()))
		return fmt.Errorf("process interrupted")
	}

	// read registers
	regs := &syscall.PtraceRegs{}
	if err := syscall.PtraceGetRegs(t.pid, regs); err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	call := parseSyscall(regs)
	call.pid = t.pid

	if call.number == -1 {
		return fmt.Errorf("expecting syscall but received -1 - did we miss a signal?")
	}

	if t.isExit && t.lastCall != nil {
		if call.number == t.lastCall.number {
			call.args = t.lastCall.args
		} else {
			return fmt.Errorf("syscall exit mismatch: %d != %d - this is likely a bug in grace due to an unprocessed signal", call.number, t.lastCall.number)
		}
	}

	if err := call.populate(t.isExit); err != nil {
		return fmt.Errorf("populate failed: %w", err)
	}

	if t.isExit {
		if t.handlers.syscallExit != nil {
			t.handlers.syscallExit(call)
		}
	} else if t.handlers.syscallEnter != nil {
		t.handlers.syscallEnter(call)
	}
	t.lastCall = call
	t.isExit = !t.isExit
	return nil
}
