package schema

import (
	"net/http"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/stretchr/testify/assert"
)

func TestFromStrings(t *testing.T) {
	var schemas []localroast.Schema
	var err error

	definitions := []string{
		"GET / 200",
		"POST /user 201",
	}
	schemas, err = FromStrings(definitions)
	assert.Nil(t, err)
	assert.Equal(t, len(definitions), len(schemas))

	assert.Equal(t, "GET", schemas[0].Method)
	assert.Equal(t, "POST", schemas[1].Method)

	assert.Equal(t, "/", schemas[0].Path)
	assert.Equal(t, "/user", schemas[1].Path)

	assert.Equal(t, http.StatusOK, schemas[0].StatusCode)
	assert.Equal(t, http.StatusCreated, schemas[1].StatusCode)

	definitions = []string{
		"GET / 200",
		"POST 201",
	}
	_, err = FromStrings(definitions)
	assert.NotNil(t, err)
}

func TestFromString(t *testing.T) {
	var definition string
	var schema localroast.Schema
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
