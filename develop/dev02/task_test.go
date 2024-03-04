package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnpack(t *testing.T) {
	_, err := Unpack("4")
	require.Error(t, err)

	_, err = Unpack("a\\m")
	require.Error(t, err)

	_, err = Unpack("a\\")
	require.Error(t, err)

	_, err = Unpack("4a")
	require.Error(t, err)

	s1, err := Unpack("a4bc2d5e")
	require.NoError(t, err)
	assert.Equal(t, "aaaabccddddde", s1)

	s2, err := Unpack("abcd")
	require.NoError(t, err)
	assert.Equal(t, "abcd", s2)

	s3, err := Unpack("")
	require.NoError(t, err)
	assert.Equal(t, "", s3)

	s4, err := Unpack("qwe\\4\\5")
	require.NoError(t, err)
	assert.Equal(t, "qwe45", s4)

	s5, err := Unpack("qwe\\45")
	require.NoError(t, err)
	assert.Equal(t, "qwe44444", s5)

	s6, err := Unpack("qwe\\\\5")
	require.NoError(t, err)
	assert.Equal(t, "qwe\\\\\\\\\\", s6)

	s7, err := Unpack("qwe\\\\")
	require.NoError(t, err)
	assert.Equal(t, "qwe\\", s7)
}
