package json

import (
	"net/http"
	"strings"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/stretchr/testify/assert"
	"io"
)

var validJSON = `
[
    {
        "method": "GET",
        "path": "/",
        "status": 200,
        "response": {
            "success": true
        }
    },
    {
        "method": "GET",
        "path": "/users",
        "status": 200,
        "response": {
            "success": true,
            "ids": [1, 2, 3]
        }
	}
]`

func TestParse(t *testing.T) {
	input := make(chan io.Reader)
	output := make(chan []localroast.Schema)

	go func() {
		p := Parser{}
		err := p.Parse(input, output)
		assert.Nil(t, err)
	}()

	input <- strings.NewReader(validJSON)
	schemas := <-output

	assert.Equal(t, 2, len(schemas))

	assert.Equal(t, http.MethodGet, schemas[0].Method)
	assert.Equal(t, http.MethodGet, schemas[1].Method)

	assert.Equal(t, "/", schemas[0].Path)
	assert.Equal(t, "/users", schemas[1].Path)

	assert.Equal(t, http.StatusOK, schemas[0].Status)
	assert.Equal(t, http.StatusOK, schemas[1].Status)

	assert.JSONEq(t, `{ "success": true }`, string(schemas[0].Response))
	assert.JSONEq(t, `{ "success": true, "ids": [1, 2, 3] }`, string(schemas[1].Response))
}

var missingKeys = `
[
	{
		"response": {
			"success": true
		}
	}
]
`

func TestParseSchemaWithMissingKeys(t *testing.T) {
	input := make(chan io.Reader)
	go func() {
		p := Parser{}
		err := p.Parse(input, make(chan []localroast.Schema))
		assert.NotNil(t, err)
		assert.Equal(t, "missing required fields: method, path, status", err.Error())
	}()
	input <- strings.NewReader(missingKeys)
}

var invalidJSON = `
[
    {
        method: "GET",
        path: "/",
        status: 200,
        response: {
            "success": true
        }
    }
]
`

func TestParseSchemaFromInvalidJSON(t *testing.T) {
	input := make(chan io.Reader)
	go func() {
		p := Parser{}
		err := p.Parse(input, make(chan []localroast.Schema))
		assert.NotNil(t, err)
	}()
	input <- strings.NewReader(invalidJSON)
}

var jsonObject = `
{
	"method": "POST",
	"path": "/users",
	"status": 201,
	"response": {
		"success": true,
		"id": 4
	}
}
`

func TestParseSchemaFromJSONObject(t *testing.T) {
	input := make(chan io.Reader)
	go func() {
		p := Parser{}
		err := p.Parse(input, make(chan []localroast.Schema))
		assert.NotNil(t, err)
	}()
	input <- strings.NewReader(jsonObject)
}
