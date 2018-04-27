package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	var bytes []byte
	var err error
	r := &FileReader{}

	t.Run("valid file", func(t *testing.T) {
		bytes, err = r.Read([]string{"../examples/stubs.json"})
		assert.Nil(t, err)
		assert.NotNil(t, bytes)
	})

	t.Run("without argument", func(t *testing.T) {
		bytes, err = r.Read([]string{})
		assert.NotNil(t, err)
		assert.Equal(t, "A file is required", err.Error())
	})

	t.Run("with incorrect file", func(t *testing.T) {
		bytes, err = r.Read([]string{"unknownfile"})
		assert.NotNil(t, err)
	})
}
