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
	"time"
)

func GetURLsForSearch() []string {
	return []string{
		"https://drushcommands.com/",
		"https://ru.wikipedia.org/wiki/Drush",
		"https://music.yandex.ru/home",
		"https://geekbrains.ru/education",
		"https://drushcommands.com/",
		"https://ru.wikipedia.org/wiki/Drush",
		"https://music.yandex.ru/home",
		"https://geekbrains.ru/education",
	}
}

// A struct to store the results for each request including url and response body.
type result struct {
	url  string
	body string
}

func getResultsParallel(search string, urls []string) []result {
	// The channel collect the http request results
	resultsCh := make(chan *result)

	// Make sure we close the channel when we're done
	defer func() {
		close(resultsCh)
	}()

	for _, url := range urls {
		go func(url string) {

			// Send the request and put the response in a result struct
			// along with the url.
			resp, _ := http.Get(url)
			defer resp.Body.Close()

			bts, _ := ioutil.ReadAll(resp.Body)
			respText := string(bts)
			result := &result{url, respText}

			// Send the result struct through the resultsCh
			resultsCh <- result

		}(url)
	}

	var results []result

	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsCh
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(urls) {
			break
		}
	}

	return results
}

func findMatches(search string, urls []string) (matches []string) {
	start := time.Now()
	results := getResultsParallel(search, urls)

	for _, result := range results {
		if strings.Contains(result.body, search) {
			matches = append(matches, result.url)
		}
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
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
