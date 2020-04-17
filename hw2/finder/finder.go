package finder

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// A struct to store the results for each request including url and response body.
type result struct {
	url  string
	body string
}

/**
Find matches of search string in provided urls
*/
func FindMatches(search string, urls []string) (matches []string) {
	results := getResultsParallel(urls)

	for _, result := range results {
		if strings.Contains(result.body, search) {
			matches = append(matches, result.url)
		}
	}

	return
}

func getResultsParallel(urls []string) []result {
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
