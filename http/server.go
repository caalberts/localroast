package http

import (
	"log"
	"net/http"
)

func NewServer(port string) *http.Server {
	log.Println("localghost:" + port)
	return &http.Server{Addr: ":" + port}
}
