package main

import (
	"flag"
	"log"

	"github.com/caalberts/localroast/json"
)

var (
	version string
)

func main() {
	log.Printf("localroast %v", version)

	port := flag.String("port", "8080", "port number")

	flag.Parse()
	args := flag.Args()

	err := json.NewCommand().Execute(*port, args)
	if err != nil {
		log.Fatal(err)
	}
}
