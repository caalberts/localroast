package schema

import (
	"net/http"
	"testing"

	"github.com/caalberts/localghost"
	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	var definition string
	var schema localghost.Schema

	definition = "GET / 200"
	schema, err := FromString(definition)
	assert.Nil(t, err)
	assert.Equal(t, "GET", schema.Method)
	assert.Equal(t, "/", schema.Path)
	assert.Equal(t, http.StatusOK, schema.StatusCode)
}
