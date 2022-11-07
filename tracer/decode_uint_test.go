package tracer

import "testing"

func Test_DecodeUint(t *testing.T) {

	tests := []struct {
		name string
		data []byte
		want uint64
	}{
		{
			name: "zero length",
			data: []byte{},
			want: 0,
		},
		{
			name: "filled byte",
			data: []byte{
				0xff,
			},
			want: 0xff,
		},
		{
			name: "byte order",
			data: []byte{
				0x11,
				0xff,
				0x00,
			},
			want: 0xff11,
		},
		{
			name: "byte max value",
			data: []byte{
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
			},
			want: 0xffffffffffffffff,
		},
		{
			name: "overflow ignored",
			data: []byte{
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
				0xff,
				0x1,
			},
			want: 0xffffffffffffffff,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := decodeUint(test.data); got != test.want {
				t.Errorf("decodeUint() = %v, want %v", got, test.want)
			}
		})
	}

}
