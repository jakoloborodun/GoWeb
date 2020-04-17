package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"goweb/hw2/finder"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Search struct {
	Search string   `json:"search"`
	Sites  []string `json:"sites"`
}

func main() {
	stopChan := make(chan os.Signal)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Post("/search", PostSearchHandler)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan
	log.Print("Shutting down")
}

func PostSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var search Search
	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	search.Sites = finder.FindMatches(search.Search, search.Sites)

	json.NewEncoder(w).Encode(search)
}
