package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dbName := "local.db"
	primaryUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_TOKEN")
	syncInterval := time.Minute

	tursoCfg := TursoEmbeddedReplicaConfig{
		dbName:       dbName,
		primaryUrl:   primaryUrl,
		authToken:    authToken,
		syncInterval: syncInterval,
	}

	log.Println(tursoCfg.createTables())

	restaurants, err := tursoCfg.selectAllRestaurants()
	if err != nil {
		log.Println(err)
	}

	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var menuItems []string
		for _, urlValue := range r.Form {
			menuItems = urlValue
		}
		coffeeTemplate := coffee(menuItems, restaurants)
		coffeeTemplate.Render(context.Background(), w)
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: sm,
	}

	log.Fatalln(server.ListenAndServe())
	defer os.RemoveAll(dir)
}
