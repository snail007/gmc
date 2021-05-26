package gbytes

import (
	"reflect"
	"testing"
)

func TestByteSize_Bytes(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want uint64
	}{
		{"B", 1 * KB, uint64(1024 * B)},
		{"MB", 1 * MB, uint64(1024 * KB)},
		{"GB", 1 * GB, uint64(1024 * MB)},
		{"TB", 1 * TB, uint64(1024 * GB)},
		{"PB", 1 * PB, uint64(1024 * TB)},
		{"EB", 1 * EB, uint64(1024 * PB)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Bytes(); got != tt.want {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_EBytes(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want float64
	}{
		{"B", 1024 * 1024 * 1024 * 1024 * 1024 * KB, 1},
		{"MB", 1024 * 1024 * 1024 * 1024 * MB, 1},
		{"GB", 1024 * 1024 * 1024 * GB, 1},
		{"TB", 1024 * 1024 * TB, 1},
		{"PB", 1024 * PB, 1},
		{"EB", 1 * EB, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.EBytes(); got != tt.want {
				t.Errorf("EBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_GBytes(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want float64
	}{
		{"B", 1024 * 1024 * KB, 1},
		{"MB", 1024 * MB, 1},
		{"GB", GB, 1},
		{"TB", TB, float64(1024)},
		{"PB", PB, float64(1024 * 1024)},
		{"EB", EB, float64(1024 * 1024 * 1024)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GBytes(); got != tt.want {
				t.Errorf("GBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_HumanReadable(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want string
	}{
		{"B", 0, "0 B"},
		{"B", 2000 * B, "1.95 KB"},
		{"B", 2048 * B, "2 KB"},
		{"B", 2100 * B, "2.05 KB"},
		{"KB", 2000 * KB, "1.95 MB"},
		{"KB", 2048 * KB, "2 MB"},
		{"KB", 2100 * KB, "2.05 MB"},
		{"MB", 2000 * MB, "1.95 GB"},
		{"MB", 2048 * MB, "2 GB"},
		{"MB", 2100 * MB, "2.05 GB"},
		{"GB", 2000 * GB, "1.95 TB"},
		{"GB", 2048 * GB, "2 TB"},
		{"GB", 2100 * GB, "2.05 TB"},
		{"TB", 2000 * TB, "1.95 PB"},
		{"TB", 2048 * TB, "2 PB"},
		{"TB", 2100 * TB, "2.05 PB"},
		{"PB", 2000 * PB, "1.95 EB"},
		{"PB", 2048 * PB, "2 EB"},
		{"PB", 2100 * PB, "2.05 EB"},
		{"EB", 2 * EB, "2 EB"},
		{"EB", 10 * EB, "10 EB"},
		{"EB", 15 * EB, "15 EB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.HR(); got != tt.want {
				t.Errorf("HumanReadable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_KBytes(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want float64
	}{
		{"B", B, 0},
		{"B", 200 * B, 0.19},
		{"KB", KB, 1},
		{"MB", MB, 1024},
		{"GB", GB, 1024 * 1024},
		{"TB", TB, 1024 * 1024 * 1024},
		{"PB", PB, 1024 * 1024 * 1024 * 1024},
		{"EB", EB, 1024 * 1024 * 1024 * 1024 * 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.KBytes(); got != tt.want {
				t.Errorf("KBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_MBytes(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want float64
	}{
		{"B", B, 0},
		{"KB", KB, 0},
		{"KB", 200 * KB, 0.19},
		{"MB", MB, 1},
		{"GB", GB, 1024},
		{"TB", ByteSize(1.5 * float64(TB)), 1.5 * 1024 * 1024},
		{"PB", PB, 1024 * 1024 * 1024},
		{"EB", EB, 1024 * 1024 * 1024 * 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.MBytes(); got != tt.want {
				t.Errorf("MBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_PBytes(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want float64
	}{
		{"B", B, 0},
		{"KB", KB, 0},
		{"MB", MB, 0},
		{"GB", GB, 0},
		{"TB", TB, 0},
		{"PB", PB, 1},
		{"EB", EB, 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.PBytes(); got != tt.want {
				t.Errorf("PBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_String(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want string
	}{
		{"B", B, "1B"},
		{"KB", KB, "1KB"},
		{"MB", MB, "1MB"},
		{"GB", GB, "1GB"},
		{"TB", TB, "1TB"},
		{"PB", PB, "1PB"},
		{"EB", EB, "1EB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_TBytes(t *testing.T) {
	tests := []struct {
		name string
		b    ByteSize
		want float64
	}{
		{"B", B, 0},
		{"KB", KB, 0},
		{"MB", MB, 0},
		{"GB", GB, 0},
		{"TB", TB, 1},
		{"PB", PB, 1024},
		{"EB", EB, 1024 * 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.TBytes(); got != tt.want {
				t.Errorf("TBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteSize_Parse(t *testing.T) {
	type args struct {
		t []byte
	}
	tests := []struct {
		name    string
		b       ByteSize
		args    args
		wantErr bool
	}{
		{"B", 0, args{t: []byte("100000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000")}, true},
		{"B", 0, args{t: []byte("100000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000B")}, true},
		{"B", 0, args{t: []byte("1x")}, true},
		{"B", 0, args{t: []byte("1Kb")}, true},
		{"B", 1, args{t: []byte("1.1B")}, false},
		{"B", B, args{t: []byte("1B")}, false},
		{"KB", KB, args{[]byte("1KB")}, false},
		{"MB", MB, args{[]byte("1MB")}, false},
		{"GB", GB, args{[]byte("1GB")}, false},
		{"TB", TB, args{[]byte("1TB")}, false},
		{"PB", PB, args{[]byte("1PB")}, false},
		{"EB", EB, args{[]byte("1EB")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := tt.b
			if tt.b.MustParse(string(tt.args.t)); !reflect.DeepEqual(tt.b, old) {
				t.Errorf("MustParse() = %v, want %v", tt.b, old)
			}
			if err := tt.b.Parse(string(tt.args.t)); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.b != old {
				t.Errorf("Parse() want = %v, actual %v", tt.b, old)
			}
		})
	}
}

func TestSizeStr(t *testing.T) {
	type args struct {
		bytes uint64
	}
	tests := []struct {
		name    string
		args    args
		wantS   string
		wantErr bool
	}{
		{"B", args{1}, "1 B", false},
		{"KB", args{uint64(1.5 * float64(KB))}, "1.5 KB", false},
		{"KB", args{uint64(2 * float64(KB))}, "2 KB", false},
		{"MB", args{uint64(1.5 * float64(MB))}, "1.5 MB", false},
		{"MB", args{uint64(2 * float64(MB))}, "2 MB", false},
		{"GB", args{uint64(1.5 * float64(GB))}, "1.5 GB", false},
		{"GB", args{uint64(2 * float64(GB))}, "2 GB", false},
		{"TB", args{uint64(1.5 * float64(TB))}, "1.5 TB", false},
		{"TB", args{uint64(2 * float64(TB))}, "2 TB", false},
		{"PB", args{uint64(1.5 * float64(PB))}, "1.5 PB", false},
		{"PB", args{uint64(2 * float64(PB))}, "2 PB", false},
		{"EB", args{uint64(1.5 * float64(EB))}, "1.5 EB", false},
		{"EB", args{uint64(2 * float64(EB))}, "2 EB", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := SizeStr(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("SizeStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("SizeStr() gotS = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		wantBytes uint64
		wantErr   bool
	}{
		{"Kb", args{""}, 0, true},
		{"B", args{""}, 0, true},
		{"B", args{"1.x"}, 0, true},
		{"B", args{"1 B"}, 1, false},
		{"KB", args{"1.5 KB"}, uint64(1.5 * float64(KB)), false},
		{"KB", args{"2 KB"}, uint64(2 * float64(KB)), false},
		{"MB", args{"1.5 MB"}, uint64(1.5 * float64(MB)), false},
		{"MB", args{"2 MB"}, uint64(2 * float64(MB)), false},
		{"GB", args{"1.5 GB"}, uint64(1.5 * float64(GB)), false},
		{"GB", args{"2 GB"}, uint64(2 * float64(GB)), false},
		{"TB", args{"1.5 TB"}, uint64(1.5 * float64(TB)), false},
		{"TB", args{"2 TB"}, uint64(2 * float64(TB)), false},
		{"PB", args{"1.5 PB"}, uint64(1.5 * float64(PB)), false},
		{"PB", args{"2 PB"}, uint64(2 * float64(PB)), false},
		{"EB", args{"1.5 EB"}, uint64(1.5 * float64(EB)), false},
		{"EB", args{"2 EB"}, uint64(2 * float64(EB)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := ParseSize(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBytes != tt.wantBytes {
				t.Errorf("Parse() gotBytes = %v, want %v", gotBytes, tt.wantBytes)
			}
		})
	}
}

func TestByteSize_MustParse(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		b    ByteSize
		args args
		want *ByteSize
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.MustParse(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustParse() = %v, want %v", got, tt.want)
			}
		})
	}
}
