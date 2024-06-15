package workerpool

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Represents a task containing a URL to which an HTTP request should be made
type Job struct {
	URL string
}

// Contains the results of the request: URL, status code, response time and error (if any). The Info method returns a string represnentation of the result
type Result struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

// Method Info
func (r Result) Info() string {
	if r.Error != nil {
		return fmt.Sprintf("[ERROR] - [%s] - %s", r.URL, r.Error.Error())
	}

	return fmt.Sprintf("[SUCCESS] - [%s] - Status: %d, Response Time: %s", r.URL, r.StatusCode, r.ResponseTime.String())
}

// Manages a pool of workers, contains channels for jobs (jobs) and results (results), a counter of synchronization (wg) and a stop flag (stopped)
type Pool struct {
	worker       *worker
	workersCount int

	jobs    chan Job
	results chan Result

	wg      *sync.WaitGroup
	stopped bool
}

// Pool constructor, initializes a pool with a given number of workers, a timeout for the HTTP client, and a channel for the results.
func New(workersCount int, timeout time.Duration, results chan Result) *Pool {
	return &Pool{
		worker:       newWorker(timeout),
		workersCount: workersCount,
		jobs:         make(chan Job),
		results:      results,
		wg:           new(sync.WaitGroup),
	}
}

// Starts the specified number of worker threads
func (p *Pool) Init() {
	for i := 0; i < p.workersCount; i++ {
		go p.initWorker(i)
	}
}

// Adds a new task to the task queue and increases the task counter.
func (p *Pool) Push(j Job) {
	if p.stopped {
		return
	}

	p.jobs <- j
	p.wg.Add(1)
}

// Stops the pool by closing the task channel and waiting for all tasks to complete
func (p *Pool) Stop() {
	p.stopped = true
	close(p.jobs)
	p.wg.Wait()
}

// Basic workflow logic: waiting for tasks from the channel, processing the tasks and sending the results to the results channel.
func (p *Pool) initWorker(id int) {
	for job := range p.jobs {
		time.Sleep(time.Second)
		p.results <- p.worker.process(job)
		p.wg.Done()
	}

	log.Printf("[worker ID %d] finished proccesing", id)
}
