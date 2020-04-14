package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func GetURLsForSearch() []string {
	return []string{
		"https://drushcommands.com/",
		"https://ru.wikipedia.org/wiki/Drush",
		"https://music.yandex.ru/home",
		"https://geekbrains.ru/education",
	}
}

func findMatches(search string, urls []string) (matches []string) {
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		respText := string(bts)
		if strings.Contains(respText, search) {
			matches = append(matches, url)
		}
	}

	return
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	t, _ := template.ParseFiles("search.html")
	_ = t.Execute(w, nil)
	fmt.Fprintf(w, "Pages to search for: %v<br>", GetURLsForSearch()) // write data to response
	if r.Method == "POST" {
		_ = r.ParseForm()
		fmt.Println(r.Form) // print information on server side.

		searchQuery := r.PostFormValue("search")
		fmt.Fprintf(w, "Your search query: %s<br>", searchQuery) // write data to response

		matches := findMatches(searchQuery, GetURLsForSearch())
		fmt.Fprintln(w, "Found matches on the followed pages:<br>")
		for _, match := range matches {
			fmt.Fprintf(w, "%s<br>", match)
		}
	}
}

func route() {
	http.HandleFunc("/search", SearchHandler)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	go route()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt
	log.Println("Shutting down")
}
