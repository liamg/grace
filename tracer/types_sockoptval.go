package tracer

func init() {
	registerTypeHandler(argTypeSockoptval, func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error {
		if raw > 0 {
			rawSockOptVal, err := readSize(pid, raw, next)
			if err != nil {
				return err
			}

			switch prev {
			// TODO: expect specific types for specific options?
			default:
				if next == 4 { // if it's exactly 4 bytes, it's probably an int
					arg.raw = uintptr(decodeInt(rawSockOptVal))
					arg.t = ArgTypeInt
				} else { // otherwise it's a big pile of misc data
					arg.data = rawSockOptVal
					arg.t = ArgTypeData
				}
			}
		} else {
			arg.t = ArgTypeAddress
			arg.annotation = "NULL"
			arg.replace = true
		}
		return nil
	})
}
