package main

import (
	"flag"
	"fmt"
	"net/http"
)

var port = flag.String("p", "8080", "port number")

func main() {
	flag.Parse()

	fmt.Println("localghost:" + *port)
	http.ListenAndServe(":"+*port, nil)
}
