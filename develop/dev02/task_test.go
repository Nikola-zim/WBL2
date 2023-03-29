package main

import "testing"

func TestUnpackString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1_test",
			args: args{
				s: "a4bc2d5e",
			},
			want: "aaaabccddddde",
		},
		{
			name: "2_test",
			args: args{
				s: "abcd",
			},
			want: "abcd",
		},
		{
			name: "3_test",
			args: args{
				s: "32",
			},
			want: "",
		},
		{
			name: "4_test",
			args: args{
				s: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnpackString(tt.args.s); got != tt.want {
				t.Errorf("UnpackString() = %s, want %s", got, tt.want)
			}
		})
	}
}
