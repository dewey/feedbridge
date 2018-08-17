package scrape

import (
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// TaskDocumentResponse is the result of a fetched task and it's originating URL
type TaskDocumentResponse struct {
	Document goquery.Document
	URL      string
}

// Task is a scrape task
type Task struct {
	URL     string
	Payload string
}

type Statistic struct {
	Items int
}

// URLToDocument is a convinience function to directly get a goquery document(s) from a list of URLs
func URLToDocument(c *http.Client, tasks []Task) ([]TaskDocumentResponse, error) {
	if len(tasks) < 1 {
		return nil, errors.New("urls can't be empty")
	}
	var result []TaskDocumentResponse
	for _, task := range tasks {
		req, err := http.NewRequest("GET", task.URL, nil)
		if err != nil {
			continue
		}
		// They block requests without valid user agent
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.2 Safari/605.1.15")
		req.Header.Set("Accept-Language", "en-us")

		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, err
		}
		rt := TaskDocumentResponse{
			Document: *doc,
			URL:      task.URL,
		}
		result = append(result, rt)
	}
	return result, nil
}

// URLtoTask is a helper function to convert a list of URLs to a list of tasks
func URLtoTask(urls []string) []Task {
	var tasks []Task
	for _, u := range urls {
		tasks = append(tasks, Task{
			URL: u,
		})
	}
	return tasks
}
