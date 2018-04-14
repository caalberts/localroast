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
	definitions := parseInput(args)
	schemas, err := schema.FromStrings(definitions)
	handle(err)

	server := http.NewServer(*port, schemas)
	log.Fatal(server.ListenAndServe())
}

func parseInput(args []string) []string {
	if len(args) < 1 {
		log.Fatal(errors.New("Please define an endpoint in the format '<METHOD> <PATH> <STATUS_CODE>'. e.g 'GET / 200'"))
	}
	return args
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
