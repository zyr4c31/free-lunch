package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	db, dir, connector := db()
	defer os.RemoveAll(dir)
	defer connector.Close()

	log.Println(createTables(db))

	restaurants, err := selectAllRestaurants(db)
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
}
