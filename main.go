package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

func db() {
	dbName := "local.db"
	primaryUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		log.Println("Error creating temporary directory: ", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl, libsql.WithAuthToken(authToken))
	if err != nil {
		log.Println("Error creating connector: ", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()
}

func main() {
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var menuItems []string
		for _, urlValue := range r.Form {
			menuItems = urlValue
		}
		coffeeTemplate := coffee(menuItems)
		coffeeTemplate.Render(context.Background(), w)

	})

	sm.HandleFunc("/trial", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		transactionType := r.FormValue("transaction-type")

		url := r.URL
		w.Write([]byte(fmt.Sprintf("%v\n%v\n%v", url, user, transactionType)))
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: sm,
	}

	log.Fatalln(server.ListenAndServe())
}
