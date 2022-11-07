package tracer

func init() {
	registerTypeHandler(argTypeIntOrErrorCode, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if arg.Int() < 0 {
			arg.t = ArgTypeErrorCode
		} else {
			arg.t = ArgTypeInt
		}
		return nil
	})
}
