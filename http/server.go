package http

import (
	"log"
	"net/http"

	"github.com/caalberts/localroast"
)

func NewServer(port string, schema localroast.Schema) *http.Server {
	mux := NewMux([]localroast.Schema{schema})

	log.Println("localroast:" + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func NewMux(schemas []localroast.Schema) *http.ServeMux {
	mux := http.NewServeMux()
	for _, schema := range schemas {
		mux.HandleFunc(schema.Path, FromSchema(schema))
	}
	return mux
}
