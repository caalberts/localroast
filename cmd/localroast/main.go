package main

import (
	"flag"
	"log"

	"github.com/caalberts/localroast"
	"github.com/caalberts/localroast/http"
	"github.com/caalberts/localroast/json"
	"github.com/caalberts/localroast/strings"
)

type command interface {
	Execute(args []string) ([]localroast.Schema, error)
}

var port = flag.String("port", "8080", "port number")
var useJSON = flag.Bool("json", false, "json")

func main() {
	flag.Parse()
	args := flag.Args()

	var cmd command
	if *useJSON {
		cmd = json.NewCommand()
	} else {
		cmd = strings.NewCommand()
	}

	schema, err := cmd.Execute(args)
	if err != nil {
		log.Fatal(err)
	}

	err = http.NewServer(*port, schema).ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
