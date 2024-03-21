package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/zyr4c31/free-lunch/sqlc"
)

func newServer(db *sql.DB) *http.Server {

	sm := http.NewServeMux()

	sm.HandleFunc("GET /restaurants", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := r.Form
		var name string
		if form.Has("name") {
			name = form.Get("name")
		}

		queries := sqlc.New(db)
		restaurants, err := queries.ListRestaurants(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		template := Restaurants(name, restaurants)
		template.Render(r.Context(), w)
	})

	sm.HandleFunc("POST /restaurants", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := r.Form
		var name string
		if form.Has("name") == false {
			template := htmlError("no input")
			template.Render(r.Context(), w)
			return
		}

		queries := sqlc.New(db)
		err := queries.CreateRestaurant(context.Background(), name)
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Add("HX-Refresh", "true")
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	return server
}
