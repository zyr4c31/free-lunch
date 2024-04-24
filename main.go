package main

import (
	"log"
)

func main() {
	server := newServer()
	log.Fatalln(server.ListenAndServe())
}
