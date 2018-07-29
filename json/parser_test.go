package json

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	Success bool   `json:"success"`
	ID      int    `json:"id"`
	IDs     []int  `json:"ids"`
	Message string `json:"message"`
}

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
	schemas, err := p.Parse(strings.NewReader(validJSON))

	assert.Nil(t, err)
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
	_, err := p.Parse(strings.NewReader(missingKeys))
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
	_, err := p.Parse(strings.NewReader(invalidJSON))
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
	_, err := p.Parse(strings.NewReader(jsonObject))
	assert.NotNil(t, err)
}
