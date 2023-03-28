package main

import (
	"reflect"
	"testing"
)

func Test_cutUtil(t *testing.T) {
	type args struct {
		lines   [][]string
		myFlags flags
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "test1",
			args: args{lines: [][]string{{"245:789", "4567", "M:4540"}},
				myFlags: flags{fields: 1},
			},
			want: [][]string{{"4567"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutUtil(tt.args.lines, tt.args.myFlags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cutUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}
