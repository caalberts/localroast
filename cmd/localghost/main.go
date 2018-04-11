package main

import (
	"flag"

	"github.com/caalberts/localghost/http"
)

var port = flag.String("p", "8080", "port number")

func main() {
	flag.Parse()
	http.NewServer(*port).ListenAndServe()
}
