package json

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"io"
)

func TestParser_Watch(t *testing.T) {
	input := make(chan io.Reader)
	p := NewParser()
	p.Watch(input)

	t.Run("valid input", func(t *testing.T) {
		input <- strings.NewReader(validJSON)
		schemas := <-p.Output()

		assert.Equal(t, 2, len(schemas))

		assert.Equal(t, http.MethodGet, schemas[0].Method)
		assert.Equal(t, http.MethodGet, schemas[1].Method)

		assert.Equal(t, "/", schemas[0].Path)
		assert.Equal(t, "/users", schemas[1].Path)

		assert.Equal(t, http.StatusOK, schemas[0].Status)
		assert.Equal(t, http.StatusOK, schemas[1].Status)

		assert.JSONEq(t, `{ "success": true }`, string(schemas[0].Response))
		assert.JSONEq(t, `{ "success": true, "ids": [1, 2, 3] }`, string(schemas[1].Response))
	})

	t.Run("with multiple input", func(t *testing.T) {
		input <- strings.NewReader(validJSON)
		schemas := <-p.Output()

		assert.Equal(t, 2, len(schemas))

		input <- strings.NewReader(validJSON2)
		schemas = <-p.Output()

		assert.Equal(t, 3, len(schemas))
	})

	t.Run("with missing keys in json does not block", func(t *testing.T) {
		input <- strings.NewReader(missingKeys)
		assert.Empty(t, p.Output())
	})

	t.Run("with invalid json does not block", func(t *testing.T) {
		input <- strings.NewReader(invalidJSON)
		assert.Empty(t, p.Output())
	})

	t.Run("with json object does not block", func(t *testing.T) {
		input <- strings.NewReader(jsonObject)
		assert.Empty(t, p.Output())
	})
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

var validJSON2 = `
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
	},
	{
        "method": "POST",
        "path": "/users",
        "status": 201,
        "response": {
            "success": true
        }
	}
]`

var missingKeys = `
[
	{
		"response": {
			"success": true
		}
	}
]
`

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
