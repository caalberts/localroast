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

type route struct {
	method string
	path   string
}
type Mux map[route]int

func NewMux(schemas []localroast.Schema) Mux {
	mux := make(Mux)
	for _, schema := range schemas {
		route := route{
			method: schema.Method,
			path:   schema.Path,
		}
		mux[route] = schema.Status
	}
	return mux
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := route{
		method: r.Method,
		path:   r.URL.Path,
	}
	responseCode, found := m[route]
	if !found {
		log.Printf("%v: Not Found\n", route)
		http.NotFound(w, r)
		return
	}

	log.Printf("%v: %d\n", route, responseCode)
	w.WriteHeader(responseCode)
}
