package workerpool

import (
	"net/http"
	"time"
)

// A worker using an HTTP client with a specified timeout
type worker struct {
	client *http.Client
}

// Creates a new worker with an HTTP client that has a timeout set.
func newWorker(timeout time.Duration) *worker {
	return &worker{
		&http.Client{
			Timeout: timeout,
		},
	}
}

// Processes a task: makes an HTTP request, measures the response time and returns the result.
func (w worker) process(j Job) Result {
	result := Result{URL: j.URL}

	now := time.Now()

	resp, err := w.client.Get(j.URL)
	if err != nil {
		result.Error = err

		return result
	}

	result.StatusCode = resp.StatusCode
	result.ResponseTime = time.Since(now)

	return result
}
