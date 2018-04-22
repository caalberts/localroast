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
			Method: "GET",
			Path:   "/",
			Status: 200,
		},
		localroast.Schema{
			Method: "GET",
			Path:   "/users",
			Status: 200,
		},
		localroast.Schema{
			Method: "POST",
			Path:   "/users",
			Status: 201,
		},
	}

	mux := NewMux(schemas)

	for _, schema := range schemas {
		req := httptest.NewRequest(schema.Method, schema.Path, nil)
		resp := httptest.NewRecorder()
		mux.ServeHTTP(resp, req)

		assert.Equal(t, schema.Status, resp.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
