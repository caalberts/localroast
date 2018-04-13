package http

import (
	"net/http"

	"github.com/caalberts/localroast"
)

func FromSchema(schema localroast.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == schema.Method {
			w.WriteHeader(schema.StatusCode)
		}
	}
}
