package json

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	Success bool   `json:"success"`
	ID      int    `json:"id"`
	IDs     []int  `json:"ids"`
	Message string `json:"message"`
}

var validJSON, _ = ioutil.ReadFile("../examples/stubs.json")

func TestParseSchemaFromJSON(t *testing.T) {
	p := &Parser{}
	schemas, err := p.Parse(validJSON)

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

	var expected, response testData
	json.Unmarshal([]byte(`{"success": true}`), &expected)
	json.Unmarshal([]byte(schemas[0].Response), &response)
	assert.Equal(t, expected, response)

	json.Unmarshal([]byte(`{
		"success": true,
		"ids": [1, 2, 3]
	}`), &expected)
	json.Unmarshal([]byte(schemas[1].Response), &response)
	assert.Equal(t, expected, response)

	json.Unmarshal([]byte(`{
		"success": true,
		"id": 4
	}`), &expected)
	json.Unmarshal([]byte(schemas[2].Response), &response)
	assert.Equal(t, expected, response)

	json.Unmarshal([]byte(`{
		"success": false,
		"message": "unauthorized"
	}`), &expected)
	json.Unmarshal([]byte(schemas[3].Response), &response)
	assert.Equal(t, expected, response)
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
	p := &Parser{}
	_, err := p.Parse([]byte(missingKeys))
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
	p := &Parser{}
	_, err := p.Parse([]byte(invalidJSON))
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
	p := &Parser{}
	_, err := p.Parse([]byte(jsonObject))
	assert.NotNil(t, err)
}
