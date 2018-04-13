package main

import (
	"errors"
	"flag"
	"log"

	"github.com/caalberts/localroast/http"
	"github.com/caalberts/localroast/schema"
)

func main() {
	port := flag.String("p", "8080", "port number")
	flag.Parse()
	args := flag.Args()
	definition := parseInput(args)
	schema, err := schema.FromString(definition)
	handle(err)

	server := http.NewServer(*port, schema)
	log.Fatal(server.ListenAndServe())
}

func parseInput(args []string) string {
	if len(args) < 1 {
		log.Fatal(errors.New("Please define an endpoint in the format '<METHOD> <PATH> <STATUS_CODE>'. e.g 'GET / 200'"))
	}

	return args[0]
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
