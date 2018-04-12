package http

import (
	"log"
	"net/http"

	"github.com/caalberts/localghost"
)

func NewServer(port string, schema localghost.Schema) *http.Server {
	mux := NewMux([]localghost.Schema{schema})

	log.Println("localghost:" + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func NewMux(schemas []localghost.Schema) *http.ServeMux {
	mux := http.NewServeMux()
	for _, schema := range schemas {
		mux.HandleFunc(schema.Path, FromSchema(schema))
	}
	return mux
}
