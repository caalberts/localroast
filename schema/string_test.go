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
	var err error

	definition = "GET / 200"
	schema, err = FromString(definition)
	assert.Nil(t, err)
	assert.Equal(t, "GET", schema.Method)
	assert.Equal(t, "/", schema.Path)
	assert.Equal(t, http.StatusOK, schema.StatusCode)

	definition = "POST / 201"
	schema, err = FromString(definition)
	assert.Nil(t, err)
	assert.Equal(t, "POST", schema.Method)
	assert.Equal(t, "/", schema.Path)
	assert.Equal(t, http.StatusCreated, schema.StatusCode)
}

func TestValidMatch(t *testing.T) {
	var source string
	var err error

	source = "GET / 200"
	_, err = ValidMatch(source)
	assert.Nil(t, err)

	source = "POST / 201"
	_, err = ValidMatch(source)
	assert.Nil(t, err)

	source = "GET /index 200"
	_, err = ValidMatch(source)
	assert.Nil(t, err)

	source = "GET /"
	_, err = ValidMatch(source)
	assert.NotNil(t, err)

	source = "GET / abc"
	_, err = ValidMatch(source)
	assert.NotNil(t, err)

	source = "GET abc"
	_, err = ValidMatch(source)
	assert.NotNil(t, err)

	source = "SEND / 200"
	_, err = ValidMatch(source)
	assert.NotNil(t, err)
}
