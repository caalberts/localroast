package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caalberts/localroast"
	"github.com/stretchr/testify/assert"
)

func TestFromSchema(t *testing.T) {
	schema := localroast.Schema{
		Method:     "GET",
		Path:       "/",
		StatusCode: http.StatusNotFound,
	}

	handler := FromSchema(schema)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	handler(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
}
