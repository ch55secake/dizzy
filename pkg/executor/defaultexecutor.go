package executor

import (
	"github.com/ch55secake/dizzy/pkg/client"
	"github.com/ch55secake/dizzy/pkg/input"
	"github.com/ch55secake/dizzy/pkg/job"
	log "github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

type ExecutionContext struct {
	Filepath       string
	Url            string
	ResponseLength int
	Timeout        time.Duration
	Method         string
	Headers        map[string]string
}

type DefaultExecutor struct {
	job.Dispatcher
	WorkerCount int
	QueueSize   int
}

// Execute will use the dispatcher provided on the struct and then kick off all jobs that are ready to be dispatched
// Will also rely on the queue size and worker count provided when using the dispatcher to execute any jobs
func Execute(ctx ExecutionContext) {

	wl := &input.WordList{}
	err := wl.NewWordList(ctx.Filepath)
	if err != nil {
		return
	}

	log.Debugf("generated a wordlist with size %d\n", wl.Size())

	requests, err := wl.TransformWordListToRequests(ctx.Url)

	var jobs []*job.Job
	for _, request := range requests {
		jobs = append(jobs, job.NewJob(rand.Int(), request))
	}

	// Convert int to float for division, round it, convert back to int
	dispatcher := job.NewDispatcher(int(math.Round(float64(len(requests))/3)), len(requests))

	for _, jobToSubmit := range jobs {
		dispatcher.Submit(jobToSubmit)
	}

	r := client.NewRequester(ctx.Timeout, ctx.Method, ctx.Headers)
	dispatcher.Run(r)

	dispatcher.Wait()
	log.Infof("jobs completed, total jobs: %d\n", len(jobs))
}
