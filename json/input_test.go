package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	var err error
	v := Validator{}

	t.Run("valid file", func(t *testing.T) {
		err = v.Validate([]string{"../examples/stubs.json"})
		assert.Nil(t, err)
	})

	t.Run("non json file", func(t *testing.T) {
		err = v.Validate([]string{"stubs.txt"})
		assert.NotNil(t, err)
		assert.Equal(t, "Input must be a JSON file", err.Error())
	})

	t.Run("without argument", func(t *testing.T) {
		err = v.Validate([]string{})
		assert.NotNil(t, err)
		assert.Equal(t, "A file is required", err.Error())
	})

	t.Run("too many arguments", func(t *testing.T) {
		err = v.Validate([]string{"abc.json", "def.txt"})
		assert.NotNil(t, err)
		assert.Equal(t, "Too many arguments", err.Error())
	})

	t.Run("with incorrect file", func(t *testing.T) {
		err = v.Validate([]string{"unknownfile"})
		assert.NotNil(t, err)
	})
}
