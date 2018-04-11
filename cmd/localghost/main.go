package main

import (
	"flag"
	"os"

	"github.com/caalberts/localghost/http"
	"github.com/caalberts/localghost/schema"
)

var port = flag.String("p", "8080", "port number")

func main() {
	flag.Parse()
	definition := os.Args[1]
	schema, _ := schema.FromString(definition)
	http.NewServer(*port, schema).ListenAndServe()
}
