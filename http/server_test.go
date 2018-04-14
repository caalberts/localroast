package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caalberts/localroast"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	schema := localroast.Schema{Path: "/"}
	port := "8888"
	server := NewServer(port, []localroast.Schema{schema})
	assert.Equal(t, ":8888", server.Addr)
}

func TestNewMux(t *testing.T) {
	schemas := []localroast.Schema{
		localroast.Schema{
			Method:     "GET",
			Path:       "/",
			StatusCode: 200,
		},
		localroast.Schema{
			Method:     "POST",
			Path:       "/user",
			StatusCode: 201,
		},
	}

	mux := NewMux(schemas)
	server := httptest.NewServer(mux)
	defer server.Close()

	for _, schema := range schemas {

		var resp *http.Response
		var err error

		switch schema.Method {
		case http.MethodGet:
			resp, err = http.Get(server.URL + schema.Path)
		case http.MethodPost:
			resp, err = http.Post(server.URL+schema.Path, "", nil)
		}

		assert.Nil(t, err)
		assert.Equal(t, schema.StatusCode, resp.StatusCode)
	}
}
