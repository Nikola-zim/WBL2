package main

import (
	"reflect"
	"testing"
)

func Test_findAnagramSets(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "test1",
			args: args{words: []string{"listen", "silent", "elbow", "below", "state", "taste"}},
			want: map[string][]string{
				"below":  {"below", "elbow"},
				"listen": {"listen", "silent"},

				"state": {"state", "taste"},
			},
		},
		{
			name: "test2",
			args: args{words: []string{"tea", "bat"}},
			want: map[string][]string{},
		},
		{
			name: "test2",
			args: args{words: []string{"tea", "bat", "tab", "tab"}},
			want: map[string][]string{
				"bat": {"bat", "tab"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findAnagramSets(tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAnagramSets() = %v, want %v", got, tt.want)
			}
		})
	}
}
