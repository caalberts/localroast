package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCLI(t *testing.T) {
	var res []string
	var err error
	r := &CLIReader{}

	t.Run("valid input", func(t *testing.T) {
		res, err = r.Read([]string{"1", "2"})
		assert.Nil(t, err)
		assert.Equal(t, []string{"1", "2"}, res)
	})

	t.Run("missing input", func(t *testing.T) {
		_, err = r.Read([]string{})
		assert.NotNil(t, err)
		assert.Equal(t, "Please define an endpoint in the format '<METHOD> <PATH> <STATUS_CODE>'. e.g 'GET / 200'", err.Error())
	})
}
