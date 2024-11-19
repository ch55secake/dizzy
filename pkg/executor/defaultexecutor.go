package executor

import (
	"fmt"
	"github.com/ch55secake/dizzy/pkg/client"
	"github.com/ch55secake/dizzy/pkg/input"
	"github.com/ch55secake/dizzy/pkg/job"
	"log"
	"math/rand"
	"time"
)

type DefaultExecutor struct {
	Filepath string
	job.Dispatcher
	WorkerCount int
	QueueSize   int
}

// Execute will use the dispatcher provided on the struct and then kick off all jobs that are ready to be dispatched
// Will also rely on the queue size and worker count provided when using the dispatcher to execute any jobs
func (e *DefaultExecutor) Execute(filepath string, url string) {

	wl := &input.WordList{}
	err := wl.NewWordList(filepath)
	if err != nil {
		return
	}

	log.Printf("generated a wordlist with size %d\n", wl.Size())

	requests, err := wl.TransformWordListToRequests(url)

	id := rand.Int()

	var jobs []*job.Job
	for _, request := range requests {
		jobs = append(jobs, job.NewJob(id, request))
	}

	for _, jobToSubmit := range jobs {
		e.Dispatcher.Submit(jobToSubmit)
	}

	r := &client.Requester{
		Method:  "GET",
		Timeout: 10 * time.Second,
		Headers: map[string]string{"Accept": "application/json"},
	}
	e.Dispatcher.Run(r)

	e.Dispatcher.Wait()
	fmt.Printf("jobs completed, total jobs: %d\n", len(jobs))
}
