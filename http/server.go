package http

import (
	"log"
	"net/http"

	"github.com/caalberts/localroast"
)

type route struct {
	method   string
	response int
}

var registry = make(map[string][]route)

func NewServer(port string, schemas []localroast.Schema) *http.Server {
	mux := NewMux(schemas)

	log.Println("Localroast brewing on port " + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func NewMux(schemas []localroast.Schema) *http.ServeMux {
	register(schemas)
	return initMux()
}

func register(schemas []localroast.Schema) {
	for _, schema := range schemas {
		r := route{method: schema.Method, response: schema.StatusCode}
		registry[schema.Path] = append(registry[schema.Path], r)
	}
}

func initMux() *http.ServeMux {
	mux := http.NewServeMux()

	for path, routes := range registry {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			for _, route := range routes {
				if r.Method == route.method {
					w.WriteHeader(route.response)
				}
			}
		})
	}

	return mux
}
