package main

import (
	"flag"
	"log"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
	"github.com/caalberts/localroast/json"
)

type command interface {
	Execute(args []string) ([]localroast.Schema, error)
}

var port = flag.String("port", "8080", "port number")

func main() {
	flag.Parse()
	args := flag.Args()

	schema, err := json.NewCommand().Execute(args)
	if err != nil {
		log.Fatal(err)
	}

	err = http.NewServer(*port, schema).ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
