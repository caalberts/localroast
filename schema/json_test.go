package schema

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validJSON, _ = ioutil.ReadFile("../examples/stubs.json")

func TestBytesToSchema(t *testing.T) {
	schemas, err := BytesToSchema(validJSON)

	assert.Nil(t, err)
	assert.Equal(t, 4, len(schemas))

	assert.Equal(t, http.MethodGet, schemas[0].Method)
	assert.Equal(t, http.MethodGet, schemas[1].Method)
	assert.Equal(t, http.MethodPost, schemas[2].Method)
	assert.Equal(t, http.MethodGet, schemas[3].Method)

	assert.Equal(t, "/", schemas[0].Path)
	assert.Equal(t, "/users", schemas[1].Path)
	assert.Equal(t, "/users", schemas[2].Path)
	assert.Equal(t, "/admin", schemas[3].Path)

	assert.Equal(t, http.StatusOK, schemas[0].Status)
	assert.Equal(t, http.StatusOK, schemas[1].Status)
	assert.Equal(t, http.StatusCreated, schemas[2].Status)
	assert.Equal(t, http.StatusUnauthorized, schemas[3].Status)
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

func TestJSONWithMissingKeys(t *testing.T) {
	_, err := BytesToSchema([]byte(missingKeys))
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

func TestInvalidJSON(t *testing.T) {
	_, err := BytesToSchema([]byte(invalidJSON))
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

func TestJSONObject(t *testing.T) {
	_, err := BytesToSchema([]byte(jsonObject))
	assert.NotNil(t, err)
}
