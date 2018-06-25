package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caalberts/localroast"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	schema := localroast.Schema{Path: "/"}
	port := "8888"
	server := NewServer(port, []localroast.Schema{schema}).(*http.Server)
	assert.Equal(t, ":8888", server.Addr)
}

func TestNewRouter(t *testing.T) {
	schemas := []localroast.Schema{
		localroast.Schema{
			Method:   "GET",
			Path:     "/",
			Status:   200,
			Response: []byte(`{"success": true}`),
		},
		localroast.Schema{
			Method:   "GET",
			Path:     "/users",
			Status:   200,
			Response: []byte(`{"success": true, "ids": [1, 2]}`),
		},
		localroast.Schema{
			Method:   "POST",
			Path:     "/users",
			Status:   201,
			Response: []byte(`{"success": true, "id": 3}`),
		},
	}

	type testData struct {
		Success bool   `json:"success"`
		ID      int    `json:"id"`
		IDs     []int  `json:"ids"`
		Message string `json:"message"`
	}
	router := newRouter(schemas)

	var expected, actual testData

	for _, schema := range schemas {
		t.Run(schema.Method+schema.Path, func(t *testing.T) {
			req := httptest.NewRequest(schema.Method, schema.Path, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
			assert.Equal(t, schema.Status, resp.Code)

			body := resp.Result().Body
			defer body.Close()
			json.NewDecoder(body).Decode(&actual)
			json.Unmarshal(schema.Response, &expected)
			assert.Equal(t, expected, actual)
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	body := resp.Result().Body
	defer body.Close()
	bodyBytes, _ := ioutil.ReadAll(body)
	assert.Equal(t, "404 page not found\n", string(bodyBytes))
}

func TestPathParam(t *testing.T) {
	schema := localroast.Schema{
		Method:   "GET",
		Path:     "/users/:id",
		Status:   200,
		Response: []byte(`{"success": true}`),
	}

	type testData struct {
		Success bool `json:"success"`
	}
	mux := newRouter([]localroast.Schema{schema})

	testPath := "/users/1"

	req := httptest.NewRequest(schema.Method, testPath, nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	assert.Equal(t, schema.Status, resp.Code)
}
