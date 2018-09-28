package main

import (
	"flag"
	log "github.com/sirupsen/logrus"

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

	cmd, err := json.NewCommand()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Execute(*port, args)
	if err != nil {
		log.Fatal(err)
	}
}
