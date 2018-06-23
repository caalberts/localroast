package http

import (
	"log"
	"net/http"

	"github.com/caalberts/localroast"
	"github.com/julienschmidt/httprouter"
)

type Server interface {
	ListenAndServe() error
}

type ServerFunc func(port string, schemas []localroast.Schema) Server

func NewServer(port string, schemas []localroast.Schema) Server {
	router := NewRouter(schemas)

	log.Println("Localroast brewing on port " + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}

func NewRouter(schemas []localroast.Schema) http.Handler {
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
