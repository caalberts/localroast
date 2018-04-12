package main

import (
	"errors"
	"flag"
	"log"

	"github.com/caalberts/localghost/http"
	"github.com/caalberts/localghost/schema"
)

func main() {
	port := flag.String("p", "8080", "port number")
	flag.Parse()
	args := flag.Args()
	definition := parseInput(args)
	schema, _ := schema.FromString(definition)
	server := http.NewServer(*port, schema)
	log.Fatal(server.ListenAndServe())
}

func parseInput(args []string) string {
	if len(args) < 1 {
		log.Fatal(errors.New("Please define an endpoint in the format '<METHOD> <PATH> <STATUS_CODE>'. e.g 'GET / 200'"))
	}

	return args[0]
}
