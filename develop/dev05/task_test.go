package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		name       string
		params     Params
		readerName string
		expr       string
		input      string
		want       string
	}{
		{
			name:   "Without flags",
			params: Params{},
			expr:   "text",
			input:  "la la Text la\nthis is text one\ntra TEXT ta ta\ntext two\n",
			want:   "this is text one\ntext two\n",
		},
		{
			name:   "With ignore case",
			params: Params{IgnoreCase: true},
			expr:   "text",
			input:  "la la Text la\nthis is text one\ntra TEXT ta ta\ntext two\n",
			want:   "la la Text la\nthis is text one\ntra TEXT ta ta\ntext two\n",
		},
		{
			name:   "With invert flag",
			params: Params{Invert: true},
			expr:   "text",
			input:  "la la Text la\nthis is text one\ntra TEXT ta ta\ntext two\n",
			want:   "la la Text la\ntra TEXT ta ta\n",
		},
		{
			name:   "With Fixed flag",
			params: Params{Fixed: true},
			expr:   "[",
			input:  "la la [Text] la\nthis is text one\ntra TEXT ta ta\ntext [two]\n",
			want:   "la la [Text] la\ntext [two]\n",
		},
		{
			name:   "With Fixed and ignore case flags",
			params: Params{Fixed: true, IgnoreCase: true},
			expr:   "[text",
			input:  "la la [Text] la\nthis is text one\ntra [TEXT ta ta\ntext [two]\n",
			want:   "la la [Text] la\ntra [TEXT ta ta\n",
		},
		{
			name:   "With Fixed and invert flags",
			params: Params{Fixed: true, Invert: true},
			expr:   "[",
			input:  "la la [Text] la\nthis is text one\ntra TEXT ta ta\ntext [two]\n",
			want:   "this is text one\ntra TEXT ta ta\n",
		},
		{
			name:   "With Fixed and invert and ignore case flags",
			params: Params{Fixed: true, Invert: true, IgnoreCase: true},
			expr:   "[text]",
			input:  "la la [Text] la\nthis is text one\ntra TEXT ta ta\ntext [two]\n",
			want:   "this is text one\ntra TEXT ta ta\ntext [two]\n",
		},
		{
			name:   "With Count flag",
			params: Params{Count: true},
			expr:   "text",
			input:  "la la Text la\nthis is text one\ntra TEXT ta ta\ntext two\n",
			want:   "2\n",
		},
		{
			name:       "With line number flag",
			params:     Params{LineNum: true},
			readerName: "buf",
			expr:       "text",
			input:      "la la Text la\nthis is text one\ntra TEXT ta ta\ntext two\n",
			want:       "buf:2:this is text one\nbuf:4:text two\n",
		},
		{
			name:   "With Before number",
			params: Params{Before: 2},
			expr:   "text",
			input:  "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntext two\n",
			want:   "this is tEXT one\ntra TEXT ta ta\ntext two\n",
		},
		{
			name:       "With line number and Before number",
			params:     Params{LineNum: true, Before: 2},
			readerName: "buf",
			expr:       "text",
			input:      "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntext two\n",
			want:       "buf:2:this is tEXT one\nbuf:3:tra TEXT ta ta\nbuf:4:text two\n",
		},
		{
			name:       "With line number, Before number and after number",
			params:     Params{LineNum: true, Before: 2, After: 1},
			readerName: "buf",
			expr:       "tEXT",
			input:      "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntEXT two\n",
			want:       "buf:1:la la Text la\nbuf:2:this is tEXT one\nbuf:3:tra TEXT ta ta\nbuf:4:tEXT two\n",
		},
		{
			name:       "With line number and after number",
			params:     Params{LineNum: true, After: 2},
			readerName: "buf",
			expr:       "Text",
			input:      "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntext two\n",
			want:       "buf:1:la la Text la\nbuf:2:this is tEXT one\nbuf:3:tra TEXT ta ta\n",
		},
		{
			name:   "With line number and Before number",
			params: Params{Before: 2},
			expr:   "text",
			input:  "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntext two\n",
			want:   "this is tEXT one\ntra TEXT ta ta\ntext two\n",
		},
		{
			name:   "With line number, Before number and after number",
			params: Params{Before: 2, After: 1},
			expr:   "tEXT",
			input:  "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntEXT two\n",
			want:   "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntEXT two\n",
		},
		{
			name:   "With after number",
			params: Params{After: 2},
			expr:   "Text",
			input:  "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntext two\n",
			want:   "la la Text la\nthis is tEXT one\ntra TEXT ta ta\n",
		},
		{
			name:   "Without matches",
			params: Params{Before: 2},
			expr:   "notfound",
			input:  "la la Text la\nthis is tEXT one\ntra TEXT ta ta\ntext two\n",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := bytes.NewBufferString(tt.input)
			var buf bytes.Buffer
			err := Grep(tt.expr, tt.params, tt.readerName, in, &buf)
			require.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
