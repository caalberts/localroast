package http

import (
	"net/http"

	"github.com/caalberts/localroast"
	"github.com/julienschmidt/httprouter"
	"sync"
)

// Server interface.
type Server interface {
	ListenAndServe() error
	Watch() chan<- []localroast.Schema
}

type server struct {
	*http.Server
	router *router
}

// ServerFunc is a constructor for a new server.
type ServerFunc func(port string) Server

// NewServer creates a http server running on given port with handlers based on given schema.
func NewServer(port string) Server {
	router := newRouter()

	return &server{
		Server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
		router: router,
	}
}

func (s *server) Watch() chan<- []localroast.Schema {
	updateChan := make(chan []localroast.Schema)
	go func() {
		for {
			schemas := <-updateChan
			s.router.updateSchema(schemas)
		}
	}()
	return updateChan
}

type router struct {
	sync.Mutex
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
	rtr.Lock()
	defer rtr.Unlock()
	rtr.Handler.ServeHTTP(w, r)
}

func (rtr *router) updateSchema(schemas []localroast.Schema) {
	rtr.Lock()
	defer rtr.Unlock()
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
