package collection

import (
	"reflect"
	"testing"
)

func TestCartesianProduct(t *testing.T) {
	type args struct {
		sets [][]interface{}
	}
	tests := []struct {
		name string
		args args
		want [][]interface{}
	}{
		{
			name: "test1",
			args: args{sets: [][]interface{}{{1, 2}, {3, 4}}},
			want: [][]interface{}{{1, 3}, {1, 4}, {2, 3}, {2, 4}},
		},
		{
			name: "test1",
			args: args{sets: [][]interface{}{{1, 1}, {3, 4}}},
			want: [][]interface{}{{1, 3}, {1, 4}, {1, 3}, {1, 4}},
		},
		{
			name: "test2",
			args: args{sets: [][]interface{}{{1, 2}}},
			want: [][]interface{}{{1}, {2}},
		},
		{
			name: "test3",
			args: args{sets: [][]interface{}{{}}},
			want: nil,
		},
		{
			name: "test4",
			args: args{sets: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CartesianProduct(tt.args.sets...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartesianProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
