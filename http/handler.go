package http

import (
	"net/http"

	"github.com/caalberts/localghost"
)

func FromSchema(schema localghost.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == schema.Method {
			w.WriteHeader(schema.StatusCode)
		}
	}
}
