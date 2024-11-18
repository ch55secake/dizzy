package job

import (
	"github.com/ch55secake/dizzy/pkg/client"
	"sync"
)

// Worker that will complete the tasks on the jobChannel/Queue
type Worker struct {
	ID         int
	JobChannel chan *Job
	Requester  *client.Requester
	wg         *sync.WaitGroup
}

// Start will kick off the processing loop for a given job, will stop when the job has been executed
func (w *Worker) Start() {
	go func() {
		for job := range w.JobChannel {
			job.Execute(w.Requester)
			w.wg.Done()
		}
	}()
}
