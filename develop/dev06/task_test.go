package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCutString(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		input  string
		want   string
	}{
		{
			name:   "Without flags",
			params: Params{},
			input:  "one two three",
			want:   "one two three",
		},
		{
			name:   "With delimiter flag and without delimiters in string",
			params: Params{Fields: []int{1}, Delimiter: ":"},
			input:  "one two three",
			want:   "one two three",
		},
		{
			name:   "With d, s flags and without delimiters in string",
			params: Params{Fields: []int{1}, Delimiter: ":", Separated: true},
			input:  "one two three",
			want:   "",
		},
		{
			name:   "With f and d flags",
			params: Params{Fields: []int{1, 4}, Delimiter: ":"},
			input:  "one:two:three:four:five",
			want:   "two:five",
		},
		{
			name:   "With wrong f flag",
			params: Params{Fields: []int{-1, 6}, Delimiter: ":"},
			input:  "one:two:three:four:five",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CutString(tt.input, tt.params)
			assert.Equal(t, tt.want, res)
		})
	}
}
