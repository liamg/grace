package tracer

import "fmt"

type Syscall struct {
	pid     int
	number  int
	rawArgs [6]uintptr
	args    []Arg
	rawRet  uintptr
	ret     Arg
	unknown bool
}

type SyscallMetadata struct {
	Name        string
	Args        []ArgMetadata
	ReturnValue ReturnMetadata
}

func LookupSyscall(number int) (*SyscallMetadata, error) {
	meta, ok := sysMap[number]
	if !ok {
		return nil, fmt.Errorf("unknown syscall %d", number)
	}
	return &meta, nil
}

func (s *Syscall) Number() int {
	return s.number
}

func (s *Syscall) Name() string {
	meta, ok := sysMap[s.number]
	if !ok {
		return fmt.Sprintf("unknown_syscall_%d", s.number)
	}
	return meta.Name
}

func (s *Syscall) Args() []Arg {
	return s.args
}

func (s *Syscall) Return() Arg {
	return s.ret
}

func (s *Syscall) Unknown() bool {
	return s.unknown
}

func (s *Syscall) populate(exit bool) error {
	meta, ok := sysMap[s.number]
	if !ok {
		s.unknown = true
	}
	ret, err := processArgument(s.rawRet, 0, 0, ArgMetadata(meta.ReturnValue), s.pid, exit)
	if err != nil {
		return fmt.Errorf("failed to set return value of syscall %s (%d): %w", meta.Name, s.number, err)
	}
	s.ret = *ret
	for i, argMeta := range meta.Args {
		if exit && !argMeta.Destination && i < len(s.args) {
			continue
		}

		var next uintptr
		if i < len(meta.Args)-1 {
			next = s.rawArgs[i+1]
		}
		arg, err := processArgument(s.rawArgs[i], next, s.rawRet, argMeta, s.pid, exit)
		if err != nil {
			return fmt.Errorf("failed to set argument %d (%s) of syscall %s (%d): %w", i, argMeta.Name, meta.Name, s.number, err)
		}
		if i >= len(s.args) {
			s.args = append(s.args, *arg)
		} else {
			s.args[i] = *arg
		}
	}

	// strip off trailing optional args if they have no value
	var lastIndex int
	for i, arg := range s.args {
		meta := meta.Args[i]
		if !meta.Optional || arg.Raw() > 0 {
			lastIndex = i
		}
	}
	if lastIndex < len(s.args)-1 {
		s.args = s.args[:lastIndex+1]
	}

	return nil
}
