package json

import (
	"net/http"
	"strings"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/stretchr/testify/assert"
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
	p := Parser{}
	output := make(chan []localroast.Schema)

	go func() {
		err := p.Parse(strings.NewReader(validJSON), output)
		assert.Nil(t, err)
	}()

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
	p := Parser{}
	err := p.Parse(strings.NewReader(missingKeys), make(chan []localroast.Schema))
	assert.NotNil(t, err)
	assert.Equal(t, "Missing required fields: method, path, status", err.Error())
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
	p := Parser{}
	err := p.Parse(strings.NewReader(invalidJSON), make(chan []localroast.Schema))
	assert.NotNil(t, err)
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
	p := Parser{}
	err := p.Parse(strings.NewReader(jsonObject), make(chan []localroast.Schema))
	assert.NotNil(t, err)
}
