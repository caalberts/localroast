package strings

import (
	"net/http"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/stretchr/testify/assert"
)

func TestParseSchemaFromString(t *testing.T) {
	var schemas []localroast.Schema
	var err error

	p := &Parser{}
	input := []string{
		"GET / 200",
		"POST /users 201",
		"PUT /user/1 403",
	}
	schemas, err = p.Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, len(input), len(schemas))

	assert.Equal(t, http.MethodGet, schemas[0].Method)
	assert.Equal(t, http.MethodPost, schemas[1].Method)
	assert.Equal(t, http.MethodPut, schemas[2].Method)

	assert.Equal(t, "/", schemas[0].Path)
	assert.Equal(t, "/users", schemas[1].Path)
	assert.Equal(t, "/user/1", schemas[2].Path)

	assert.Equal(t, http.StatusOK, schemas[0].Status)
	assert.Equal(t, http.StatusCreated, schemas[1].Status)
	assert.Equal(t, http.StatusForbidden, schemas[2].Status)

	p = &Parser{}
	input = []string{
		"GET / 200",
		"POST 201",
	}
	_, err = p.Parse(input)
	assert.NotNil(t, err)
}

var schemaTests = []struct {
	in       string
	expected localroast.Schema
}{
	{
		"GET / 200",
		localroast.Schema{
			Method: http.MethodGet,
			Path:   "/",
			Status: http.StatusOK,
		},
	},
	{
		"POST / 201",
		localroast.Schema{
			Method: http.MethodPost,
			Path:   "/",
			Status: http.StatusCreated,
		},
	},
	{
		"PUT /user/1 403",
		localroast.Schema{
			Method: http.MethodPut,
			Path:   "/user/1",
			Status: http.StatusForbidden,
		},
	},
}

func TestConvertStringToSchema(t *testing.T) {
	for _, test := range schemaTests {
		schema, err := toSchema(test.in)
		assert.Nil(t, err)
		assert.Equal(t, test.expected.Method, schema.Method)
		assert.Equal(t, test.expected.Path, schema.Path)
		assert.Equal(t, test.expected.Status, schema.Status)
	}
}

var validMatchTests = []struct {
	in    string
	valid bool
}{
	{"GET / 200", true},
	{"POST / 201", true},
	{"PUT /user/1 400", true},
	{"PATCH /user/2 403", true},
	{"DELETE /account/3 405", true},
	{"GET /index 200", true},
	{"GET /", false},
	{"GET / abc", false},
	{"GET abc", false},
	{"SEND / 200", false},
}

func TestValidStringMatch(t *testing.T) {
	for _, test := range validMatchTests {
		_, err := validMatch(test.in)
		if test.valid {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
