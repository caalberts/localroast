package http

import (
	"log"
	"net/http"

	"github.com/caalberts/localroast"
	"github.com/julienschmidt/httprouter"
)

// Server interface.
type Server interface {
	ListenAndServe() error
}

// ServerFunc is a constructor for a new server.
type ServerFunc func(port string, schemas []localroast.Schema) Server

// NewServer creates a http server running on given port with handlers based on given schema.
func NewServer(port string, schemas []localroast.Schema) Server {
	router := newRouter(schemas)

	log.Println("Localroast brewing on port " + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}

func newRouter(schemas []localroast.Schema) http.Handler {
	router := httprouter.New()

	handlerFunc := func(s localroast.Schema) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(s.Status)
			w.Write(s.Response)
		}
	}

	for _, schema := range schemas {
		router.Handle(schema.Method, schema.Path, handlerFunc(schema))
	}

	return router
}
