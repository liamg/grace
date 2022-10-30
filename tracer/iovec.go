package tracer

type iovec struct {
	Base uintptr /* Starting address */
	Len  uint    /* Number of bytes to transfer */
}

func convertIovecs(vecs []iovec) []Arg {
	var output []Arg
	for _, fd := range vecs {
		output = append(output, convertIovec(fd))
	}
	return output
}

func convertIovec(vec iovec) Arg {
	return Arg{
		t: ArgTypeObject,
		obj: &Object{
			Name: "iovec",
			Properties: []Arg{
				{
					name: "base",
					t:    ArgTypeAddress,
					raw:  vec.Base,
				},
				{
					name: "len",
					t:    ArgTypeUnsignedInt,
					raw:  uintptr(vec.Len),
				},
			},
		},
		known: true,
	}
}
