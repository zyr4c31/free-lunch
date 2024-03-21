package main

import (
	"log"
)

func main() {
	db := newDBConnection()
	defer db.Close()

	server := newServer(db)

	log.Fatalln(server.ListenAndServe())
}
