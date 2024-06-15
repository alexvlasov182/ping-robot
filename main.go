package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GOLANG-NINJA/pingrobot/workerpool"
)

const (
	INTERVAL        = time.Second * 10
	REQUEST_TIMEOUT = time.Second * 2
	WORKERS_COUNT   = 3
)

var urls = []string{
	"https://facebook-instagram.com/",
	"https://oleksandr-vlasov.com/",
	"https://google.com/",
	"https://golang.org/",
}

// Main function:
// 1. Creates a channel for results and a pool of workers.
// 2. Initializes the pool.
// 3. Runs goroutines to generate tasks and process results.
// 4. Waits for a termination signal (SIGTERM or SIGINT) and stops the pool

func main() {
	results := make(chan workerpool.Result)
	workerPool := workerpool.New(WORKERS_COUNT, REQUEST_TIMEOUT, results)

	workerPool.Init()

	go generateJobs(workerPool)
	go proccessResults(results)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	workerPool.Stop()
}

// Process the results from the results channel, outputing information
func proccessResults(results chan workerpool.Result) {
	go func() {
		for result := range results {
			fmt.Println(result.Info())
		}
	}()
}

// Periodically adds tasks (URLs) to the pool.
func generateJobs(wp *workerpool.Pool) {
	for {
		for _, url := range urls {
			wp.Push(workerpool.Job{URL: url})
		}

		time.Sleep(INTERVAL)
	}
}
