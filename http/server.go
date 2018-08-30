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
	router := newRouter()
	router.UpdateSchema(schemas)

	log.Println("brewing on port " + port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}

type router struct {
	http.Handler
}

func newRouter() *router {
	rtr := router{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		}),
	}
	return &rtr
}

func (rtr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.Handler.ServeHTTP(w, r)
}

func (rtr *router) UpdateSchema(schemas []localroast.Schema) {
	router := httprouter.New()
	for _, schema := range schemas {
		router.Handle(schema.Method, schema.Path, handlerFunc(schema))
	}
	rtr.Handler = router
}

func handlerFunc(schema localroast.Schema) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(schema.Status)
		w.Write(schema.Response)
	}
}
