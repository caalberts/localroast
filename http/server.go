package http

import (
	"log"
	"net/http"

	"github.com/caalberts/localroast"
)

func NewServer(port string, schemas []localroast.Schema) *http.Server {
	mux := NewMux(schemas)

	log.Println("Localroast brewing on port " + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

type Mux map[string]responses
type responses map[string]int

func NewMux(schemas []localroast.Schema) Mux {
	mux := make(Mux)
	for _, schema := range schemas {
		if _, exists := mux[schema.Path]; !exists {
			mux[schema.Path] = make(responses)
		}
		mux[schema.Path][schema.Method] = schema.StatusCode
	}
	return mux
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := m[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
	}

	response, ok := handler[r.Method]
	if !ok {
		http.NotFound(w, r)
	}

	w.WriteHeader(response)
}
