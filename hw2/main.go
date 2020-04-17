package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
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
		r.Get("/search", ReadCookieHandler)
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan
	log.Print("Shutting down")
}

func PostSearchHandler(w http.ResponseWriter, r *http.Request) {
	SetCookieHandler(w, r)
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

func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "GoName",
		Value: uuid.Must(uuid.NewV4()).String(),
		Path:  "/search",
	}

	http.SetCookie(w, cookie)
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("GoName"); err == nil {
		w.Write([]byte(cookie.Value))
	}
}
