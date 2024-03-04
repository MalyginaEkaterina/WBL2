package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		input  string
		want   string
	}{
		{
			name:   "Without flags",
			params: Params{ColumnNum: -1},
			input:  "grt\nabc\ncbv\n",
			want:   "abc\ncbv\ngrt\n",
		},
		{
			name:   "With reverse",
			params: Params{ColumnNum: -1, Reverse: true},
			input:  "grt\nabc\ncbv\n",
			want:   "grt\ncbv\nabc\n",
		},
		{
			name:   "With numbers",
			params: Params{ColumnNum: -1, IsNumber: true},
			input:  "123\n45\n300\n",
			want:   "45\n123\n300\n",
		},
		{
			name:   "With numbers and reverse",
			params: Params{ColumnNum: -1, IsNumber: true, Reverse: true},
			input:  "123\n45\n300\n",
			want:   "300\n123\n45\n",
		},
		{
			name:   "Without flags and with spaces",
			params: Params{ColumnNum: -1},
			input:  "  hhh\nfff\nabc\n",
			want:   "  hhh\nabc\nfff\n",
		},
		{
			name:   "With ignoring spaces",
			params: Params{ColumnNum: -1, IgnoreSpaces: true},
			input:  "  hhh\nfff\nabc\n",
			want:   "abc\nfff\n  hhh\n",
		},
		{
			name:   "With column",
			params: Params{ColumnNum: 1},
			input:  "hhh kk\nfff aa\nabc yz\n",
			want:   "fff aa\nhhh kk\nabc yz\n",
		},
		{
			name:   "With column and numbers and reverse",
			params: Params{ColumnNum: 1, IsNumber: true, Reverse: true},
			input:  "hhh 123\nfff 45\nabc 90\n",
			want:   "hhh 123\nabc 90\nfff 45\n",
		},
		{
			name:   "With column and numbers and reverse and deleting duplicates",
			params: Params{ColumnNum: 1, IsNumber: true, Reverse: true, NoDuplicates: true},
			input:  "hhh 123\nfff 45\nabc 90\nfff 45\n",
			want:   "hhh 123\nabc 90\nfff 45\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := bytes.NewBufferString(tt.input)
			var buf bytes.Buffer
			err := Sort(tt.params, in, &buf)
			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
