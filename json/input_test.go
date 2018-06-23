package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	var bytes []byte
	var err error
	r := FileReader{}

	t.Run("valid file", func(t *testing.T) {
		bytes, err = r.Read([]string{"../examples/stubs.json"})
		assert.Nil(t, err)
		assert.NotNil(t, bytes)
	})

	t.Run("non json file", func(t *testing.T) {
		bytes, err = r.Read([]string{"stubs.txt"})
		assert.NotNil(t, err)
		assert.Equal(t, "Input must be a JSON file", err.Error())
	})

	t.Run("without argument", func(t *testing.T) {
		bytes, err = r.Read([]string{})
		assert.NotNil(t, err)
		assert.Equal(t, "A file is required", err.Error())
	})

	t.Run("too many arguments", func(t *testing.T) {
		bytes, err = r.Read([]string{"abc.json", "def.txt"})
		assert.NotNil(t, err)
		assert.Equal(t, "Too many arguments", err.Error())
	})

	t.Run("with incorrect file", func(t *testing.T) {
		bytes, err = r.Read([]string{"unknownfile"})
		assert.NotNil(t, err)
	})
}
