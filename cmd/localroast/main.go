package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
	"github.com/caalberts/localroast/schema"
)

type parser interface {
	CreateSchema() ([]localroast.Schema, error)
}

func main() {
	port := flag.String("port", "8080", "port number")
	json := flag.Bool("json", false, "json")

	flag.Parse()
	args := flag.Args()

	var p parser
	if *json {
		bytes := readJSON(args)
		p = &schema.JSON{Bytes: bytes}
	} else {
		input := parseInput(args)
		p = &schema.String{Strings: input}
	}

	schemas, err := p.CreateSchema()
	handle(err)

	server := http.NewServer(*port, schemas)
	log.Fatal(server.ListenAndServe())
}

func readJSON(args []string) []byte {
	filepath := args[0]
	bytes, err := ioutil.ReadFile(filepath)
	handle(err)

	return bytes
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
