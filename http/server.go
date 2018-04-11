package http

import (
	"log"
	"net/http"

	"github.com/caalberts/localghost"
)

func NewServer(port string, schema localghost.Schema) *http.Server {
	mux := NewMux(schema)

	log.Println("localghost:" + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func NewMux(schema localghost.Schema) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(schema.Path, FromSchema(schema))
	return mux
}
